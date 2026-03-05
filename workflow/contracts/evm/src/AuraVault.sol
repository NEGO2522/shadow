// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import {IERC20} from "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import {CCIPReceiver} from "@chainlink/contracts-ccip/src/v0.8/ccip/applications/CCIPReceiver.sol";
import {Client} from "@chainlink/contracts-ccip/src/v0.8/ccip/libraries/Client.sol";
import {ReceiverTemplate} from "./ReceiverTemplate.sol";

/// @title Shadow — Shielded Multi-Chain Pay-by-Human
contract Shadow is CCIPReceiver, ReceiverTemplate {
    IERC20 public immutable usdc;

    event ShieldedDeposit(address indexed sender, bytes32 encryptedRecipient, uint256 amount);
    event ShieldedPayout(address indexed recipient, uint256 amount);

    constructor(
        address _usdc, 
        address _forwarder, 
        address _router
    ) 
        CCIPReceiver(_router)        // CCIP Router for cross-chain
        ReceiverTemplate(_forwarder) // CRE Forwarder for same-chain
    {
        require(_usdc != address(0), "Shadow: zero USDC address");
        usdc = IERC20(_usdc);
    }

    /**
 * @dev Required by Solidity to resolve the conflict between CCIPReceiver and ReceiverTemplate.
 * It ensures the contract correctly identifies that it supports both interfaces.
 */
function supportsInterface(bytes4 interfaceId)
    public
    view
    virtual
    override(CCIPReceiver, ReceiverTemplate)
    returns (bool)
{
    return super.supportsInterface(interfaceId);
}

    /// @notice Same-chain deposit
    function deposit(bytes32 _encryptedRecipient, uint256 _amount) external {
        require(_amount > 0, "Shadow: zero amount");
        require(usdc.transferFrom(msg.sender, address(this), _amount), "Transfer failed");
        emit ShieldedDeposit(msg.sender, _encryptedRecipient, _amount);
    }

    /// @notice Handler for SAME-CHAIN settlement (from CRE Forwarder)
    /// This overrides the virtual function in ReceiverTemplate
    function _processReport(bytes calldata report) internal override {
        (address recipient, uint256 amount) = abi.decode(report, (address, uint256));
        _executePayout(recipient, amount);
    }

    /// @notice Handler for CROSS-CHAIN settlement (from CCIP Router)
    /// This overrides the virtual function in CCIPReceiver
    function _ccipReceive(Client.Any2EVMMessage memory message) internal override {
        // Decode recipient from the data payload
        address recipient = abi.decode(message.data, (address));
        // Get amount from the tokens actually sent via CCIP
        uint256 amount = message.destTokenAmounts[0].amount;
        
        _executePayout(recipient, amount);
    }

    /// @dev Internal helper to keep logic DRY
    function _executePayout(address _recipient, uint256 _amount) internal {
        require(_recipient != address(0), "Shadow: zero recipient");
        require(usdc.transfer(_recipient, _amount), "Payout failed");
        emit ShieldedPayout(_recipient, _amount);
    }
}