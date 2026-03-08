'use client';

import Link from 'next/link';
import { useState } from 'react';
import { ArrowLeft, Delete, Shield } from 'lucide-react';
import { MiniKit } from '@worldcoin/minikit-js';
import shadowAbi from '@/abi/shadow.abi.json';

const SHADOW_ADDRESS = (process.env.NEXT_PUBLIC_SHADOW_CONTRACT ??
  '0x7aD3A7BeCe749eBaBEDF20fe2c688C525479e5De') as `0x${string}`;
const USDC_ADDRESS = (process.env.NEXT_PUBLIC_USDC_SEPOLIA ??
  '0x1c7D4B196Cb0C7B01d743Fbc6116a902379C7238') as `0x${string}`;

const ERC20_APPROVE_ABI = [
  {
    name: 'approve',
    type: 'function',
    stateMutability: 'nonpayable',
    inputs: [
      { name: 'spender', type: 'address' },
      { name: 'amount', type: 'uint256' },
    ],
    outputs: [{ name: '', type: 'bool' }],
  },
] as const;

type BtnState = 'idle' | 'pending' | 'success' | 'failed';

export default function TransferPage() {
  const [amount, setAmount] = useState('0');
  const [recipient, setRecipient] = useState('');
  const [btnState, setBtnState] = useState<BtnState>('idle');

  const handleNumberClick = (num: number | string) => {
    setAmount((prev) => {
      if (prev === '0') return num.toString();
      if (prev.length >= 7) return prev;
      return prev + num;
    });
  };

  const handleDelete = () => {
    setAmount((prev) => (prev.length <= 1 ? '0' : prev.slice(0, -1)));
  };

  const handleSend = async () => {
    if (amount === '0' || !amount) {
      alert('Enter an amount.');
      return;
    }
    if (!recipient || !/^0x[0-9a-fA-F]{40}$/.test(recipient)) {
      alert('Enter a valid recipient address (0x…).');
      return;
    }

    setBtnState('pending');
    try {
      const encRes = await fetch('/api/shielded-deposit', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ recipient, amount }),
      });
      if (!encRes.ok) throw new Error('Encryption failed');
      const { encryptedRecipient, amountRaw } = await encRes.json();

      const { finalPayload } = await MiniKit.commandsAsync.sendTransaction({
        transaction: [
          {
            address: USDC_ADDRESS,
            abi: ERC20_APPROVE_ABI,
            functionName: 'approve',
            args: [SHADOW_ADDRESS, amountRaw],
          },
          {
            address: SHADOW_ADDRESS,
            abi: shadowAbi,
            functionName: 'deposit',
            args: [encryptedRecipient, amountRaw],
          },
        ],
      });

      if (finalPayload.status === 'success') {
        setBtnState('success');
        setAmount('0');
        setRecipient('');
      } else {
        setBtnState('failed');
      }
    } catch (err) {
      console.error(err);
      setBtnState('failed');
    } finally {
      setTimeout(() => setBtnState('idle'), 3000);
    }
  };

  const formattedAmount = Number(amount).toLocaleString();

  return (
    <div className="app-screen">
      {/* Header */}
      <div className="screen-header">
        <Link href="/" className="btn-icon-round">
          <ArrowLeft size={18} strokeWidth={2.5} />
        </Link>
        <h1 className="text-base font-semibold tracking-wide flex items-center gap-2">
          <Shield size={14} className="text-emerald-400" />
          Shielded Transfer
        </h1>
        <div className="w-10" />
      </div>

      {/* Recipient Input */}
      <div className="px-6 mb-4">
        <div className="bg-[#1e1e24] border border-white/5 rounded-3xl p-4">
          <label className="label-mono mb-2 block">Recipient Address</label>
          <input
            type="text"
            placeholder="0x…"
            value={recipient}
            onChange={(e) => setRecipient(e.target.value)}
            className="w-full bg-transparent text-white text-sm font-mono placeholder-white/20 outline-none"
          />
        </div>
      </div>

      {/* Amount Display */}
      <div className="flex flex-col items-center justify-center py-6 h-36">
        <h2 className="text-6xl font-light tracking-tight text-white mb-4 shadow-sm flex items-center overflow-hidden max-w-full px-4">
          <span className="text-white/40 mr-1 text-5xl">$</span>
          <span className="truncate">{formattedAmount}</span>
          <span className="text-white/40">.00</span>
          <div className="w-0.5 h-12 bg-emerald-500 animate-pulse ml-1 opacity-80 shrink-0" />
        </h2>
        <p className="subtitle-sm font-medium">USDC · Ethereum Sepolia</p>
      </div>

      {/* Shielded note */}
      <div className="mx-6 mb-4 flex items-start gap-2 bg-emerald-500/10 border border-emerald-500/20 rounded-2xl p-3">
        <Shield size={13} className="text-emerald-400 mt-0.5 shrink-0" />
        <p className="text-[11px] text-emerald-300/80 leading-relaxed">
          Recipient is encrypted on-chain. CRE DON routes payout privately.
        </p>
      </div>

      {/* Numpad */}
      <div className="grid grid-cols-3 gap-2 px-6 mt-4 mb-24">
        {[1, 2, 3, 4, 5, 6, 7, 8, 9].map((num) => (
          <button
            key={num}
            onClick={() => handleNumberClick(num)}
            className="keypad-btn"
          >
            {num}
          </button>
        ))}
        <button className="h-16 flex items-center justify-center text-2xl font-medium text-white/30 rounded-2xl active:bg-white/5 transition-colors">
          .
        </button>
        <button onClick={() => handleNumberClick(0)} className="keypad-btn text-3xl">
          0
        </button>
        <button onClick={handleDelete} className="keypad-btn">
          <Delete size={20} strokeWidth={2} />
        </button>
      </div>

      {/* Send Button */}
      <div className="absolute bottom-6 left-0 right-0 px-6">
        <button
          onClick={handleSend}
          disabled={btnState === 'pending'}
          className={`w-full text-white rounded-full py-5 text-sm font-bold tracking-wide active:scale-[0.98] transition-all shadow-xl disabled:opacity-60 ${
            btnState === 'success'
              ? 'bg-emerald-600'
              : btnState === 'failed'
              ? 'bg-red-600'
              : 'bg-emerald-500/80 hover:bg-emerald-500'
          }`}
        >
          {btnState === 'pending'
            ? 'Sending…'
            : btnState === 'success'
            ? 'Sent!'
            : btnState === 'failed'
            ? 'Failed — try again'
            : 'Send Shielded'}
        </button>
      </div>
    </div>
  );
}
