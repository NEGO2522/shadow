// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "@uniswap/v3-periphery/contracts/interfaces/ISwapRouter.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

/// @title SdpSettlement
/// @notice Settles matched dark-pool fills atomically via Uniswap v3.
///         Only the TEE relayer wallet may call `settle()`.
/// Sepolia addresses:
///   Uniswap v3 Router : 0xE592427A0AEce92De3Edee1F18E0157C05861564
///   WETH              : 0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14
///   USDC              : 0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238
contract SdpSettlement {
    ISwapRouter public immutable router;
    address public immutable WETH;
    address public immutable USDC;

    /// Only the TEE relayer wallet (derived from MNEMONIC index 0) may settle.
    address public immutable teeWallet;

    event Settled(
        bytes32 indexed buyId,
        bytes32 indexed sellId,
        uint256 amountIn,
        uint256 amountOut,
        uint256 blockNumber
    );

    error NotTee();
    error SwapFailed();

    constructor(
        address _router,
        address _weth,
        address _usdc,
        address _tee
    ) {
        router    = ISwapRouter(_router);
        WETH      = _weth;
        USDC      = _usdc;
        teeWallet = _tee;
    }

    /// @notice Called by the TEE relayer after a dark-pool match.
    ///         Pulls USDC from the TEE wallet, swaps for WETH on Uniswap v3,
    ///         and emits a Settled event for the fill record.
    ///
    /// @param buyOrderId   UUID of the buy-side order (bytes32, left-padded).
    /// @param sellOrderId  UUID of the sell-side order (bytes32, left-padded).
    /// @param amountIn     USDC to spend (6-decimal, e.g. 2841_000000 = $2841).
    /// @param amountOutMin Minimum WETH to receive (18-decimal); 0.5% slippage guard.
    function settle(
        bytes32 buyOrderId,
        bytes32 sellOrderId,
        uint256 amountIn,
        uint256 amountOutMin
    ) external {
        if (msg.sender != teeWallet) revert NotTee();

        // Pull USDC from the TEE wallet into this contract.
        IERC20(USDC).transferFrom(msg.sender, address(this), amountIn);
        IERC20(USDC).approve(address(router), amountIn);

        // Exact-input single-hop swap: USDC → WETH, 0.05 % pool.
        uint256 amountOut = router.exactInputSingle(
            ISwapRouter.ExactInputSingleParams({
                tokenIn:           USDC,
                tokenOut:          WETH,
                fee:               500,         // 0.05 % pool
                recipient:         msg.sender,  // TEE wallet receives WETH
                deadline:          block.timestamp + 60,
                amountIn:          amountIn,
                amountOutMinimum:  amountOutMin,
                sqrtPriceLimitX96: 0
            })
        );

        emit Settled(buyOrderId, sellOrderId, amountIn, amountOut, block.number);
    }
}
