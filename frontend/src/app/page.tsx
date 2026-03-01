import Image from 'next/image';
import Link from 'next/link';
import { ArrowRight, Lock, ChevronRight, LayoutDashboard, ArrowLeftRight } from 'lucide-react';

export default function Home() {
  return (
    <div className="flex h-[100dvh] flex-col items-center justify-between bg-[#121212] tracking-wide text-white overflow-hidden relative max-w-md mx-auto sm:border sm:border-neutral-800">

      {/* Background Graphic Lines Simulation */}
      <div className="absolute top-0 left-0 w-full h-1/2 overflow-hidden pointer-events-none opacity-20">
        <svg viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg" className="absolute -top-10 -left-10 w-full animate-pulse blur-[1px]">
          <path fill="none" stroke="#6C63FF" strokeWidth="0.5" d="M0,100 C50,20 150,180 200,100" />
          <path fill="none" stroke="#fff" strokeWidth="0.2" d="M0,120 C80,30 120,200 200,80" />
        </svg>
      </div>

      {/* Modern Floating Slidebar / Navigation */}
      <div className="absolute top-8 left-0 right-0 z-50 flex justify-center w-full px-8">
        <div className="flex items-center rounded-full bg-[#1b1b22]/90 backdrop-blur-xl p-1 shadow-2xl border border-white/10 mx-auto">
          <Link
            href="/dashboard"
            className="flex items-center gap-2 text-[11px] font-medium tracking-wider uppercase text-white hover:bg-white/10 transition-colors px-5 py-2.5 rounded-full"
          >
            <LayoutDashboard size={13} className="text-[#8b7df0]" />
            <span>Dashboard</span>
          </Link>

          <div className="w-[1px] h-4 bg-white/10 mx-1"></div>

          <Link
            href="/transactions"
            className="flex items-center gap-2 text-[11px] font-medium tracking-wider uppercase text-neutral-400 hover:text-white hover:bg-white/10 transition-colors px-5 py-2.5 rounded-full"
          >
            <ArrowLeftRight size={13} className="text-[#6ea0f5]" />
            <span>Transactions</span>
          </Link>
        </div>
      </div>

      <main className="flex w-full grow flex-col items-center justify-center p-6 z-10 animate-in fade-in slide-in-from-bottom-5 duration-1000">

        {/* Illustration Section */}
        <div className="relative mb-12 flex items-center justify-center w-full max-w-[280px] aspect-square">
          <Image
            src="/illustration.png"
            alt="Finance Illustration"
            width={280}
            height={280}
            className="object-contain drop-shadow-2xl"
            priority
          />
        </div>

        {/* Text Section */}
        <div className="w-full text-left space-y-4 max-w-[320px]">
          <h1 className="text-3xl font-semibold leading-snug">
            Take Control of Your<br />Finances Today!
          </h1>
          <p className="text-[#a1a1aa] text-sm leading-relaxed">
            With our app, you can easily track your income and
            expenses, set financial goals, and make informed
            decisions about your money.
          </p>
        </div>
      </main>

      {/* Bottom Actions */}
      <footer className="w-full p-6 pt-0 z-10 w-full mt-auto">
        <button
          className="flex w-full items-center justify-between rounded-[40px] bg-[#1a1a20] p-3 transition-opacity hover:opacity-90 hover:bg-[#202028] shadow-lg border border-white/5 group"
        >
          {/* Arrow Circle */}
          <div className="flex h-[52px] w-[52px] items-center justify-center rounded-full bg-[#8b7df0] text-white shadow-[0_0_15px_rgba(139,125,240,0.5)] group-hover:scale-105 transition-transform">
            <ArrowRight size={24} />
          </div>

          <div className="flex flex-1 items-center px-4 font-medium tracking-wide">
            Get Started
          </div>

          {/* Faint Chevrons */}
          <div className="flex text-[#8b7df0] opacity-40 -space-x-2 mr-3 group-hover:opacity-80 transition-opacity">
            <ChevronRight size={20} />
            <ChevronRight size={20} />
            <ChevronRight size={20} />
          </div>

          {/* Lock Icon */}
          <div className="flex h-12 w-12 items-center justify-center rounded-full bg-[#2a2a32]">
            <Lock size={18} className="text-neutral-400" />
          </div>
        </button>
      </footer>
    </div>
  );
}
