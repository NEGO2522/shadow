'use client';

import Image from 'next/image';
import Link from 'next/link';
import { useState } from 'react';
import { MessageSquare, ArrowUpRight, Plus, Send, MoreHorizontal, Home, CreditCard, Receipt, User, TrendingUp, QrCode } from 'lucide-react';

export default function Dashboard() {
  const [hoveredNav, setHoveredNav] = useState<string | null>(null);

  return (
    <div className="flex h-[100dvh] flex-col bg-[#141416] tracking-wide text-white overflow-y-auto relative mx-auto sm:border sm:border-neutral-800 max-w-md pb-24 font-sans">

      {/* Header */}
      <div className="flex items-center justify-between px-6 pt-12 pb-6">
        <div>
          <h1 className="text-xl font-medium tracking-tight">
            Hi 👋, Daniel
          </h1>
          <p className="text-[#a1a1aa] text-sm mt-1 tracking-wide">
            Welcome Back
          </p>
        </div>
        <div className="flex items-center gap-4">
          <div className="h-10 w-10 rounded-full bg-[#1a1a20] flex items-center justify-center p-1 border border-white/5 overflow-hidden ring-2 ring-transparent transition hover:ring-white/20">
            <Image
              src="https://api.dicebear.com/7.x/notionists/svg?seed=Daniel"
              alt="User Avatar"
              width={32}
              height={32}
              className="rounded-full bg-white/10"
            />
          </div>
          <button onClick={() => alert('Messages')} className="h-10 w-10 flex items-center justify-center rounded-full bg-[#1e1e24] text-gray-300 relative border border-white/5 active:scale-95 transition-transform">
            <MessageSquare size={18} strokeWidth={2} />
            <div className="absolute top-2 right-2.5 h-1.5 w-1.5 bg-red-500 rounded-full border border-[#1e1e24]"></div>
          </button>
        </div>
      </div>

      {/* Main Card */}
      <div className="px-6 relative z-10">
        <div className="rounded-[32px] bg-gradient-to-br from-[#6b7b96] to-[#45546e] p-6 shadow-2xl relative overflow-hidden">

          {/* Card background shapes */}
          <div className="absolute top-0 right-0 w-40 h-40 bg-white/5 rounded-full blur-3xl -mr-10 -mt-10"></div>

          <div className="flex justify-between items-start">
            <div className="space-y-1">
              <p className="text-white/80 text-sm font-medium">Total Balance</p>
              <h2 className="text-[34px] font-semibold tracking-tight text-white">
                $12,650.00
              </h2>
            </div>
            {/* MasterCard Logo Fake */}
            <div className="flex -space-x-3 opacity-90">
              <div className="w-8 h-8 rounded-full bg-red-500/80 mix-blend-multiply"></div>
              <div className="w-8 h-8 rounded-full bg-yellow-500/80 mix-blend-multiply"></div>
            </div>
          </div>

          <div className="mt-4 flex items-center gap-2">
            <div className="flex items-center gap-1.5 bg-white/10 px-2 py-1 rounded-full backdrop-blur-md">
              <TrendingUp size={12} className="text-white" />
              <span className="text-[10px] text-white font-medium">2.00% than last week</span>
            </div>
          </div>

          <div className="mt-8 flex items-end justify-between text-white/90 font-mono text-sm tracking-widest relative z-10">
            <p>**** **** **** 9946</p>
            <p className="font-sans text-xs tracking-wide">12 / 25</p>
          </div>

          {/* Action Row inside Card */}
          <div className="mt-5 flex gap-3 relative z-10">
            <button onClick={() => alert('Pay flow')} className="flex items-center gap-2 bg-[#1b1c20] text-white px-6 py-3 rounded-[24px] shadow-lg text-sm font-semibold flex-1 justify-center active:scale-95 transition-transform group hover:bg-[#25262c]">
              <div className="w-6 h-6 rounded-full bg-white/10 flex items-center justify-center group-hover:bg-emerald-400/20 transition-colors">
                <QrCode size={14} className="text-emerald-400" />
              </div>
              Pay
            </button>
            <Link href="/transfer" className="flex items-center gap-2 bg-[#8b9dba] text-black px-6 py-3 rounded-[24px] shadow-lg text-sm font-semibold flex-1 justify-center active:scale-95 transition-transform hover:bg-[#a8b8d4]">
              <Send size={14} />
              Transfer
            </Link>
          </div>
        </div>
      </div>

      {/* Latest Transaction */}
      <div className="px-6 mt-10">
        <h3 className="text-sm font-medium tracking-wide text-white/90 mb-4">Latest Transaction</h3>
        <div className="space-y-3">

          <div className="flex items-center justify-between bg-[#1b1b21] p-4 rounded-3xl border border-white/5">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-full bg-[#2a2a32] flex items-center justify-center overflow-hidden">
                <Image src="https://api.dicebear.com/7.x/avataaars/svg?seed=Felix" alt="Felix" width={36} height={36} />
              </div>
              <div>
                <p className="font-medium text-white text-sm">Transfer for Work</p>
                <p className="text-xs text-neutral-500 mt-1">1 Sep 9:29 | <span className="text-neutral-400">Pending</span></p>
              </div>
            </div>
            <p className="text-white font-medium text-sm">+$120</p>
          </div>

          <div className="flex items-center justify-between bg-[#1b1b21] p-4 rounded-3xl border border-white/5">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-full bg-[#fcd4b6]/10 flex items-center justify-center overflow-hidden">
                <Image src="https://api.dicebear.com/7.x/avataaars/svg?seed=Aneka" alt="Aneka" width={36} height={36} />
              </div>
              <div>
                <p className="font-medium text-white text-sm">Transfer for Work</p>
                <p className="text-xs text-neutral-500 mt-1">1 Sep 9:29 | <span className="text-emerald-500">Succes</span></p>
              </div>
            </div>
            <p className="text-white font-medium text-sm">-$120</p>
          </div>

          <div className="flex items-center justify-between bg-[#1b1b21] p-4 rounded-3xl border border-white/5">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 rounded-full bg-[#bde0ca]/10 flex items-center justify-center overflow-hidden">
                <Image src="https://api.dicebear.com/7.x/avataaars/svg?seed=Sara" alt="Sara" width={36} height={36} />
              </div>
              <div>
                <p className="font-medium text-white text-sm">Transfer for Work</p>
                <p className="text-xs text-neutral-500 mt-1">1 Sep 9:29 | <span className="text-emerald-500">Succes</span></p>
              </div>
            </div>
            <p className="text-white font-medium text-sm">+$120</p>
          </div>

        </div>
      </div>

      {/* Bottom Floating Navigation */}
      <div className="fixed bottom-6 left-0 right-0 z-50 flex justify-center w-full px-8 pointer-events-none">
        <div
          className="flex items-center justify-between gap-2 rounded-full bg-[#212128] p-2 shadow-2xl border border-white/5 mx-auto pointer-events-auto"
          onMouseLeave={() => setHoveredNav(null)}
        >

          <Link
            href="/"
            onMouseEnter={() => setHoveredNav('home')}
            className={`flex items-center justify-center rounded-full px-3 py-3 transition-all duration-300 ${hoveredNav === 'home' || hoveredNav === null ? 'bg-white text-black' : 'text-neutral-400 hover:text-white'}`}
          >
            <Home size={18} className={`shrink-0 ${hoveredNav === 'home' || hoveredNav === null ? 'text-black' : ''}`} />
            <span className={`overflow-hidden text-xs font-medium transition-all duration-300 ${hoveredNav === 'home' || hoveredNav === null ? 'w-[42px] opacity-100 ml-1.5 text-black' : 'w-0 opacity-0'}`}>Home</span>
          </Link>

          <Link
            href="/transactions"
            onMouseEnter={() => setHoveredNav('transactions')}
            className={`flex items-center justify-center rounded-full px-3 py-3 transition-all duration-300 ${hoveredNav === 'transactions' ? 'bg-white text-black' : 'text-neutral-400 hover:text-white'}`}
          >
            <Receipt size={18} className={`shrink-0 ${hoveredNav === 'transactions' ? 'text-black' : ''}`} />
            <span className={`overflow-hidden text-xs font-medium transition-all duration-300 ${hoveredNav === 'transactions' ? 'w-[75px] opacity-100 ml-1.5 text-black' : 'w-0 opacity-0'}`}>Transactions</span>
          </Link>

          <Link
            href="/history"
            onMouseEnter={() => setHoveredNav('history')}
            className={`flex items-center justify-center rounded-full px-3 py-3 transition-all duration-300 ${hoveredNav === 'history' ? 'bg-white text-black' : 'text-neutral-400 hover:text-white'}`}
          >
            <CreditCard size={18} className={`shrink-0 ${hoveredNav === 'history' ? 'text-black' : ''}`} />
            <span className={`overflow-hidden text-xs font-medium transition-all duration-300 ${hoveredNav === 'history' ? 'w-[48px] opacity-100 ml-1.5 text-black' : 'w-0 opacity-0'}`}>History</span>
          </Link>

          <Link
            href="/profile"
            onMouseEnter={() => setHoveredNav('profile')}
            className={`flex items-center justify-center rounded-full px-3 py-3 transition-all duration-300 ${hoveredNav === 'profile' ? 'bg-white text-black' : 'text-neutral-400 hover:text-white'}`}
          >
            <User size={18} className={`shrink-0 ${hoveredNav === 'profile' ? 'text-black' : ''}`} />
            <span className={`overflow-hidden text-xs font-medium transition-all duration-300 ${hoveredNav === 'profile' ? 'w-[42px] opacity-100 ml-1.5 text-black' : 'w-0 opacity-0'}`}>Profile</span>
          </Link>

        </div>
      </div>

    </div>
  );
}
