import Link from "next/link";
import { User, Wallet, Home as HomeIcon, ArrowLeftRight } from "lucide-react";

export default function Dashboard() {
    return (
        <div className="flex h-[100dvh] flex-col bg-[#121212] tracking-wide text-white overflow-hidden relative max-w-md mx-auto sm:border sm:border-neutral-800 pb-[80px]">
            <main className="flex w-full grow flex-col overflow-y-auto p-6 scrollbar-hide">
                {/* Header */}
                <header className="flex items-center justify-between mb-8">
                    <div className="flex items-center gap-3">
                        <div className="h-10 w-10 rounded-full bg-neutral-700 overflow-hidden flex items-center justify-center">
                            <span className="text-xl">🧑🏽‍🦱</span>
                        </div>
                        <div>
                            <p className="text-xs text-neutral-400">Hi, Arip!</p>
                            <h1 className="text-lg font-medium">Monthly Budget</h1>
                        </div>
                    </div>
                    <div className="flex items-center gap-2 rounded-full bg-[#1b1b22] px-3 py-1.5 border border-white/5">
                        <div className="h-2 w-2 rounded-full bg-emerald-400"></div>
                        <span className="text-xs font-medium text-neutral-300">My Balance</span>
                    </div>
                </header>

                {/* Planned Expenses & Chart */}
                <section className="flex justify-between items-start mb-8">
                    <div>
                        <h2 className="text-sm font-medium text-neutral-300 mb-1">Planned Expenses</h2>
                        <p className="text-3xl font-semibold mb-2">$1,475<span className="text-xl text-neutral-400">.00</span></p>
                        <div className="inline-block rounded-full bg-[#1e231e] px-2.5 py-1">
                            <span className="text-[10px] font-medium text-emerald-400">$32 Left to budget</span>
                        </div>
                    </div>

                    {/* SVG Pie Chart for cross-browser compatibility */}
                    <div className="relative h-16 w-16 mr-2 flex-shrink-0">
                        <svg viewBox="0 0 36 36" className="w-full h-full transform -rotate-90 drop-shadow-lg">
                            {/* background ring */}
                            <circle cx="18" cy="18" r="16" fill="none" className="stroke-white/10" strokeWidth="4" />
                            {/* Purple Segment (61%) */}
                            <circle cx="18" cy="18" r="16" fill="none" className="stroke-[#8b7df0]" strokeWidth="4" strokeDasharray="61 100" strokeDashoffset="0" />
                            {/* Blue Segment (19%) gap of 2 then 19 */}
                            <circle cx="18" cy="18" r="16" fill="none" className="stroke-[#6ea0f5]" strokeWidth="4" strokeDasharray="19 100" strokeDashoffset="-63" />
                        </svg>
                        <div className="absolute top-1/2 left-1/2 w-4 h-4 bg-white rounded-full transform -translate-x-1/2 -translate-y-1/2 opacity-90 shadow-md"></div>
                    </div>
                </section>

                {/* Categories Cards */}
                <section className="flex gap-3 overflow-x-auto pb-6 mb-2 scrollbar-hide -mx-6 px-6 snap-x">
                    <div className="flex h-[130px] min-w-[60px] snap-center items-center justify-center rounded-[28px] bg-[#1b1b22] border border-white/5 opacity-60">
                        <span className="text-xl">+</span>
                    </div>

                    <div className="flex flex-col h-[130px] min-w-[120px] snap-center rounded-[28px] bg-[#8b7df0] p-4 text-white relative overflow-hidden shadow-lg shadow-[#8b7df0]/30">
                        <h3 className="text-sm font-medium mb-1 relative z-10">Housing</h3>
                        <p className="text-xl font-semibold mb-auto relative z-10">$965<span className="text-sm opacity-70">.00</span></p>
                        <div className="inline-flex w-fit items-center justify-center rounded-full bg-white/20 px-2 py-0.5 backdrop-blur-sm relative z-10 mt-2">
                            <span className="text-[10px] font-medium">61%</span>
                        </div>
                        <div className="absolute -bottom-4 -right-4 h-20 w-20 rounded-full bg-white/10 blur-xl"></div>
                    </div>

                    <div className="flex flex-col h-[130px] min-w-[120px] snap-center rounded-[28px] bg-[#a8e6cf] p-4 text-[#121212] relative overflow-hidden">
                        <h3 className="text-sm font-medium mb-1">Food</h3>
                        <p className="text-xl font-semibold mb-auto">$300<span className="text-sm opacity-70">.00</span></p>
                        <div className="inline-flex w-fit items-center justify-center rounded-full bg-black/10 px-2 py-0.5 mt-2">
                            <span className="text-[10px] font-medium">19%</span>
                        </div>
                    </div>

                    <div className="flex flex-col h-[130px] min-w-[120px] snap-center rounded-[28px] bg-white p-4 text-[#121212] relative overflow-hidden shadow-[0_4px_20px_rgba(255,255,255,0.15)] origin-bottom z-10 scale-[1.02]">
                        <h3 className="text-sm font-medium mb-1">Saving</h3>
                        <p className="text-xl font-semibold mb-auto">$200<span className="text-sm opacity-70">.00</span></p>
                        <div className="inline-flex w-fit items-center justify-center rounded-full bg-black/5 px-2 py-0.5 mt-2">
                            <span className="text-[10px] font-medium text-neutral-600">13%</span>
                        </div>
                    </div>
                </section>

                {/* My Income Section */}
                <section className="mb-8">
                    <h2 className="text-base font-medium mb-4">My Income</h2>
                    <div className="flex gap-4">
                        <div className="flex-1 rounded-3xl bg-[#1a1a20] p-4 border border-white/5">
                            <div className="flex justify-between items-start mb-6">
                                <div className="h-8 w-8 rounded-full bg-[#2a2a32] flex items-center justify-center">
                                    <span className="font-medium text-sm">$</span>
                                </div>
                                <span className="text-neutral-500 text-xs tracking-widest">•••</span>
                            </div>
                            <p className="text-xs text-neutral-400 mb-1">Salary</p>
                            <p className="text-lg font-medium">$1.500<span className="text-sm text-neutral-500">,00</span></p>
                        </div>

                        <div className="flex-1 rounded-3xl bg-[#1a1a20] p-4 border border-white/5">
                            <div className="flex justify-between items-start mb-6">
                                <div className="h-8 w-8 rounded-full bg-[#2a2a32] flex items-center justify-center">
                                    <Wallet size={14} className="text-white" />
                                </div>
                                <span className="text-neutral-500 text-xs tracking-widest">•••</span>
                            </div>
                            <p className="text-xs text-neutral-400 mb-1">Interest</p>
                            <p className="text-lg font-medium">$240<span className="text-sm text-neutral-500">,00</span></p>
                        </div>
                    </div>
                </section>

                {/* Spending March Section - Preview */}
                <section>
                    <div className="flex justify-between items-end mb-4">
                        <h2 className="text-base font-medium">Spending March</h2>
                        <Link href="/transactions" className="text-xs text-[#8b7df0] font-medium hover:underline">See All</Link>
                    </div>

                    <div className="flex items-center justify-between mb-4">
                        <div className="flex items-center gap-3">
                            <div className="h-10 w-10 flex items-center justify-center rounded-full bg-[#1b1b22] border border-white/5">
                                <div className="text-xl">🎵</div>
                            </div>
                            <div>
                                <p className="font-medium text-sm">Spotify</p>
                                <p className="text-[10px] text-neutral-500 mt-0.5">10:00 am • Mar 26th 2023</p>
                            </div>
                        </div>
                        <p className="font-medium text-sm">$54.99</p>
                    </div>

                    <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                            <div className="h-10 w-10 flex items-center justify-center rounded-full bg-[#1b1b22] border border-white/5">
                                <div className="text-lg">🎨</div>
                            </div>
                            <div>
                                <p className="font-medium text-sm">Figma</p>
                                <p className="text-[10px] text-neutral-500 mt-0.5">8:00 am • Mar 21th 2023</p>
                            </div>
                        </div>
                        <p className="font-medium text-sm">$8.99</p>
                    </div>
                </section>
            </main>

            {/* Floating Bottom Nav */}
            <div className="absolute bottom-6 left-0 right-0 z-50 flex justify-center w-full px-6">
                <nav className="flex items-center gap-8 rounded-full bg-[#1b1b22]/90 backdrop-blur-md px-8 py-3.5 shadow-xl border border-white/5">
                    <Link href="/dashboard" className="h-10 w-10 flex items-center justify-center rounded-full bg-[#8b7df0] text-white shadow-[0_0_15px_rgba(139,125,240,0.5)]">
                        <HomeIcon size={20} />
                    </Link>
                    <Link href="/transactions" className="h-10 w-10 flex items-center justify-center text-neutral-500 hover:text-white transition-colors">
                        <ArrowLeftRight size={20} />
                    </Link>
                    <button className="h-10 w-10 flex items-center justify-center text-neutral-500 hover:text-white transition-colors">
                        <User size={20} />
                    </button>
                </nav>
            </div>
        </div>
    );
}
