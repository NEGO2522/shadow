'use client';

import { useState, useEffect } from 'react';
import { ArrowLeft, TrendingUp, TrendingDown, Shield, X } from 'lucide-react';
import { BottomNav } from '@/components/BottomNav';
import TradingViewWidget from '@/components/TradingViewWidget';
import { MiniKit } from '@worldcoin/minikit-js';
import shadowAbi from '@/abi/shadow.abi.json';

const SHADOW_ADDRESS = (process.env.NEXT_PUBLIC_SHADOW_CONTRACT ??
  '0x7aD3A7BeCe749eBaBEDF20fe2c688C525479e5De') as `0x${string}`;

interface Asset {
  id: string;
  symbol: string;
  name: string;
  tvSymbol: string;
  icon: string;
  price: number;
  change24h: number;
}

const ASSETS: Omit<Asset, 'price' | 'change24h'>[] = [
  { id: 'bitcoin',       symbol: 'BTC',  name: 'Bitcoin',   tvSymbol: 'BINANCE:BTCUSDT',  icon: '₿' },
  { id: 'ethereum',      symbol: 'ETH',  name: 'Ethereum',  tvSymbol: 'BINANCE:ETHUSDT',  icon: 'Ξ' },
  { id: 'solana',        symbol: 'SOL',  name: 'Solana',    tvSymbol: 'BINANCE:SOLUSDT',  icon: '◎' },
  { id: 'binancecoin',   symbol: 'BNB',  name: 'BNB',       tvSymbol: 'BINANCE:BNBUSDT',  icon: '⬡' },
  { id: 'avalanche-2',   symbol: 'AVAX', name: 'Avalanche', tvSymbol: 'BINANCE:AVAXUSDT', icon: '▲' },
  { id: 'uniswap',       symbol: 'UNI',  name: 'Uniswap',   tvSymbol: 'BINANCE:UNIUSDT',  icon: '🦄' },
  { id: 'chainlink',     symbol: 'LINK', name: 'Chainlink', tvSymbol: 'BINANCE:LINKUSDT', icon: '⬡' },
  { id: 'aave',          symbol: 'AAVE', name: 'Aave',      tvSymbol: 'BINANCE:AAVEUSDT', icon: '👻' },
  { id: 'arbitrum',      symbol: 'ARB',  name: 'Arbitrum',  tvSymbol: 'BINANCE:ARBUSDT',  icon: '🔵' },
  { id: 'worldcoin-wld', symbol: 'WLD',  name: 'Worldcoin', tvSymbol: 'BINANCE:WLDUSDT',  icon: '🌍' },
];

const COINGECKO_URL =
  `https://api.coingecko.com/api/v3/simple/price` +
  `?ids=${ASSETS.map((a) => a.id).join(',')}` +
  `&vs_currencies=usd&include_24hr_change=true`;

function fmt(n: number) {
  if (n >= 1000) return n.toLocaleString('en-US', { maximumFractionDigits: 2 });
  if (n >= 1) return n.toFixed(2);
  return n.toFixed(4);
}

type OrderSide = 'buy' | 'sell';
type OrderState = 'idle' | 'pending' | 'success' | 'failed';

