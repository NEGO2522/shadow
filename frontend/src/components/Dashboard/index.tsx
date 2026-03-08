'use client';

import Image from 'next/image';
import Link from 'next/link';
import { useState, useEffect } from 'react';
import {
  MessageSquare,
  Send,
  TrendingUp,
  Shield,
  ArrowDownLeft,
  ArrowUpRight,
} from 'lucide-react';
import { BottomNav } from '@/components/BottomNav';
import { useMultiChainBalance } from '@/hooks/useMultiChainBalance';

const WALLET_STORAGE = 'shadow_wld_wallet';

interface Tx {
  type: 'sent' | 'received';
  txHash: string;
  blockNumber: string;
  amount: string;
}

function shortHash(hash: string) {
  return hash ? `${hash.slice(0, 6)}…${hash.slice(-4)}` : '';
}

export default function Dashboard() {
  const [isShielded, setIsShielded] = useState(false);
  const [walletAddress, setWalletAddress] = useState<`0x${string}` | undefined>();
  const [txList, setTxList] = useState<Tx[]>([]);
  const [txLoading, setTxLoading] = useState(false);

  const { balance, loading } = useMultiChainBalance(walletAddress);

  // Load wallet address from localStorage (set after World ID verification)
  useEffect(() => {
    const addr = localStorage.getItem(WALLET_STORAGE);
    if (addr) setWalletAddress(addr as `0x${string}`);
  }, []);

  // Fetch real on-chain transactions
  useEffect(() => {
    if (!walletAddress) return;
    setTxLoading(true);
    fetch(`/api/transactions?address=${walletAddress}`)
      .then((r) => r.json())
      .then((data) => setTxList(data.transactions ?? []))
      .catch(() => {})
      .finally(() => setTxLoading(false));
  }, [walletAddress]);

  const displayAddress = walletAddress
    ? `${walletAddress.slice(0, 6)}…${walletAddress.slice(-4)}`
    : '—';

  return (
    <div className="app-screen pb-24">
      {/* ── Header ────────────────────────────────────────────────── */}
      <div className="screen-header">
        <div>
          <h1 className="text-xl font-medium tracking-tight">Hi 👋</h1>
          <p className="subtitle-sm mt-1 font-mono text-xs">{displayAddress}</p>
        </div>

        <div className="flex items-center gap-2">
          <Link href="/transfer" className="header-action-btn header-action-btn--primary">
            <ArrowUpRight size={12} />
            Pay
          </Link>

          <button className="header-action-btn header-action-btn--dark">
            <ArrowDownLeft size={12} />
            Receive
          </button>

          <button
            onClick={() => alert('Messages')}
            className="btn-icon-round relative"
          >
            <MessageSquare size={18} strokeWidth={2} />
            <div className="absolute top-2 right-2.5 h-1.5 w-1.5 bg-red-500 rounded-full border border-[#1e1e24]" />
          </button>

          <div className="relative">
            <div className="h-10 w-10 rounded-full bg-[#1a1a20] flex items-center justify-center p-1 border border-white/5 overflow-hidden ring-2 ring-emerald-400/40">
              <Image
                src={`https://api.dicebear.com/7.x/notionists/svg?seed=${walletAddress ?? 'shadow'}`}
                alt="User Avatar"
                width={32}
                height={32}
                className="rounded-full bg-white/10"
              />
            </div>
            <div className="verified-badge">
              <Shield size={7} className="text-black" />
            </div>
          </div>
        </div>
      </div>

      {/* ── Balance Card ──────────────────────────────────────────── */}
      <div className="px-6 relative z-10">
        <div className="gradient-card">
          <div className="absolute top-0 right-0 w-40 h-40 bg-white/5 rounded-full blur-3xl -mr-10 -mt-10" />

          <div className="flex justify-between items-center mb-5 relative z-10">
            <div className="flex items-center gap-2">
              <Shield
                size={13}
                className={isShielded ? 'text-emerald-400' : 'text-white/30'}
              />
              <span className="text-[11px] font-semibold text-white/70 tracking-wide">
                {isShielded ? 'Shielded Mode' : 'Standard Mode'}
              </span>
            </div>
            <button
              onClick={() => setIsShielded(!isShielded)}
              className={`shield-toggle ${isShielded ? 'shield-toggle--on' : 'shield-toggle--off'}`}
              aria-label="Toggle shielded mode"
            >
              <div
                className={`shield-toggle__thumb ${isShielded ? 'translate-x-5' : 'translate-x-0.5'}`}
              />
            </button>
          </div>

          <div className="flex justify-between items-start relative z-10">
            <div className="space-y-1">
              <p className="text-white/80 text-sm font-medium">Total Balance</p>
              <h2 className="title-xl">
                {loading ? (
                  <span className="text-white/40 animate-pulse">···</span>
                ) : (
                  `$${balance}`
                )}
              </h2>
            </div>
            <div className="flex -space-x-3 opacity-90">
              <div className="w-8 h-8 rounded-full bg-red-500/80 mix-blend-multiply" />
              <div className="w-8 h-8 rounded-full bg-yellow-500/80 mix-blend-multiply" />
            </div>
          </div>

          <div className="mt-4 flex items-center gap-2 relative z-10">
            <div className="flex items-center gap-1.5 bg-white/10 px-2 py-1 rounded-full backdrop-blur-md">
              <TrendingUp size={12} className="text-white" />
              <span className="text-[10px] text-white font-medium">
                Unified across chains
              </span>
            </div>
          </div>

          <div className="mt-8 flex items-end justify-between text-white/90 font-mono text-sm tracking-widest relative z-10">
            <p className="text-xs tracking-wider">{displayAddress}</p>
            <p className="font-sans text-xs tracking-wide">Sepolia</p>
          </div>

          <div className="mt-6 flex gap-3 relative z-10">
            <button className="btn-action-dark group">
              <ArrowDownLeft size={14} className="text-emerald-400" />
              Receive
            </button>
            <Link
              href="/transfer"
              className={`btn-action-primary ${isShielded ? 'btn-action-shielded' : ''}`}
            >
              <Send size={14} />
              {isShielded ? 'Shielded Pay' : 'Transfer'}
            </Link>
          </div>
        </div>
      </div>

      {/* ── Latest Transactions ───────────────────────────────────── */}
      <div className="px-6 mt-10">
        <div className="flex justify-between items-center mb-4 px-1">
          <h3 className="text-sm font-medium tracking-wide text-white/90">
            Latest Transactions
          </h3>
          <Link href="/history" className="text-[10px] text-white/40 hover:text-white/70 font-mono uppercase tracking-wider">
            See all
          </Link>
        </div>

        {txLoading ? (
          <div className="space-y-3">
            {[1, 2, 3].map((i) => (
              <div key={i} className="list-item-card animate-pulse">
                <div className="flex items-center gap-4">
                  <div className="w-12 h-12 rounded-full bg-white/5" />
                  <div className="space-y-2">
                    <div className="h-3 w-28 bg-white/5 rounded-full" />
                    <div className="h-2 w-16 bg-white/5 rounded-full" />
                  </div>
                </div>
                <div className="h-3 w-14 bg-white/5 rounded-full" />
              </div>
            ))}
          </div>
        ) : txList.length === 0 ? (
          <p className="text-center text-white/20 text-sm py-8">No transactions yet</p>
        ) : (
          <div className="space-y-3">
            {txList.slice(0, 5).map((tx, idx) => (
              <div key={idx} className="list-item-card">
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
                    <p className="text-[10px] uppercase tracking-wider text-neutral-500 mt-1 font-bold font-mono">
                      {shortHash(tx.txHash)}
                    </p>
                  </div>
                </div>
                <p className={`font-semibold text-sm ${tx.type === 'sent' ? 'text-red-400' : 'text-emerald-400'}`}>
                  {tx.type === 'sent' ? '-' : '+'}${tx.amount}
                </p>
              </div>
            ))}
          </div>
        )}
      </div>

      <BottomNav />
    </div>
  );
}
