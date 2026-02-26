'use client';

import { LogIn, LogOut, User } from 'lucide-react';
import Link from 'next/link';
import { useEffect } from 'react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useAuth } from '@/hooks/useAuth';

export function AuthButton() {
  const { user, isAuthenticated, logout, initialize } = useAuth();

  useEffect(() => {
    // コンポーネントマウント時にCookieから認証状態を初期化
    initialize();
  }, [initialize]);

  if (isAuthenticated && user) {
    return (
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" className="relative h-10 w-10 rounded-full">
            <Avatar className="h-10 w-10">
              <AvatarImage src={user.icon_url || undefined} alt={user.name} />
              <AvatarFallback className="bg-primary text-primary-foreground">
                {user.name?.slice(0, 2).toUpperCase() || 'U'}
              </AvatarFallback>
            </Avatar>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-56" align="end" forceMount>
          <DropdownMenuLabel className="font-normal">
            <div className="flex flex-col space-y-1">
              <p className="text-sm font-medium leading-none">{user.name}</p>
              <p className="text-xs leading-none text-muted-foreground">{user.email}</p>
            </div>
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          <DropdownMenuItem asChild className="cursor-pointer">
            <Link href={`/users/${user.id}`}>
              <User className="mr-2 h-4 w-4" />
              <span>プロフィール</span>
            </Link>
          </DropdownMenuItem>
          <DropdownMenuItem onClick={logout} className="cursor-pointer">
            <LogOut className="mr-2 h-4 w-4" />
            <span>ログアウト</span>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    );
  }

  return (
    <div className="flex gap-2">
      <Button variant="ghost" asChild>
        <Link href="/auth">
          <LogIn className="mr-2 h-4 w-4" />
          ログイン
        </Link>
      </Button>
    </div>
  );
}
