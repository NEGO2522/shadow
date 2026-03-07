// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "forge-std/Script.sol";
import "../src/SdpSettlement.sol";

/// @notice Deploys SdpSettlement to Sepolia.
///
/// Usage:
///   forge script script/Deploy.s.sol:DeployScript \
///     --rpc-url $RPC_URL \
///     --private-key $PRIVATE_KEY \
///     --broadcast \
///     --verify
///
/// Required env vars:
///   PRIVATE_KEY         — deployer key (funded with Sepolia ETH)
///   TEE_WALLET          — address derived from the MNEMONIC index 0
///   RPC_URL             — Alchemy Sepolia endpoint
///   ETHERSCAN_API_KEY   — for --verify
contract DeployScript is Script {
    // Sepolia contract addresses.
    address constant UNISWAP_V3_ROUTER = 0xE592427A0AEce92De3Edee1F18E0157C05861564;
    address constant WETH              = 0xfFf9976782d46CC05630D1f6eBAb18b2324d6B14;
    address constant USDC              = 0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238;

    function run() external {
        address teeWallet = vm.envAddress("TEE_WALLET");

        vm.startBroadcast();
        SdpSettlement settlement = new SdpSettlement(
            UNISWAP_V3_ROUTER,
            WETH,
            USDC,
            teeWallet
        );
        vm.stopBroadcast();

        console.log("SdpSettlement deployed at:", address(settlement));
        console.log("  teeWallet :", teeWallet);
        console.log("  router    :", UNISWAP_V3_ROUTER);
        console.log("  WETH      :", WETH);
        console.log("  USDC      :", USDC);
    }
}
