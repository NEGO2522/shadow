// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

interface IERC20 {
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
    function transfer(address to, uint256 amount) external returns (bool);
}

/// @title AuraVault — Shielded Pay-by-Human
/// @notice Users deposit USDC with an encrypted recipient. The Chainlink CRE workflow
///         verifies humanity via World ID (in a TEE), decrypts the recipient, and
///         triggers a private payout through the CRE forwarder.
contract AuraVault {
    IERC20 public immutable usdc;

    /// @notice The Chainlink CRE forwarder that is allowed to call onReport.
    address public immutable forwarder;

    /// @notice Emitted when a user makes a shielded deposit.
    /// @param sender    The depositor's address.
    /// @param encryptedRecipient  ECIES/AES-encrypted recipient address (32 bytes).
    /// @param amount    USDC amount deposited (6-decimal units).
    event ShieldedDeposit(
        address indexed sender,
        bytes32 encryptedRecipient,
        uint256 amount
    );

    /// @notice Emitted when the CRE workflow settles a payout.
    /// @param recipient The decrypted recipient address.
    /// @param amount    USDC amount paid out.
    event ShieldedPayout(address indexed recipient, uint256 amount);

    constructor(address _usdc, address _forwarder) {
        require(_usdc != address(0), "AuraVault: zero USDC address");
        require(_forwarder != address(0), "AuraVault: zero forwarder address");
        usdc = IERC20(_usdc);
        forwarder = _forwarder;
    }

    /// @notice Deposit USDC with an encrypted recipient address.
    /// @dev    Caller must have approved this contract to spend `_amount` of USDC first.
    ///         The `_encryptedRecipient` should be the recipient address encrypted
    ///         with the DON's public key so only the TEE can decrypt it.
    /// @param _encryptedRecipient  ECIES/AES-encrypted recipient (32 bytes).
    /// @param _amount              USDC amount to deposit.
    function deposit(bytes32 _encryptedRecipient, uint256 _amount) external {
        require(_amount > 0, "AuraVault: zero amount");
        require(
            usdc.transferFrom(msg.sender, address(this), _amount),
            "AuraVault: transfer failed"
        );
        emit ShieldedDeposit(msg.sender, _encryptedRecipient, _amount);
    }

    /// @notice Called by the Chainlink CRE forwarder to settle a payout.
    /// @dev    The forwarder verifies the DON report before calling this.
    ///         `report` is ABI-encoded as `(address recipient, uint256 amount)`.
    /// @param report  ABI-encoded payout data from the CRE workflow.
    function onReport(bytes calldata report) external {
        require(msg.sender == forwarder, "AuraVault: only forwarder");
        (address recipient, uint256 amount) = abi.decode(report, (address, uint256));
        require(recipient != address(0), "AuraVault: zero recipient");
        require(amount > 0, "AuraVault: zero amount");
        require(
            usdc.transfer(recipient, amount),
            "AuraVault: payout failed"
        );
        emit ShieldedPayout(recipient, amount);
    }
}
