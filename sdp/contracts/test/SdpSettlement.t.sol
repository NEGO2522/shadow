// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

import "forge-std/Test.sol";
import "../src/SdpSettlement.sol";

/// @notice Unit tests for SdpSettlement access control.
///         Full swap tests require a Sepolia fork (run with --fork-url $RPC_URL).
contract SdpSettlementTest is Test {
    SdpSettlement settlement;

    address tee    = address(0xTEE);
    address router = address(0x1);
    address weth   = address(0x2);
    address usdc   = address(0x3);

    function setUp() public {
        settlement = new SdpSettlement(router, weth, usdc, tee);
    }

    function test_onlyTeeCanSettle() public {
        vm.prank(address(0xBAD));
        vm.expectRevert(SdpSettlement.NotTee.selector);
        settlement.settle(bytes32(0), bytes32(0), 1e6, 0);
    }

    function test_immutablesSetCorrectly() public view {
        assertEq(address(settlement.router()), router);
        assertEq(settlement.WETH(), weth);
        assertEq(settlement.USDC(), usdc);
        assertEq(settlement.teeWallet(), tee);
    }
}
