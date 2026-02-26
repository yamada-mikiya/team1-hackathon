'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useAuth } from '@/hooks/useAuth';
import { AuthButton } from './AuthButton';

export function Header() {
  const pathname = usePathname();
  const { isAuthenticated } = useAuth();

  return (
    <header className="sticky top-0 z-50 w-full border-b bg-white/95 backdrop-blur supports-[backdrop-filter]:bg-white/60">
      <div className="container mx-auto flex h-16 items-center justify-between px-4">
        <Link href="/" className="flex items-center space-x-2">
          <span className="text-2xl font-bold text-primary">えーブログ</span>
        </Link>

        <nav className="flex items-center gap-6">
          <Link
            href="/"
            className={`text-sm font-medium transition-colors hover:text-primary ${
              pathname === '/' ? 'text-primary' : 'text-slate-600'
            }`}
          >
            記事一覧
          </Link>
          {/* 認証済みユーザーのみAuthButtonを表示 */}
          {isAuthenticated && <AuthButton />}
        </nav>
      </div>
    </header>
  );
}
