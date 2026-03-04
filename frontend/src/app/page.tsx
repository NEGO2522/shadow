'use client';

import Image from 'next/image';
import Link from 'next/link';
import { useState } from 'react';
import {
  MessageSquare,
  Send,
  Home,
  CreditCard,
  Receipt,
  User,
  TrendingUp,
  QrCode
} from 'lucide-react';

import { IDKit, orbLegacy } from "@worldcoin/idkit-core";

export default function Dashboard() {
  const [hoveredNav, setHoveredNav] = useState<string | null>(null);
  const [isVerifying, setIsVerifying] = useState(false);

  const handleVerify = async () => {
    setIsVerifying(true);
    const actionName = process.env.NEXT_PUBLIC_WLD_ACTION_NAME || "payment-verification";
    const appId = process.env.NEXT_PUBLIC_APP_ID as `app_${string}`;

    try {
      const rpSig = await fetch("/api/rp-signature", {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ action: actionName }),
      }).then((r) => r.json());

      const request = await IDKit.request({
        app_id: appId,
        action: actionName,
        rp_context: {
          rp_id: "rp_3dc77f0d48e60297", // Your app's `rp_id` from the Developer Portal
          nonce: rpSig.nonce,
          created_at: rpSig.created_at,
          expires_at: rpSig.expires_at,
          signature: rpSig.sig,
        },
        allow_legacy_proofs: true,
        environment: "production", // simulator uses production too
      }).preset(orbLegacy({ signal: "" }));

      const connectUrl = request.connectorURI;
      console.log("Connect URL:", connectUrl);
      window.open(connectUrl, '_blank'); // Open the URL so the developer can run it in simulator

      const response = await request.pollUntilCompletion();
      console.log("Proof collected:", response);

      const res = await fetch('/api/verify-proof', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          payload: response,
          action: actionName,
          signal: "",
        }),
      });

      if (!res.ok) {
        throw new Error('Verification failed on the backend');
      }

      console.log("Proof validated successfully!");
      alert('Successfully verified with World ID!');
    } catch (error) {
      console.error("Error during verification:", error);
      alert("Verification failed. Please try again.");
    } finally {
      setIsVerifying(false);
    }
  };

  return (
    <div className="app-screen pb-24">
      {/* Header */}
      <div className="screen-header">
        <div>
          <h1 className="text-xl font-medium tracking-tight">
            Hi 👋, Daniel
          </h1>
          <p className="subtitle-sm mt-1">
            Welcome Back
          </p>
        </div>
        <div className="flex items-center gap-4">
          <div className="h-10 w-10 rounded-full bg-[#1a1a20] flex items-center justify-center p-1 border border-white/5 overflow-hidden ring-2 ring-transparent transition hover:ring-white/20">
            <Image
              src="https://api.dicebear.com/7.x/notionists/svg?seed=Daniel"
              alt="User Avatar"
              width={32}
              height={32}
              className="rounded-full bg-white/10"
            />
          </div>
          <button onClick={() => alert('Messages')} className="btn-icon-round relative">
            <MessageSquare size={18} strokeWidth={2} />
            <div className="absolute top-2 right-2.5 h-1.5 w-1.5 bg-red-500 rounded-full border border-[#1e1e24]"></div>
          </button>
        </div>
      </div>

      {/* Main Card */}
      <div className="px-6 relative z-10">
        <div className="gradient-card">
          <div className="absolute top-0 right-0 w-40 h-40 bg-white/5 rounded-full blur-3xl -mr-10 -mt-10"></div>

          <div className="flex justify-between items-start">
            <div className="space-y-1">
              <p className="text-white/80 text-sm font-medium">Total Balance</p>
              <h2 className="title-xl">$12,650.00</h2>
            </div>
            <div className="flex -space-x-3 opacity-90">
              <div className="w-8 h-8 rounded-full bg-red-500/80 mix-blend-multiply"></div>
              <div className="w-8 h-8 rounded-full bg-yellow-500/80 mix-blend-multiply"></div>
            </div>
          </div>

          <div className="mt-4 flex items-center gap-2">
            <div className="flex items-center gap-1.5 bg-white/10 px-2 py-1 rounded-full backdrop-blur-md">
              <TrendingUp size={12} className="text-white" />
              <span className="text-[10px] text-white font-medium">2.00% than last week</span>
            </div>
          </div>

          <div className="mt-8 flex items-end justify-between text-white/90 font-mono text-sm tracking-widest relative z-10">
            <p>**** **** **** 9946</p>
            <p className="font-sans text-xs tracking-wide">12 / 25</p>
          </div>

          {/* Action Row inside Card */}
          <div className="mt-6 flex gap-3 relative z-10">

            {/* WORLD ID INTEGRATION START */}
            <button
              onClick={handleVerify}
              disabled={isVerifying}
              className="btn-action-dark group disabled:opacity-50"
            >
              <div className="w-6 h-6 rounded-full bg-white/10 flex items-center justify-center group-hover:bg-emerald-400/20 transition-colors">
                <QrCode size={14} className={`text-emerald-400 ${isVerifying ? 'animate-pulse' : ''}`} />
              </div>
              {isVerifying ? 'Verifying...' : 'Pay'}
            </button>
            {/* WORLD ID INTEGRATION END */}

            <Link href="/transfer" className="btn-action-primary">
              <Send size={14} />
              Transfer
            </Link>
          </div>
        </div>
      </div>

      {/* Latest Transaction */}
      <div className="px-6 mt-10">
        <h3 className="text-sm font-medium tracking-wide text-white/90 mb-4 px-1">Latest Transaction</h3>
        <div className="space-y-3">
          {[
            { name: 'Transfer for Work', date: '1 Sep 9:29', status: 'Pending', statusColor: 'text-neutral-400', amount: '+$120', avatar: 'Felix' },
            { name: 'Transfer for Work', date: '1 Sep 9:29', status: 'Success', statusColor: 'text-emerald-500', amount: '-$120', avatar: 'Aneka' },
            { name: 'Transfer for Work', date: '1 Sep 9:29', status: 'Success', statusColor: 'text-emerald-500', amount: '+$120', avatar: 'Sara' }
          ].map((tx, idx) => (
            <div key={idx} className="list-item-card">
              <div className="flex items-center gap-4">
                <div className="w-12 h-12 rounded-full bg-[#1e1e24] flex items-center justify-center overflow-hidden border border-white/5">
                  <Image src={`https://api.dicebear.com/7.x/avataaars/svg?seed=${tx.avatar}`} alt={tx.avatar} width={36} height={36} />
                </div>
                <div>
                  <p className="font-medium text-white text-sm">{tx.name}</p>
                  <p className="text-[10px] uppercase tracking-wider text-neutral-500 mt-1 font-bold">
                    {tx.date} | <span className={tx.statusColor}>{tx.status}</span>
                  </p>
                </div>
              </div>
              <p className="text-white font-semibold text-sm">{tx.amount}</p>
            </div>
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
              className={`nav-link ${hoveredNav === nav.id || (hoveredNav === null && nav.id === 'home') ? 'nav-link-active' : 'nav-link-inactive'}`}
            >
              <nav.icon size={18} className="shrink-0" />
              <span className={`overflow-hidden text-xs font-semibold transition-all duration-300 ${hoveredNav === nav.id || (hoveredNav === null && nav.id === 'home') ? nav.width + ' opacity-100 ml-1.5' : 'w-0 opacity-0'}`}>
                {nav.label}
              </span>
            </Link>
          ))}
        </div>
      </div>
    </div>
  );
}