export default function TradingPage() {
  const [assets, setAssets] = useState<Asset[]>([]);
  const [loading, setLoading] = useState(true);
  const [selected, setSelected] = useState<Asset | null>(null);
  const [showOrderModal, setShowOrderModal] = useState(false);

  // Order form
  const [orderSide, setOrderSide] = useState<OrderSide>('buy');
  const [orderPrice, setOrderPrice] = useState('');
  const [orderQty, setOrderQty] = useState('');
  const [orderState, setOrderState] = useState<OrderState>('idle');

  useEffect(() => {
    async function fetchPrices() {
      try {
        const res = await fetch(COINGECKO_URL);
        const data = await res.json();
        setAssets(
          ASSETS.map((a) => ({
            ...a,
            price: data[a.id]?.usd ?? 0,
            change24h: data[a.id]?.usd_24h_change ?? 0,
          })),
        );
      } catch {
        setAssets(ASSETS.map((a) => ({ ...a, price: 0, change24h: 0 })));
      } finally {
        setLoading(false);
      }
    }
    fetchPrices();
    const id = setInterval(fetchPrices, 60_000);
    return () => clearInterval(id);
  }, []);

  // Pre-fill order price from current market price
  const openOrderModal = (asset: Asset) => {
    setSelected(asset);
    setOrderPrice(asset.price > 0 ? fmt(asset.price) : '');
    setOrderQty('');
    setOrderState('idle');
    setShowOrderModal(true);
  };

  const handlePlaceOrder = async () => {
    if (!selected || !orderPrice || !orderQty) {
      alert('Fill in price and quantity.');
      return;
    }
    setOrderState('pending');
    try {
      // Encrypt the order server-side
      const encRes = await fetch('/api/shielded-order', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          side: orderSide,
          asset: selected.symbol,
          price: parseFloat(orderPrice),
          quantity: parseFloat(orderQty),
        }),
      });
      if (!encRes.ok) throw new Error('Order encryption failed');
      const { encryptedOrder, orderId } = await encRes.json();

      // Submit to Shadow.placeOrder on-chain (no USDC transfer needed)
      const { finalPayload } = await MiniKit.commandsAsync.sendTransaction({
        transaction: [
          {
            address: SHADOW_ADDRESS,
            abi: shadowAbi,
            functionName: 'placeOrder',
            args: [encryptedOrder, orderId],
          },
        ],
      });

      if (finalPayload.status === 'success') {
        setOrderState('success');
        setTimeout(() => setShowOrderModal(false), 1500);
      } else {
        setOrderState('failed');
      }
    } catch (err) {
      console.error(err);
      setOrderState('failed');
    } finally {
      setTimeout(() => { if (orderState !== 'success') setOrderState('idle'); }, 3000);
    }
  };

  // ── Chart view ────────────────────────────────────────────────
  if (selected && !showOrderModal) {
    return (
      <div className="app-screen pb-24">
        <div className="screen-header">
          <button onClick={() => setSelected(null)} className="btn-icon-round">
            <ArrowLeft size={18} strokeWidth={2.5} />
          </button>
          <div className="text-center">
            <p className="text-base font-semibold">{selected.name}</p>
            <p className="text-[11px] text-white/40">{selected.symbol}/USDT</p>
          </div>
          <div className="text-right">
            <p className="text-sm font-semibold">${fmt(selected.price)}</p>
            <p className={`text-[11px] font-bold ${selected.change24h >= 0 ? 'text-emerald-400' : 'text-red-400'}`}>
              {selected.change24h >= 0 ? '+' : ''}{selected.change24h.toFixed(2)}%
            </p>
          </div>
        </div>

        <div className="px-4" style={{ height: 'calc(100dvh - 108px - 140px)' }}>
          <div className="rounded-[20px] overflow-hidden border border-white/5 h-full">
            <TradingViewWidget symbol={selected.tvSymbol} />
          </div>
        </div>

        {/* Dark pool order button */}
        <div className="px-6 mt-4">
          <button
            onClick={() => openOrderModal(selected)}
            className="w-full flex items-center justify-center gap-2 py-4 rounded-[28px] bg-emerald-500/20 text-emerald-400 border border-emerald-500/30 text-sm font-bold hover:bg-emerald-500/30 transition-all active:scale-95"
          >
            <Shield size={15} />
            Place Dark Pool Order
          </button>
        </div>

        <BottomNav />
      </div>
    );
  }

  // ── Order Modal ───────────────────────────────────────────────
  if (showOrderModal && selected) {
    return (
      <div className="app-screen items-center justify-center relative">
        <div className="absolute inset-0 bg-black/60 backdrop-blur-sm" onClick={() => setShowOrderModal(false)} />
        <div className="relative z-10 mx-6 w-full max-w-sm bg-[#1a1b1f] border border-white/10 rounded-[32px] p-6 space-y-5">
          <div className="flex justify-between items-center">
            <div>
              <p className="text-base font-semibold">Dark Pool Order</p>
              <p className="text-[11px] text-white/40 mt-0.5">{selected.name} · {selected.symbol}/USDC</p>
            </div>
            <button onClick={() => setShowOrderModal(false)} className="btn-icon-round">
              <X size={16} />
            </button>
          </div>

          {/* Side toggle */}
          <div className="flex gap-2 p-1 bg-black/20 rounded-2xl">
            {(['buy', 'sell'] as OrderSide[]).map((s) => (
              <button
                key={s}
                onClick={() => setOrderSide(s)}
                className={`flex-1 py-2.5 rounded-xl text-xs font-bold transition-all capitalize ${
                  orderSide === s
                    ? s === 'buy'
                      ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'
                      : 'bg-red-500/20 text-red-400 border border-red-500/30'
                    : 'text-neutral-400 hover:text-white'
                }`}
              >
                {s}
              </button>
            ))}
          </div>

          {/* Price */}
          <div className="flex flex-col gap-1.5">
            <label className="label-mono">Price (USDC)</label>
            <input
              type="number"
              placeholder="0.00"
              value={orderPrice}
              onChange={(e) => setOrderPrice(e.target.value)}
              className="pay-input"
            />
          </div>

          {/* Quantity */}
          <div className="flex flex-col gap-1.5">
            <label className="label-mono">Quantity ({selected.symbol})</label>
            <input
              type="number"
              placeholder="0.00"
              value={orderQty}
              onChange={(e) => setOrderQty(e.target.value)}
              className="pay-input"
            />
          </div>

          <div className="flex items-start gap-2 bg-emerald-500/10 border border-emerald-500/20 rounded-2xl p-3">
            <Shield size={13} className="text-emerald-400 mt-0.5 shrink-0" />
            <p className="text-[11px] text-emerald-300/80 leading-relaxed">
              Order is encrypted on-chain. CRE DON decrypts and routes to the dark pool matcher.
            </p>
          </div>

          <button
            onClick={handlePlaceOrder}
            disabled={orderState === 'pending'}
            className={`w-full py-4 rounded-[28px] text-sm font-bold transition-all active:scale-95 disabled:opacity-60 ${
              orderSide === 'buy'
                ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'
                : 'bg-red-500/20 text-red-400 border border-red-500/30'
            }`}
          >
            {orderState === 'pending'
              ? 'Submitting…'
              : orderState === 'success'
              ? 'Order Placed!'
              : orderState === 'failed'
              ? 'Failed — try again'
              : `Place ${orderSide.charAt(0).toUpperCase() + orderSide.slice(1)} Order`}
          </button>
        </div>
      </div>
    );
  }

  // ── Asset list view ───────────────────────────────────────────
  return (
    <div className="app-screen pb-24 overflow-y-auto no-scrollbar">
      <div className="screen-header">
        <div>
          <h1 className="text-xl font-medium tracking-tight">Markets</h1>
          <p className="subtitle-sm mt-1">Crypto Assets</p>
        </div>
        <div className="text-[10px] text-white/20 font-mono uppercase tracking-widest">
          Live
        </div>
      </div>

      <div className="px-6 flex justify-between mb-3">
        <span className="label-mono">Asset</span>
        <div className="flex gap-6">
          <span className="label-mono">24h</span>
          <span className="label-mono w-20 text-right">Price</span>
        </div>
      </div>

      <div className="px-4 space-y-2">
        {loading
          ? Array.from({ length: 10 }).map((_, i) => (
              <div key={i} className="list-item-card animate-pulse">
                <div className="flex items-center gap-4">
                  <div className="w-10 h-10 rounded-full bg-white/5" />
                  <div className="space-y-2">
                    <div className="h-3 w-20 bg-white/5 rounded-full" />
                    <div className="h-2 w-12 bg-white/5 rounded-full" />
                  </div>
                </div>
                <div className="h-3 w-16 bg-white/5 rounded-full" />
              </div>
            ))
          : assets.map((asset) => {
              const up = asset.change24h >= 0;
              return (
                <button
                  key={asset.id}
                  onClick={() => setSelected(asset)}
                  className="list-item-card w-full text-left active:scale-[0.98] transition-transform"
                >
                  <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-full bg-[#1e1e24] border border-white/5 flex items-center justify-center text-lg">
                      {asset.icon}
                    </div>
                    <div>
                      <p className="text-sm font-semibold text-white">{asset.name}</p>
                      <p className="text-[10px] text-white/30 font-mono uppercase tracking-wider mt-0.5">
                        {asset.symbol}
                      </p>
                    </div>
                  </div>
                  <div className="flex items-center gap-5">
                    <div className={`flex items-center gap-1 text-xs font-bold ${up ? 'text-emerald-400' : 'text-red-400'}`}>
                      {up ? <TrendingUp size={12} /> : <TrendingDown size={12} />}
                      {up ? '+' : ''}{asset.change24h.toFixed(2)}%
                    </div>
                    <p className="text-sm font-semibold text-white w-20 text-right">
                      ${fmt(asset.price)}
                    </p>
                  </div>
                </button>
              );
            })}
      </div>

      <BottomNav />
    </div>
  );
}
