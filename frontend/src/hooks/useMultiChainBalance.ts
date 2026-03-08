'use client';

import { useState, useEffect } from 'react';
import { createPublicClient, http, formatUnits } from 'viem';
import { sepolia, arbitrumSepolia } from 'viem/chains';
import type { Chain } from 'viem';

// Worldchain Sepolia chain definition (not yet in viem's default exports)
const worldchainSepolia = {
  id: 4801,
  name: 'World Chain Sepolia',
  nativeCurrency: { name: 'Ether', symbol: 'ETH', decimals: 18 },
  rpcUrls: {
    default: {
      http: [
        'https://worldchain-sepolia.g.alchemy.com/v2/t7Oxw5b_OpDL6yQVWN70ZjxO6hTCaZeW',
      ],
    },
  },
} satisfies Chain;

const USDC_ABI = [
  {
    name: 'balanceOf',
    type: 'function',
    stateMutability: 'view',
    inputs: [{ name: 'account', type: 'address' }],
    outputs: [{ name: '', type: 'uint256' }],
  },
] as const;

const CHAINS = [
  {
    chain: sepolia,
    rpc: 'https://eth-sepolia.g.alchemy.com/v2/t7Oxw5b_OpDL6yQVWN70ZjxO6hTCaZeW',
    // Official USDC on Ethereum Sepolia
    usdc: '0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238' as `0x${string}`,
  },
  {
    chain: worldchainSepolia,
    rpc: 'https://worldchain-sepolia.g.alchemy.com/v2/t7Oxw5b_OpDL6yQVWN70ZjxO6hTCaZeW',
    // USDC on Worldchain Sepolia
    usdc: '0x79A02482A880bCE3F13e09Da970dC34db4CD24d1' as `0x${string}`,
  },
  {
    chain: arbitrumSepolia,
    rpc: 'https://arb-sepolia.g.alchemy.com/v2/t7Oxw5b_OpDL6yQVWN70ZjxO6hTCaZeW',
    // USDC on Arbitrum Sepolia
    usdc: '0x75faf114eafb1BDbe2F0316DF893fd58CE46AA4d' as `0x${string}`,
  },
] as const;

export function useMultiChainBalance(address: `0x${string}` | undefined) {
  const [balance, setBalance] = useState<string>('0.00');
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (!address) return;

    async function fetchBalances() {
      setLoading(true);
      let total = BigInt(0);

      await Promise.allSettled(
        CHAINS.map(async ({ chain, rpc, usdc }) => {
          try {
            const client = createPublicClient({
              chain: chain as Chain,
              transport: http(rpc),
            });
            const bal = await client.readContract({
              address: usdc,
              abi: USDC_ABI,
              functionName: 'balanceOf',
              args: [address!],
            });
            total = total + bal;
          } catch {
            // chain unavailable — skip
          }
        }),
      );

      setBalance(parseFloat(formatUnits(total, 6)).toFixed(2));
      setLoading(false);
    }

    fetchBalances();
  }, [address]);

  return { balance, loading };
}
