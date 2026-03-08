'use client';

import { useState, useEffect } from 'react';
import { IDKitRequestWidget, orbLegacy, type RpContext, type IDKitResult } from '@worldcoin/idkit';
import { MiniKit } from '@worldcoin/minikit-js';
import { Shield, QrCode } from 'lucide-react';
import Dashboard from '@/components/Dashboard';

const ACTION = process.env.NEXT_PUBLIC_WLD_ACTION_NAME ?? 'payment-verification';
const APP_ID = (process.env.NEXT_PUBLIC_APP_ID ?? '') as `app_${string}`;
const RP_ID  = process.env.NEXT_PUBLIC_WLD_RP_ID ?? '';
const STORAGE = 'shadow_wld_verified';
const WALLET_STORAGE = 'shadow_wld_wallet';

export default function App() {
  const [verified, setVerified]   = useState<boolean | null>(null);
  const [rpContext, setRpContext]  = useState<RpContext | null>(null);
  const [open, setOpen]           = useState(false);
  const [fetching, setFetching]   = useState(false);
  const [error, setError]         = useState('');

  useEffect(() => {
    setVerified(localStorage.getItem(STORAGE) === 'true');
  }, []);

  // Step 1 — fetch RP signature from our backend, then open the widget
  const handleVerify = async () => {
    setError('');
    setFetching(true);
    try {
      const rpSig = await fetch('/api/rp-signature', {
        method: 'POST',
        headers: { 'content-type': 'application/json' },
        body: JSON.stringify({ action: ACTION }),
      }).then((r) => r.json());

      setRpContext({
        rp_id:      RP_ID,
        nonce:      rpSig.nonce,
        created_at: rpSig.created_at,
        expires_at: rpSig.expires_at,
        signature:  rpSig.sig,
      });
      setOpen(true);
    } catch {
      setError('Could not reach server. Check your env vars.');
    } finally {
      setFetching(false);
    }
  };

  // Step 2 — widget calls this after proof; we verify on the backend
  const handleWidgetVerify = async (result: IDKitResult) => {
    const response = await fetch('/api/verify-proof', {
      method: 'POST',
      headers: { 'content-type': 'application/json' },
      body: JSON.stringify({ rp_id: RP_ID, idkitResponse: result }),
    });
    if (!response.ok) throw new Error('Backend verification failed');
  };

  // Step 3 — handleVerify resolved → mark verified + store wallet address
  const onSuccess = () => {
    localStorage.setItem(STORAGE, 'true');
    const addr = MiniKit.walletAddress;
    if (addr) localStorage.setItem(WALLET_STORAGE, addr);
    setVerified(true);
  };

  const onError = (code: string) => {
    setError(`Verification error: ${code}`);
  };

  // ── Loading ───────────────────────────────────────────────────
  if (verified === null) {
    return (
      <div className="app-screen items-center justify-center">
        <div className="w-5 h-5 border-2 border-white/10 border-t-white/60 rounded-full animate-spin" />
      </div>
    );
  }

  if (verified) return <Dashboard />;

  // ── Verify gate ───────────────────────────────────────────────
  return (
    <div className="app-screen items-center justify-center relative overflow-hidden">
      <div className="absolute inset-0 flex items-center justify-center pointer-events-none">
        <div className="w-80 h-80 bg-white/[0.025] rounded-full blur-3xl" />
      </div>

      <div className="relative z-10 flex flex-col items-center px-8 w-full">
        <div className="w-14 h-14 rounded-full bg-white flex items-center justify-center mb-10 shadow-[0_0_48px_rgba(255,255,255,0.12)]">
          <Shield size={24} className="text-black" strokeWidth={2.5} />
        </div>

        <h1 className="text-2xl font-semibold tracking-tight text-white mb-2 text-center">
          Shadow
        </h1>
        <p className="text-white/30 text-[13px] text-center mb-14 tracking-wide leading-relaxed">
          Prove you&apos;re human to continue
        </p>

        {/* IDKit widget — only renders when open=true and rpContext is ready */}
        {rpContext && (
          <IDKitRequestWidget
            open={open}
            onOpenChange={setOpen}
            app_id={APP_ID}
            action={ACTION}
            rp_context={rpContext}
            allow_legacy_proofs={true}
            preset={orbLegacy({ signal: '' })}
            environment="staging"
            handleVerify={handleWidgetVerify}
            onSuccess={onSuccess}
            onError={onError}
          />
        )}

        <button
          onClick={handleVerify}
          disabled={fetching}
          className="verify-btn"
        >
          {fetching ? (
            <>
              <div className="w-4 h-4 border-2 border-black/20 border-t-black/60 rounded-full animate-spin" />
              <span>Preparing…</span>
            </>
          ) : (
            <>
              <QrCode size={15} strokeWidth={2.5} />
              <span>Verify with World ID</span>
            </>
          )}
        </button>

        {error && (
          <p className="mt-5 text-white/30 text-xs text-center">{error}</p>
        )}

        <p className="absolute bottom-12 label-mono opacity-30">
          Powered by Worldcoin
        </p>
      </div>
    </div>
  );
}
