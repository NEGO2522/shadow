'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import {
  ArrowLeft,
  ArrowDownLeft,
  ArrowUpRight,
} from 'lucide-react';
import { BottomNav } from '@/components/BottomNav';

const WALLET_STORAGE = 'shadow_wld_wallet';

interface Tx {
  type: 'sent' | 'received';
  txHash: string;
  blockNumber: string;
  amount: string;
}

function shortHash(hash: string) {
  return hash ? `${hash.slice(0, 8)}…${hash.slice(-6)}` : '';
}

export default function HistoryPage() {
  const [txList, setTxList] = useState<Tx[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const addr = localStorage.getItem(WALLET_STORAGE);
    if (!addr) { setLoading(false); return; }
    fetch(`/api/transactions?address=${addr}`)
      .then((r) => r.json())
      .then((data) => setTxList(data.transactions ?? []))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const totalSent = txList
    .filter((t) => t.type === 'sent')
    .reduce((s, t) => s + parseFloat(t.amount), 0)
    .toFixed(2);

  const totalReceived = txList
    .filter((t) => t.type === 'received')
    .reduce((s, t) => s + parseFloat(t.amount), 0)
    .toFixed(2);

  return (
    <div className="app-screen overflow-y-auto no-scrollbar pb-24">
      {/* Header */}
      <div className="screen-header">
        <Link href="/" className="btn-icon-round">
          <ArrowLeft size={18} strokeWidth={2.5} />
        </Link>
        <h1 className="text-base font-semibold tracking-wide">History</h1>
        <div className="w-10" />
      </div>

      {/* Summary Row */}
      <div className="px-6 flex gap-4 mb-8">
        <div className="flex-1 bg-[#1e2321] p-4 rounded-3xl border border-white/5 space-y-2 shadow-lg">
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-red-500" />
            <span className="text-white/60 text-[10px] font-bold uppercase tracking-wider">Sent</span>
          </div>
          <h3 className="text-xl font-semibold">${totalSent}</h3>
        </div>
        <div className="flex-1 bg-[#231a1a] p-4 rounded-3xl border border-white/5 space-y-2 shadow-lg">
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-emerald-500" />
            <span className="text-white/60 text-[10px] font-bold uppercase tracking-wider">Received</span>
          </div>
          <h3 className="text-xl font-semibold">${totalReceived}</h3>
        </div>
      </div>

      {/* Transaction List */}
      <div className="px-6">
        <h3 className="text-sm font-medium tracking-wide text-white/90 mb-4">All Transactions</h3>

        {loading ? (
          <div className="space-y-3">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="list-item-card animate-pulse">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 rounded-full bg-white/5" />
                  <div className="space-y-2">
                    <div className="h-3 w-28 bg-white/5 rounded-full" />
                    <div className="h-2 w-40 bg-white/5 rounded-full" />
                  </div>
                </div>
                <div className="h-3 w-14 bg-white/5 rounded-full" />
              </div>
            ))}
          </div>
        ) : txList.length === 0 ? (
          <p className="text-center text-white/20 text-sm py-16">No transactions yet</p>
        ) : (
          <div className="space-y-3">
            {txList.map((tx, idx) => (
              <a
                key={idx}
                href={`https://sepolia.etherscan.io/tx/${tx.txHash}`}
                target="_blank"
                rel="noopener noreferrer"
                className="list-item-card block active:scale-[0.98] transition-transform"
              >
                <div className="flex items-center justify-between">
                  <div className="flex items-center gap-4">
                    <div className="w-12 h-12 rounded-full bg-[#1e1e24] flex items-center justify-center border border-white/5">
                      {tx.type === 'sent' ? (
                        <ArrowUpRight size={18} className="text-red-400" />
                      ) : (
                        <ArrowDownLeft size={18} className="text-emerald-400" />
                      )}
                    </div>
                    <div>
                      <p className="font-medium text-white text-sm">
                        {tx.type === 'sent' ? 'Shielded Send' : 'Shielded Receive'}
                      </p>
                      <p className="text-[10px] text-neutral-500 mt-1 font-mono">
                        {shortHash(tx.txHash)}
                      </p>
                    </div>
                  </div>
                  <p className={`font-semibold text-sm ${tx.type === 'sent' ? 'text-red-400' : 'text-emerald-400'}`}>
                    {tx.type === 'sent' ? '-' : '+'}${tx.amount}
                  </p>
                </div>
              </a>
            ))}
          </div>
        )}
      </div>

      <BottomNav />
    </div>
  );
}
