'use client';

import Image from 'next/image';
import Link from 'next/link';
import {
  ArrowLeft,
  Settings,
  Shield,
  CreditCard,
  LogOut,
  ChevronRight,
} from 'lucide-react';
import { BottomNav } from '@/components/BottomNav';

export default function ProfilePage() {
  const menuItems = [
    { icon: Shield, label: 'Security & Privacy', color: 'text-emerald-400', bg: 'bg-emerald-400/10' },
    { icon: CreditCard, label: 'Payment Methods', color: 'text-purple-400', bg: 'bg-purple-400/10' },
  ];

  return (
    <div className="app-screen overflow-y-auto no-scrollbar pb-24">
      
      {/* Header */}
      <div className="screen-header">
        <Link href="/" className="btn-icon-round">
          <ArrowLeft size={18} strokeWidth={2.5} />
        </Link>
        <h1 className="text-base font-semibold tracking-wide">
          Profile
        </h1>
        <button onClick={() => alert('Settings')} className="btn-icon-round">
          <Settings size={18} strokeWidth={2.5} />
        </button>
      </div>

      {/* Profile Info Section */}
      <div className="px-6 flex flex-col items-center mt-4 mb-8">
        <div className="relative">
          <div className="h-24 w-24 rounded-full bg-[#1a1a20] flex items-center justify-center p-1 border-2 border-emerald-500/50 overflow-hidden shadow-2xl">
            <Image
              src="https://api.dicebear.com/7.x/notionists/svg?seed=Daniel"
              alt="User Avatar"
              width={96}
              height={96}
              className="rounded-full bg-white/10"
            />
          </div>
          <div className="absolute bottom-1 right-1 h-6 w-6 bg-emerald-500 rounded-full border-4 border-[#121217] flex items-center justify-center">
            <div className="h-2 w-2 bg-white rounded-full animate-pulse"></div>
          </div>
        </div>
        
        <h2 className="text-2xl font-bold mt-4 tracking-tight">Daniel Jameson</h2>
        <p className="text-neutral-500 text-sm font-medium mt-1">daniel.design@icloud.com</p>
        
        <div className="mt-4 flex gap-2">
          <span className="bg-white/5 border border-white/10 px-3 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider text-emerald-400">
            Verified Pro
          </span>
        </div>
      </div>

      {/* Menu Sections */}
      <div className="px-6 space-y-3">
        <h3 className="text-xs font-bold uppercase tracking-widest text-neutral-500 ml-1 mb-2">Account Settings</h3>
        
        {menuItems.map((item, idx) => (
          <button 
            key={idx} 
            className="w-full flex items-center justify-between glass-card p-4 active:scale-[0.98] transition-transform"
            onClick={() => alert(item.label)}
          >
            <div className="flex items-center gap-4">
              <div className={`w-10 h-10 rounded-2xl ${item.bg} flex items-center justify-center ${item.color}`}>
                <item.icon size={20} />
              </div>
              <span className="text-sm font-semibold text-white/90">{item.label}</span>
            </div>
            <ChevronRight size={18} className="text-neutral-600" />
          </button>
        ))}

        {/* Logout Button */}
        <button 
          className="w-full mt-8 flex items-center justify-center gap-2 p-4 rounded-3xl bg-red-500/10 border border-red-500/20 text-red-500 font-bold text-sm hover:bg-red-500/20 transition-colors"
          onClick={() => alert('Logging out...')}
        >
          <LogOut size={18} />
          Sign Out
        </button>
      </div>

      <BottomNav />
    </div>
  );
}