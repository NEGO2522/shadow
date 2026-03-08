'use client';

import Link from 'next/link';
import { useState } from 'react';
import { Home, BarChart2, Clock } from 'lucide-react';
import { usePathname } from 'next/navigation';

const navItems = [
  { id: 'home',    icon: Home,     label: 'Home',    href: '/',        width: 'w-[42px]' },
  { id: 'trade',   icon: BarChart2, label: 'Trade',  href: '/trading', width: 'w-[46px]' },
  { id: 'history', icon: Clock,    label: 'History', href: '/history', width: 'w-[52px]' },
];

export function BottomNav() {
  const pathname = usePathname();
  const [hovered, setHovered] = useState<string | null>(null);

  const defaultActive =
    pathname === '/trading' ? 'trade' : pathname === '/history' ? 'history' : 'home';
  const activeId = hovered ?? defaultActive;

  return (
    <div className="bottom-nav-v2">
      <div className="nav-bar-pill" onMouseLeave={() => setHovered(null)}>
        {navItems.map((nav) => (
          <Link
            key={nav.id}
            href={nav.href}
            onMouseEnter={() => setHovered(nav.id)}
            className={`nav-link ${activeId === nav.id ? 'nav-link-active' : 'nav-link-inactive'}`}
          >
            <nav.icon size={18} className="shrink-0" />
            <span
              className={`overflow-hidden text-xs font-semibold transition-all duration-300 ${
                activeId === nav.id ? nav.width + ' opacity-100 ml-1.5' : 'w-0 opacity-0'
              }`}
            >
              {nav.label}
            </span>
          </Link>
        ))}
      </div>
    </div>
  );
}
