'use client';

import Link from 'next/link';
import { ArrowLeft, MoreHorizontal, ChevronDown, TrendingUp, TrendingDown, Pizza, Home as HomeIcon, BookOpen } from 'lucide-react';

export default function HistoryPage() {
    return (
        <div className="flex h-[100dvh] flex-col bg-[#141416] tracking-wide text-white overflow-y-auto relative mx-auto sm:border sm:border-neutral-800 max-w-md pb-6 font-sans">

            {/* Header */}
            <div className="flex items-center justify-between px-6 pt-12 pb-6">
                <Link href="/" className="h-10 w-10 flex items-center justify-center rounded-full bg-[#1e1e24] text-gray-300 relative border border-white/5 active:scale-95 transition-transform">
                    <ArrowLeft size={18} strokeWidth={2.5} />
                </Link>
                <h1 className="text-base font-semibold tracking-wide">
                    History
                </h1>
                <button onClick={() => alert('Options')} className="h-10 w-10 flex items-center justify-center rounded-full bg-[#1e1e24] text-gray-300 relative border border-white/5 active:scale-95 transition-transform">
                    <MoreHorizontal size={18} strokeWidth={2.5} />
                </button>
            </div>

            {/* Date Filter */}
            <div className="px-6 mb-6">
                <button onClick={() => alert('Change date range')} className="w-full flex items-center justify-between bg-[#1e1e24] border border-white/5 rounded-2xl p-4 active:scale-[0.98] transition-transform cursor-pointer hover:bg-white/5">
                    <div className="flex items-center gap-2">
                        <div className="text-neutral-400">📅</div>
                        <span className="text-sm font-medium text-white/90">1 November - 30 November 2022</span>
                    </div>
                    <ChevronDown size={16} className="text-neutral-400" />
                </button>
            </div>

            {/* Main Chart Card */}
            <div className="px-6 relative z-10 mb-6">
                <div className="rounded-[32px] bg-[#f2ae9a] p-6 shadow-2xl relative overflow-hidden text-black">
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
                        {/* Y Axis steps */}
                        <div className="absolute left-0 top-0 h-full w-full flex flex-col justify-between text-[10px] text-black/40 font-bold -z-10 pb-4">
                            <div className="flex items-center gap-2 w-full"><span className="w-8 shrink-0">$7,000</span><div className="h-px bg-black/10 w-full border-dashed border-b border-black/10"></div></div>
                            <div className="flex items-center gap-2 w-full"><span className="w-8 shrink-0">$5,000</span><div className="h-px bg-black/10 w-full border-dashed border-b border-black/10"></div></div>
                            <div className="flex items-center gap-2 w-full"><span className="w-8 shrink-0">$3,000</span><div className="h-px bg-black/10 w-full border-dashed border-b border-black/10"></div></div>
                            <div className="flex items-center gap-2 w-full"><span className="w-8 shrink-0">$1,000</span><div className="h-px bg-black/10 w-full border-dashed border-b border-black/10"></div></div>
                            <div className="flex items-center gap-2 w-full"><span className="w-8 shrink-0">$0</span><div className="h-px bg-black/10 w-full border-dashed border-b border-black/10"></div></div>
                        </div>

                        <div className="flex flex-col items-center gap-2">
                            <div className="w-6 h-16 rounded-t-full rounded-b-full bg-[#d89785] relative overflow-hidden">
                                <div className="absolute bottom-0 left-0 right-0 h-10 bg-[#272630]"></div>
                            </div>
                            <span className="text-[10px] text-black/60 font-medium">1-7</span>
                        </div>
                        <div className="flex flex-col items-center gap-2">
                            <div className="w-6 h-32 rounded-t-full rounded-b-full bg-[#d89785] relative overflow-hidden">
                                <div className="absolute bottom-0 left-0 right-0 h-16 bg-[#272630]"></div>
                            </div>
                            <span className="text-[10px] text-black/60 font-medium">8-14</span>
                        </div>
                        <div className="flex flex-col items-center gap-2">
                            <div className="w-6 h-12 rounded-t-full rounded-b-full bg-[#d89785] relative overflow-hidden">
                                <div className="absolute bottom-0 left-0 right-0 h-6 bg-[#272630]"></div>
                            </div>
                            <span className="text-[10px] text-black/60 font-medium">15-21</span>
                        </div>
                        <div className="flex flex-col items-center gap-2">
                            <div className="w-6 h-20 rounded-t-full rounded-b-full bg-[#d89785] relative overflow-hidden">
                                <div className="absolute bottom-0 left-0 right-0 h-8 bg-[#272630]"></div>
                            </div>
                            <span className="text-[10px] text-black/60 font-medium">22-28</span>
                        </div>
                        <div className="flex flex-col items-center gap-2">
                            <div className="w-6 h-16 rounded-t-full rounded-b-full bg-[#d89785] relative overflow-hidden">
                                <div className="absolute bottom-0 left-0 right-0 h-12 bg-[#272630]"></div>
                            </div>
                            <span className="text-[10px] text-black/60 font-medium">29-30</span>
                        </div>
                    </div>
                </div>
            </div>

            {/* Income & Expenses Row */}
            <div className="px-6 flex gap-4 mb-8">

                {/* Expenses */}
                <div className="flex-1 bg-[#1e2321] p-4 rounded-3xl border border-white/5 space-y-3">
                    <div className="flex items-center gap-2">
                        <div className="w-2 h-2 rounded-full bg-emerald-500"></div>
                        <span className="text-white/60 text-xs font-semibold uppercase tracking-wider">Expenses</span>
                    </div>
                    <h3 className="text-xl font-semibold">$ 10,280.00</h3>
                    <div className="flex items-center gap-1.5 bg-[#253229] w-max px-2 py-1 rounded-full text-emerald-500">
                        <TrendingUp size={12} strokeWidth={3} />
                        <span className="text-[10px] font-bold">2.00% than last month</span>
                    </div>
                </div>

                {/* Income */}
                <div className="flex-1 bg-[#231a1a] p-4 rounded-3xl border border-white/5 space-y-3">
                    <div className="flex items-center gap-2">
                        <div className="w-2 h-2 rounded-full bg-red-500"></div>
                        <span className="text-white/60 text-xs font-semibold uppercase tracking-wider">Income</span>
                    </div>
                    <h3 className="text-xl font-semibold">$ 11,280.00</h3>
                    <div className="flex items-center gap-1.5 bg-[#312222] w-max px-2 py-1 rounded-full text-red-400">
                        <TrendingDown size={12} strokeWidth={3} />
                        <span className="text-[10px] font-bold">2.00% than last month</span>
                    </div>
                </div>

            </div>

            {/* Popular Spending */}
            <div className="px-6">
                <h3 className="text-base font-semibold tracking-wide text-white/90 mb-4">Popular Spending</h3>
                <div className="flex gap-3 overflow-x-auto pb-4 -mx-6 px-6 no-scrollbar">

                    <button onClick={() => alert('Filter: Food')} className="flex items-center gap-3 bg-[#1e1e24] border border-white/5 rounded-full px-4 py-2 whitespace-nowrap active:scale-[0.98] transition-transform shadow-lg cursor-pointer hover:bg-white/5">
                        <div className="w-8 h-8 rounded-full bg-[#f2ae9a] flex items-center justify-center -ml-2 text-xl">
                            🍕
                        </div>
                        <span className="font-semibold text-sm">Food</span>
                    </button>

                    <button onClick={() => alert('Filter: Household')} className="flex items-center gap-3 bg-[#1e1e24] border border-white/5 rounded-full px-4 py-2 whitespace-nowrap active:scale-[0.98] transition-transform shadow-lg cursor-pointer hover:bg-white/5">
                        <div className="w-8 h-8 rounded-full bg-[#9ba4ce] flex items-center justify-center -ml-2 text-xl">
                            🏠
                        </div>
                        <span className="font-semibold text-sm">Household</span>
                    </button>

                    <button onClick={() => alert('Filter: Study')} className="flex items-center gap-3 bg-[#1e1e24] border border-white/5 rounded-full px-4 py-2 whitespace-nowrap active:scale-[0.98] transition-transform shadow-lg cursor-pointer hover:bg-white/5">
                        <div className="w-8 h-8 rounded-full bg-[#6db392] flex items-center justify-center -ml-2 text-xl">
                            📖
                        </div>
                        <span className="font-semibold text-sm">Study</span>
                    </button>

                </div>
            </div>

        </div>
    );
}
