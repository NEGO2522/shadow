'use client';

import { useState } from 'react';
import { Shield, Send, ArrowRight } from 'lucide-react';
import { MiniKit, ResponseEvent } from '@worldcoin/minikit-js';
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

type PayMode = 'standard' | 'shielded';
type BtnState = 'pending' | 'success' | 'failed' | undefined;

export const Pay = () => {
  const [mode, setMode] = useState<PayMode>('standard');
  const [btnState, setBtnState] = useState<BtnState>(undefined);
  const [recipient, setRecipient] = useState('');
  const [amount, setAmount] = useState('');

  // ── Standard MiniKit payment ──────────────────────────────────────
  const handleStandardPay = async () => {
    if (!recipient || !amount) {
      alert('Enter recipient address and amount.');
      return;
    }
    setBtnState('pending');
    try {
      const res = await fetch('/api/initiate-payment', { method: 'POST' });
      const { id } = await res.json();

      const { finalPayload } = await MiniKit.commandsAsync.pay({
        reference: id,
        to: recipient as `0x${string}`,
        tokens: [
          {
            symbol: 'USDC' as const,
            token_amount: String(Math.round(parseFloat(amount) * 1_000_000)),
          },
        ],
        description: 'Shadow payment',
      });

      setBtnState(finalPayload.status === 'success' ? 'success' : 'failed');
    } catch {
      setBtnState('failed');
    } finally {
      setTimeout(() => setBtnState(undefined), 3000);
    }
  };

  // ── Shielded payment via Shadow contract (CRE workflow) ───────────
  const handleShieldedPay = async () => {
    if (!recipient || !amount) {
      alert('Enter recipient address and amount.');
      return;
    }
    setBtnState('pending');
    try {
      // Step 1: Encrypt recipient server-side
      const encRes = await fetch('/api/shielded-deposit', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ recipient, amount }),
      });
      if (!encRes.ok) throw new Error('Encryption failed');
      const { encryptedRecipient, amountRaw } = await encRes.json();

      // Step 2: USDC approve + Shadow.deposit in one MiniKit transaction
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
        MiniKit.subscribe(ResponseEvent.MiniAppSendTransaction, () => {});
      } else {
        setBtnState('failed');
      }
    } catch (err) {
      console.error(err);
      setBtnState('failed');
    } finally {
      setTimeout(() => setBtnState(undefined), 3000);
    }
  };

  const stateLabel: Record<NonNullable<BtnState>, string> = {
    pending: 'Processing…',
    success: 'Sent!',
    failed: 'Failed — try again',
  };

  return (
    <div className="w-full space-y-4">
      <p className="text-lg font-semibold">Pay</p>

      {/* Mode toggle */}
      <div className="flex gap-2 p-1 bg-[#1a1b1f] rounded-2xl border border-white/5">
        {(['standard', 'shielded'] as PayMode[]).map((m) => (
          <button
            key={m}
            onClick={() => setMode(m)}
            className={`flex-1 flex items-center justify-center gap-1.5 py-2.5 rounded-xl text-xs font-bold transition-all ${
              mode === m
                ? m === 'shielded'
                  ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'
                  : 'bg-white text-black'
                : 'text-neutral-400 hover:text-white'
            }`}
          >
            {m === 'shielded' && <Shield size={12} />}
            {m.charAt(0).toUpperCase() + m.slice(1)}
          </button>
        ))}
      </div>

      {/* Form fields */}
      <div className="space-y-3">
        <div className="flex flex-col gap-1.5">
          <label className="label-mono">
            {mode === 'shielded' ? 'Recipient Address (private)' : 'Recipient Address'}
          </label>
          <input
            type="text"
            placeholder="0x…"
            value={recipient}
            onChange={(e) => setRecipient(e.target.value)}
            className="pay-input"
          />
        </div>
        <div className="flex flex-col gap-1.5">
          <label className="label-mono">Amount (USDC)</label>
          <input
            type="number"
            placeholder="0.00"
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            className="pay-input"
          />
        </div>
      </div>

      {mode === 'shielded' && (
        <div className="flex items-start gap-2 bg-emerald-500/10 border border-emerald-500/20 rounded-2xl p-3">
          <Shield size={14} className="text-emerald-400 mt-0.5 shrink-0" />
          <p className="text-[11px] text-emerald-300/80 leading-relaxed">
            Recipient is encrypted on-chain. The CRE DON decrypts and routes
            the payout privately.
          </p>
        </div>
      )}

      {/* CTA */}
      <button
        onClick={mode === 'shielded' ? handleShieldedPay : handleStandardPay}
        disabled={btnState === 'pending'}
        className={`w-full flex items-center justify-center gap-2 py-4 rounded-[28px] text-sm font-bold transition-all active:scale-95 disabled:opacity-60 ${
          mode === 'shielded'
            ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30 hover:bg-emerald-500/30'
            : 'btn-action-primary'
        }`}
      >
        {mode === 'shielded' ? <Shield size={15} /> : <Send size={15} />}
        {btnState ? stateLabel[btnState] : mode === 'shielded' ? 'Shielded Send' : 'Pay'}
        {!btnState && <ArrowRight size={14} />}
      </button>
    </div>
  );
};
