import { NextResponse } from 'next/server';
import { createPublicClient, http, parseAbiItem, getAddress } from 'viem';
import { sepolia } from 'viem/chains';

const SHADOW_ADDRESS = (process.env.NEXT_PUBLIC_SHADOW_CONTRACT ??
  '0x7aD3A7BeCe749eBaBEDF20fe2c688C525479e5De') as `0x${string}`;
const SEPOLIA_RPC =
  process.env.NEXT_PUBLIC_SEPOLIA_RPC ??
  'https://eth-sepolia.g.alchemy.com/v2/t7Oxw5b_OpDL6yQVWN70ZjxO6hTCaZeW';

const client = createPublicClient({ chain: sepolia, transport: http(SEPOLIA_RPC) });

const shieldedDepositEvent = parseAbiItem(
  'event ShieldedDeposit(address indexed sender, bytes32 encryptedRecipient, uint256 amount)',
);
const shieldedPayoutEvent = parseAbiItem(
  'event ShieldedPayout(address indexed recipient, uint256 amount)',
);

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const address = searchParams.get('address');
  if (!address || !/^0x[0-9a-fA-F]{40}$/.test(address)) {
    return NextResponse.json({ error: 'Invalid address' }, { status: 400 });
  }

  const userAddress = getAddress(address) as `0x${string}`;

  try {
    const [depositLogs, payoutLogs] = await Promise.all([
      client.getLogs({
        address: SHADOW_ADDRESS,
        event: shieldedDepositEvent,
        args: { sender: userAddress },
        fromBlock: BigInt(0),
      }),
      client.getLogs({
        address: SHADOW_ADDRESS,
        event: shieldedPayoutEvent,
        args: { recipient: userAddress },
        fromBlock: BigInt(0),
      }),
    ]);

    const sent = depositLogs.map((log) => ({
      type: 'sent',
      txHash: log.transactionHash,
      blockNumber: log.blockNumber?.toString(),
      amount: (Number(log.args.amount ?? 0) / 1_000_000).toFixed(2),
    }));

    const received = payoutLogs.map((log) => ({
      type: 'received',
      txHash: log.transactionHash,
      blockNumber: log.blockNumber?.toString(),
      amount: (Number(log.args.amount ?? 0) / 1_000_000).toFixed(2),
    }));

    const all = [...sent, ...received].sort(
      (a, b) => Number(b.blockNumber) - Number(a.blockNumber),
    );

    return NextResponse.json({ transactions: all });
  } catch (err) {
    console.error('transactions fetch error:', err);
    return NextResponse.json({ error: 'Failed to fetch transactions' }, { status: 500 });
  }
}
