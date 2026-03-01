import Link from 'next/link';
import { X } from 'lucide-react';

export default function Transactions() {
    return (
        <div className="flex h-[100dvh] flex-col bg-[#121212] tracking-wide text-white overflow-hidden relative max-w-md mx-auto sm:border sm:border-neutral-800">

            {/* Header and Balance Details */}
            <header className="p-6 pb-2 relative">
                <Link href="/dashboard" className="absolute top-6 right-6 h-8 w-8 rounded-full border border-white/10 flex items-center justify-center text-neutral-400 hover:text-white transition-colors">
                    <X size={16} />
                </Link>
                <p className="text-sm font-medium text-neutral-400 mb-1">My Balance</p>
                <h1 className="text-4xl font-semibold mb-3 tracking-tight">$8,822<span className="text-2xl text-neutral-400">.89</span></h1>

                <div className="inline-flex items-center gap-1 rounded-full bg-emerald-500/20 text-emerald-400 px-2.5 py-1 text-xs font-medium border border-emerald-500/10 mb-8 shadow-[0_0_10px_rgba(16,185,129,0.1)]">
                    <span className="text-[10px]">↗</span>
                    +08%
                </div>

                {/* Progress Bars */}
                <div className="flex justify-between items-center mb-3">
                    <h2 className="text-sm font-medium">Spending March</h2>
                    <span className="text-sm font-medium text-neutral-400">$744<span className="text-xs">.97</span></span>
                </div>

                <div className="flex h-2.5 w-full gap-2 rounded-full overflow-hidden bg-[#1a1a20]">
                    <div className="h-full bg-[#8b7df0] w-[45%] rounded-full shadow-[0_0_8px_rgba(139,125,240,0.6)]"></div>
                    <div className="h-full bg-[#6ea0f5] w-[25%] rounded-full opacity-90"></div>
                    <div className="h-full bg-white w-[20%] rounded-full shadow-[0_0_10px_rgba(255,255,255,0.8)] z-10"></div>
                    <div className="h-full w-[10%]"></div>
                </div>
            </header>

            {/* Transaction Lists */}
            <main className="flex-1 overflow-y-auto px-6 pt-4 pb-12 scrollbar-hide">
                {/* This Month */}
                <div className="mb-8">
                    <h3 className="text-xs font-medium text-neutral-500 mb-4 sticky top-0 bg-[#121212]/90 backdrop-blur py-2 z-10 w-full">This Month</h3>
                    <div className="space-y-6">
                        <div className="flex items-center justify-between group cursor-pointer">
                            <div className="flex items-center gap-4">
                                <div className="h-12 w-12 flex items-center justify-center rounded-full bg-[#1b1b22] border border-white/5 group-hover:bg-[#202028] transition-colors">
                                    <span className="text-xl opacity-90">🎵</span>
                                </div>
                                <div>
                                    <p className="font-medium text-[15px] mb-0.5 group-hover:text-[#8b7df0] transition-colors">Spotify</p>
                                    <p className="text-[11px] text-neutral-500">10:00 am • Mar 26th 2023</p>
                                </div>
                            </div>
                            <p className="font-medium text-[15px]">$54.99</p>
                        </div>

                        <div className="flex items-center justify-between group cursor-pointer">
                            <div className="flex items-center gap-4">
                                <div className="h-12 w-12 flex items-center justify-center rounded-full bg-[#1b1b22] border border-white/5 group-hover:bg-[#202028] transition-colors">
                                    <span className="text-xl opacity-90">🎨</span>
                                </div>
                                <div>
                                    <p className="font-medium text-[15px] mb-0.5 group-hover:text-[#8b7df0] transition-colors">Figma</p>
                                    <p className="text-[11px] text-neutral-500">8:00 am • Mar 21th 2023</p>
                                </div>
                            </div>
                            <p className="font-medium text-[15px]">$8.99</p>
                        </div>

                        <div className="flex items-center justify-between group cursor-pointer">
                            <div className="flex items-center gap-4">
                                <div className="h-12 w-12 flex items-center justify-center rounded-full bg-[#1b1b22] border border-white/5 group-hover:bg-[#202028] transition-colors">
                                    <span className="text-xl opacity-90">🛒</span>
                                </div>
                                <div>
                                    <p className="font-medium text-[15px] mb-0.5 group-hover:text-[#8b7df0] transition-colors">Online Shopping</p>
                                    <p className="text-[11px] text-neutral-500">10:00 am • Mar 11th 2023</p>
                                </div>
                            </div>
                            <p className="font-medium text-[15px]">$132.00</p>
                        </div>

                        <div className="flex items-center justify-between group cursor-pointer">
                            <div className="flex items-center gap-4">
                                <div className="h-12 w-12 flex items-center justify-center rounded-full bg-[#1b1b22] border border-white/5 group-hover:bg-[#202028] transition-colors">
                                    <span className="text-xl opacity-90">🏠</span>
                                </div>
                                <div>
                                    <p className="font-medium text-[15px] mb-0.5 group-hover:text-[#8b7df0] transition-colors">AirBnB rent</p>
                                    <p className="text-[11px] text-neutral-500">11:00 am • Mar 2th 2023</p>
                                </div>
                            </div>
                            <p className="font-medium text-[15px]">$548.99</p>
                        </div>
                    </div>
                </div>

                {/* Last Month */}
                <div>
                    <h3 className="text-xs font-medium text-neutral-500 mb-4 sticky top-0 bg-[#121212]/90 backdrop-blur py-2 z-10 w-full">Last Month</h3>
                    <div className="space-y-6">
                        <div className="flex items-center justify-between group cursor-pointer">
                            <div className="flex items-center gap-4">
                                <div className="h-12 w-12 flex items-center justify-center rounded-full bg-[#1b1b22] border border-white/5 group-hover:bg-[#202028] transition-colors">
                                    <span className="text-xl opacity-90">🎵</span>
                                </div>
                                <div>
                                    <p className="font-medium text-[15px] mb-0.5 group-hover:text-[#8b7df0] transition-colors">Spotify</p>
                                    <p className="text-[11px] text-neutral-500">10:00 am • Feb 16th 2023</p>
                                </div>
                            </div>
                            <p className="font-medium text-[15px]">$54.99</p>
                        </div>

                        <div className="flex items-center justify-between group cursor-pointer opacity-50">
                            <div className="flex items-center gap-4">
                                <div className="h-12 w-12 flex items-center justify-center rounded-full bg-[#1b1b22] border border-white/5 group-hover:bg-[#202028] transition-colors">
                                    <span className="text-xl opacity-90">🎨</span>
                                </div>
                                <div>
                                    <p className="font-medium text-[15px] mb-0.5 group-hover:text-[#8b7df0] transition-colors">Figma</p>
                                    <p className="text-[11px] text-neutral-500">8:00 am • Feb 12th 2023</p>
                                </div>
                            </div>
                            <p className="font-medium text-[15px]">$8.99</p>
                        </div>
                    </div>
                </div>
            </main>
        </div>
    );
}
