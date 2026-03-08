'use client';

import React, { useEffect, useRef, memo } from 'react';

interface Props {
  symbol: string; // e.g. "BINANCE:BTCUSDT"
}

function TradingViewWidget({ symbol }: Props) {
  const container = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!container.current) return;
    container.current.innerHTML = '';

    const script = document.createElement('script');
    script.src =
      'https://s3.tradingview.com/external-embedding/embed-widget-advanced-chart.js';
    script.type = 'text/javascript';
    script.async = true;
    script.innerHTML = JSON.stringify({
      symbol,
      interval: 'D',
      style: '1',
      theme: 'dark',
      locale: 'en',
      timezone: 'Etc/UTC',
      backgroundColor: '#0f1011',
      gridColor: 'rgba(255,255,255,0.04)',
      allow_symbol_change: false,
      hide_side_toolbar: true,
      hide_top_toolbar: false,
      hide_legend: false,
      hide_volume: false,
      save_image: false,
      autosize: true,
    });

    container.current.appendChild(script);

    return () => {
      if (container.current) container.current.innerHTML = '';
    };
  }, [symbol]);

  return (
    <div ref={container} style={{ height: '100%', width: '100%' }} />
  );
}

export default memo(TradingViewWidget);
