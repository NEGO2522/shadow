'use client';

import Link from 'next/link';
import { useState } from 'react';
import { 
  ArrowLeft, 
  MoreHorizontal, 
  ChevronDown, 
  TrendingUp, 
  TrendingDown, 
  Home, 
  Receipt, 
  CreditCard, 
  User 
} from 'lucide-react';

export default function HistoryPage() {
  const [hoveredNav, setHoveredNav] = useState<string | null>(null);

  return (
    /* Added "no-scrollbar" to the main container */
    <div className="app-screen overflow-y-auto no-scrollbar pb-24">

      {/* Header */}
      <div className="screen-header">
        <Link href="/" className="btn-icon-round">
          <ArrowLeft size={18} strokeWidth={2.5} />
        </Link>
        <h1 className="text-base font-semibold tracking-wide">
          History
        </h1>
        <button onClick={() => alert('Options')} className="btn-icon-round">
          <MoreHorizontal size={18} strokeWidth={2.5} />
        </button>
      </div>

      {/* Date Filter */}
      <div className="px-6 mb-6">
        <button onClick={() => alert('Change date range')} className="w-full flex items-center justify-between glass-card p-4">
          <div className="flex items-center gap-2">
            <div className="text-neutral-400">📅</div>
            <span className="text-sm font-medium text-white/90">1 November - 30 November 2022</span>
          </div>
          <ChevronDown size={16} className="text-neutral-400" />
        </button>
      </div>

      {/* Main Chart Card */}
      <div className="px-6 relative z-10 mb-6">
        <div className="rounded-[32px] bg-[#f2ae9a] p-6 shadow-2xl relative overflow-hidden text-black transition-transform duration-500 hover:scale-[1.01]">
          <div className="flex justify-between items-start mb-6">
            <div className="space-y-1">
              <p className="text-black/60 text-sm font-semibold tracking-wide">Total</p>
              <h2 className="text-[34px] font-semibold tracking-tight">
                $10,250.00
              </h2>
            </div>
            <div className="bg-white/30 backdrop-blur-md px-4 py-1.5 rounded-full shadow-sm text-sm font-bold">
              80% Spend
            </div>
          </div>

          {/* Chart placeholder using css bars */}
          <div className="h-40 w-full flex items-end justify-between px-2 pb-4 mt-8 relative border-b border-black/10">
            <div className="absolute left-0 top-0 h-full w-full flex flex-col justify-between text-[10px] text-black/40 font-bold -z-10 pb-4">
              <div className="flex items-center gap-2 w-full"><span className="w-8 shrink-0">$7,000</span><div className="h-px bg-black/10 w-full border-dashed border-b border-black/10"></div></div>
              <div className="flex items-center gap-2 w-full"><span className="w-8 shrink-0">$5,000</span><div className="h-px bg-black/10 w-full border-dashed border-b border-black/10"></div></div>
              <div className="flex items-center gap-2 w-full"><span className="w-8 shrink-0">$3,000</span><div className="h-px bg-black/10 w-full border-dashed border-b border-black/10"></div></div>
              <div className="flex items-center gap-2 w-full"><span className="w-8 shrink-0">$1,000</span><div className="h-px bg-black/10 w-full border-dashed border-b border-black/10"></div></div>
              <div className="flex items-center gap-2 w-full"><span className="w-8 shrink-0">$0</span><div className="h-px bg-black/10 w-full border-dashed border-b border-black/10"></div></div>
            </div>

            {[
              { h: 'h-16', fill: 'h-10', label: '1-7' },
              { h: 'h-32', fill: 'h-16', label: '8-14' },
              { h: 'h-12', fill: 'h-6', label: '15-21' },
              { h: 'h-20', fill: 'h-8', label: '22-28' },
              { h: 'h-16', fill: 'h-12', label: '29-30' }
            ].map((bar, i) => (
              <div key={i} className="flex flex-col items-center gap-2">
                <div className={`w-6 ${bar.h} rounded-t-full rounded-b-full bg-[#d89785] relative overflow-hidden shadow-inner`}>
                  <div className={`absolute bottom-0 left-0 right-0 ${bar.fill} bg-[#272630]`}></div>
                </div>
                <span className="text-[10px] text-black/60 font-medium">{bar.label}</span>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Income & Expenses Row */}
      <div className="px-6 flex gap-4 mb-8">
        <div className="flex-1 bg-[#1e2321] p-4 rounded-3xl border border-white/5 space-y-3 shadow-lg">
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-emerald-500"></div>
            <span className="text-white/60 text-[10px] font-bold uppercase tracking-wider">Expenses</span>
          </div>
          <h3 className="text-xl font-semibold">$ 10,280.00</h3>
          <div className="flex items-center gap-1.5 bg-[#253229] w-max px-2 py-1 rounded-full text-emerald-500">
            <TrendingUp size={12} strokeWidth={3} />
            <span className="text-[10px] font-bold">2.00%</span>
          </div>
        </div>

        <div className="flex-1 bg-[#231a1a] p-4 rounded-3xl border border-white/5 space-y-3 shadow-lg">
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-red-500"></div>
            <span className="text-white/60 text-[10px] font-bold uppercase tracking-wider">Income</span>
          </div>
          <h3 className="text-xl font-semibold">$ 11,280.00</h3>
          <div className="flex items-center gap-1.5 bg-[#312222] w-max px-2 py-1 rounded-full text-red-400">
            <TrendingDown size={12} strokeWidth={3} />
            <span className="text-[10px] font-bold">2.00%</span>
          </div>
        </div>
      </div>

      {/* Popular Spending */}
      <div className="px-6 mb-12">
        <h3 className="text-base font-semibold tracking-wide text-white/90 mb-4">Popular Spending</h3>
        <div className="flex gap-3 overflow-x-auto pb-4 -mx-6 px-6 no-scrollbar">
          {[
            { icon: '🍕', name: 'Food', color: 'bg-[#f2ae9a]' },
            { icon: '🏠', name: 'Household', color: 'bg-[#9ba4ce]' },
            { icon: '📖', name: 'Study', color: 'bg-[#6db392]' }
          ].map((cat, i) => (
            <button key={i} onClick={() => alert(`Filter: ${cat.name}`)} className="flex items-center gap-3 bg-[#1e1e24] border border-white/5 rounded-full px-4 py-2 whitespace-nowrap active:scale-[0.98] transition-all shadow-lg hover:bg-white/5">
              <div className={`w-8 h-8 rounded-full ${cat.color} flex items-center justify-center -ml-2 text-xl shadow-inner`}>
                {cat.icon}
              </div>
              <span className="font-semibold text-sm">{cat.name}</span>
            </button>
          ))}
        </div>
      </div>

      {/* Bottom Floating Navigation */}
      <div className="bottom-nav-v2">
        <div
          className="nav-bar-pill"
          onMouseLeave={() => setHoveredNav(null)}
        >
          {[
            { id: 'home', icon: Home, label: 'Home', href: '/', width: 'w-[42px]' },
            { id: 'activity', icon: Receipt, label: 'Activity', href: '/transactions', width: 'w-[54px]' },
            { id: 'history', icon: CreditCard, label: 'History', href: '/history', width: 'w-[48px]' },
            { id: 'profile', icon: User, label: 'Me', href: '/profile', width: 'w-[24px]' }
          ].map((nav) => (
            <Link
              key={nav.id}
              href={nav.href}
              onMouseEnter={() => setHoveredNav(nav.id)}
              className={`nav-link ${hoveredNav === nav.id || (hoveredNav === null && nav.id === 'history') ? 'nav-link-active' : 'nav-link-inactive'}`}
            >
              <nav.icon size={18} className="shrink-0" />
              <span className={`overflow-hidden text-xs font-semibold transition-all duration-300 ${hoveredNav === nav.id || (hoveredNav === null && nav.id === 'history') ? nav.width + ' opacity-100 ml-1.5' : 'w-0 opacity-0'}`}>
                {nav.label}
              </span>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}