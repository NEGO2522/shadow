'use client';

import Link from 'next/link';
import Image from 'next/image';
import { useState } from 'react';
import { ArrowLeft, ChevronRight, ChevronDown, Delete } from 'lucide-react';

export default function TransferPage() {
    const [amount, setAmount] = useState('0');
    const [showCategories, setShowCategories] = useState(false);

    const categories = [
        { id: 'food', name: 'Food', icon: '🍕', color: 'bg-[#e8a89b]' },
        { id: 'home', name: 'Household', icon: '🏠', color: 'bg-[#9ba4ce]' },
        { id: 'study', name: 'Study', icon: '📖', color: 'bg-[#6db392]' },
        { id: 'travel', name: 'Travel', icon: '✈️', color: 'bg-[#d1ab75]' },
        { id: 'bills', name: 'Bills', icon: '🧾', color: 'bg-[#b69bce]' },
        { id: 'health', name: 'Health', icon: '💊', color: 'bg-[#9bceaf]' }
    ];
    const [selectedCategory, setSelectedCategory] = useState(categories[0]);

    const handleNumberClick = (num: number | string) => {
        setAmount(prev => {
            if (prev === '0') return num.toString();
            if (prev.length >= 7) return prev; // Prevent too many digits
            return prev + num;
        });
    };

    const handleDelete = () => {
        setAmount(prev => {
            if (prev.length <= 1) return '0';
            return prev.slice(0, -1);
        });
    };

    const handleSend = () => {
        if (amount === '0') {
            alert('Please enter an amount to send.');
            return;
        }
        alert(`Successfully sent ₹${Number(amount).toLocaleString()}.00 to Selma Knight for ${selectedCategory.name}!`);
        setAmount('0');
    };

    const formattedAmount = Number(amount).toLocaleString();

    return (
        <div className="flex h-[100dvh] flex-col bg-[#141416] tracking-wide text-white relative mx-auto sm:border sm:border-neutral-800 max-w-md font-sans">

            {/* Header */}
            <div className="flex items-center justify-between px-6 pt-12 pb-6">
                <Link href="/" className="h-10 w-10 flex items-center justify-center rounded-full bg-[#1e1e24] text-gray-300 relative border border-white/5 active:scale-95 transition-transform">
                    <ArrowLeft size={18} strokeWidth={2.5} />
                </Link>
                <h1 className="text-base font-semibold tracking-wide">
                    Transfer
                </h1>
                <div className="w-10"></div> {/* Spacer for centering */}
            </div>

            {/* Amount Display */}
            <div className="flex flex-col items-center justify-center py-8 h-40">
                <h2 className="text-6xl font-light tracking-tight text-white mb-4 shadow-sm flex items-center overflow-hidden max-w-full px-4">
                    <span className="text-white/40 mr-1 text-5xl">₹</span>
                    <span className="truncate">{formattedAmount}</span><span className="text-white/40">.00</span>
                    <div className="w-0.5 h-12 bg-indigo-500 animate-pulse ml-1 opacity-80 shrink-0"></div>
                </h2>
                <p className="text-sm font-medium text-white/60 tracking-wide">
                    Your Available Balance (₹ 20,456.00)
                </p>
            </div>

            {/* Target & Category Selection */}
            <div className="px-6 space-y-3 mt-4">

                {/* Transfer Target Card */}
                <button onClick={() => alert('Select recipient')} className="flex w-full items-center justify-between bg-[#5b6a84] p-4 rounded-3xl active:scale-[0.98] transition-transform shadow-lg relative overflow-hidden group">
                    <div className="absolute top-0 right-0 w-32 h-32 bg-white/5 rounded-full blur-2xl -mr-10 -mt-10"></div>

                    <div className="flex items-center gap-4 relative z-10">
                        <div className="w-12 h-12 rounded-full bg-[#3d4860] flex items-center justify-center overflow-hidden border border-white/10 shadow-inner">
                            <Image src="https://api.dicebear.com/7.x/notionists/svg?seed=Selma" alt="Selma" width={48} height={48} className="bg-white/5" />
                        </div>
                        <div className="text-left">
                            <p className="font-semibold text-white/90 text-[15px] mb-0.5">Transfer to Selma Knight</p>
                            <p className="font-mono text-xs text-white/50 tracking-widest">**** **** 9946</p>
                        </div>
                    </div>
                    <div className="relative z-10 text-white/40 group-hover:text-white transition-colors">
                        <ChevronRight size={20} strokeWidth={2.5} />
                    </div>
                </button>

                {/* Category Card Menu */}
                <div className="relative">
                    <button onClick={() => setShowCategories(!showCategories)} className="flex w-full items-center justify-between bg-[#2a2a32] p-4 rounded-3xl border border-white/5 active:scale-[0.98] transition-all shadow-lg group">
                        <div className="flex items-center gap-4">
                            <div className={`w-12 h-12 rounded-full ${selectedCategory.color} flex items-center justify-center text-2xl shadow-inner`}>
                                {selectedCategory.icon}
                            </div>
                            <div className="text-left">
                                <p className="font-mono text-[10px] uppercase tracking-widest text-[#a1a1aa] font-bold mb-1 opacity-80">Category</p>
                                <p className="font-semibold text-white/90 text-[15px]">{selectedCategory.name}</p>
                            </div>
                        </div>
                        <div className={`text-white/40 group-hover:text-white transition-transform duration-300 ${showCategories ? 'rotate-180' : ''}`}>
                            <ChevronDown size={20} strokeWidth={2.5} />
                        </div>
                    </button>

                    {/* Expandable Categories Dropdown */}
                    <div className={`absolute top-[105%] left-0 right-0 z-50 bg-[#1e1e24] border border-white/10 rounded-2xl shadow-2xl p-2 grid grid-cols-2 gap-2 transition-all duration-300 origin-top overflow-hidden ${showCategories ? 'scale-y-100 opacity-100 visible' : 'scale-y-95 opacity-0 invisible'}`}>
                        {categories.map(cat => (
                            <button
                                key={cat.id}
                                onClick={() => { setSelectedCategory(cat); setShowCategories(false); }}
                                className={`flex items-center gap-3 p-3 rounded-xl transition-colors ${selectedCategory.id === cat.id ? 'bg-white/10 border border-white/10' : 'hover:bg-white/5 border border-transparent'}`}
                            >
                                <div className={`w-8 h-8 rounded-full ${cat.color} flex items-center justify-center text-lg`}>{cat.icon}</div>
                                <span className={`font-medium text-sm ${selectedCategory.id === cat.id ? 'text-white' : 'text-white/70'}`}>{cat.name}</span>
                            </button>
                        ))}
                    </div>
                </div>

            </div>

            {/* Numpad */}
            <div className="grid grid-cols-3 gap-2 px-6 mt-10 flex-1 content-start mb-24">
                {[1, 2, 3, 4, 5, 6, 7, 8, 9].map((num) => (
                    <button key={num} onClick={() => handleNumberClick(num)} className="h-16 flex items-center justify-center text-2xl font-light text-white rounded-2xl bg-[#1e1e24] active:bg-white/10 transition-colors border border-white/5">
                        {num}
                    </button>
                ))}
                <button onClick={() => alert('Decimals not supported in this demo')} className="h-16 flex items-center justify-center text-2xl font-medium text-white/30 rounded-2xl active:bg-white/5 transition-colors">
                    .
                </button>
                <button onClick={() => handleNumberClick(0)} className="h-16 flex items-center justify-center text-3xl font-light text-white rounded-2xl bg-[#1e1e24] active:bg-white/10 transition-colors border border-white/5">
                    0
                </button>
                <button onClick={handleDelete} className="h-16 flex items-center justify-center text-white/70 rounded-2xl bg-[#1e1e24] active:bg-white/10 transition-colors border border-white/5">
                    <Delete size={20} strokeWidth={2} />
                </button>
            </div>

            {/* Send Money Button */}
            <div className="absolute bottom-6 left-0 right-0 px-6">
                <button onClick={handleSend} className="w-full bg-[#6a73a3] text-white rounded-full py-4 text-sm font-bold tracking-wide active:scale-[0.98] transition-all shadow-xl shadow-indigo-500/10 hover:bg-[#7882b8]">
                    Send Money
                </button>
            </div>

        </div>
    );
}
