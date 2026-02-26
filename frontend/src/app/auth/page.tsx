'use client';

import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import type { SignUpRequestAffiliation } from '@/generated/models';
import { useAuth } from '@/hooks/useAuth';

export default function AuthPage() {
  const router = useRouter();
  const { login, signup } = useAuth();

  // ログインフォーム
  const [loginEmail, setLoginEmail] = useState('');
  const [loginPassword, setLoginPassword] = useState('');
  const [loginError, setLoginError] = useState('');
  const [loginLoading, setLoginLoading] = useState(false);

  // サインアップフォーム
  const [signupName, setSignupName] = useState('');
  const [signupEmail, setSignupEmail] = useState('');
  const [signupPassword, setSignupPassword] = useState('');
  const [signupConfirmPassword, setSignupConfirmPassword] = useState('');
  const [signupAffiliation, setSignupAffiliation] = useState<SignUpRequestAffiliation | ''>('');
  const [signupError, setSignupError] = useState('');
  const [signupLoading, setSignupLoading] = useState(false);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoginError('');
    setLoginLoading(true);

    try {
      await login(loginEmail, loginPassword);
      router.push('/');
    } catch (error) {
      setLoginError(error instanceof Error ? error.message : 'ログインに失敗しました');
    } finally {
      setLoginLoading(false);
    }
  };

  const handleSignup = async (e: React.FormEvent) => {
    e.preventDefault();
    setSignupError('');

    if (signupPassword !== signupConfirmPassword) {
      setSignupError('パスワードが一致しません');
      return;
    }

    if (signupPassword.length < 8) {
      setSignupError('パスワードは8文字以上である必要があります');
      return;
    }

    setSignupLoading(true);

    try {
      await signup(signupEmail, signupPassword, signupName, signupAffiliation || undefined);
      router.push('/');
    } catch (error) {
      setSignupError(error instanceof Error ? error.message : 'サインアップに失敗しました');
    } finally {
      setSignupLoading(false);
    }
  };

  return (
    <main className="min-h-screen bg-slate-50 flex items-center justify-center p-4">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle className="text-2xl text-center">えーブログ</CardTitle>
          <CardDescription className="text-center">
            サークルメンバー専用の認証ページ
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Tabs defaultValue="login" className="w-full">
            <TabsList className="grid w-full grid-cols-2">
              <TabsTrigger value="login">ログイン</TabsTrigger>
              <TabsTrigger value="signup">新規登録</TabsTrigger>
            </TabsList>

            {/* ログインタブ */}
            <TabsContent value="login">
              <form onSubmit={handleLogin} className="space-y-4 mt-4">
                <div className="space-y-2">
                  <label htmlFor="login-email" className="text-sm font-medium">
                    メールアドレス
                  </label>
                  <Input
                    id="login-email"
                    type="email"
                    placeholder="example@email.com"
                    value={loginEmail}
                    onChange={(e) => setLoginEmail(e.target.value)}
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label htmlFor="login-password" className="text-sm font-medium">
                    パスワード
                  </label>
                  <Input
                    id="login-password"
                    type="password"
                    placeholder="••••••••"
                    value={loginPassword}
                    onChange={(e) => setLoginPassword(e.target.value)}
                    required
                  />
                </div>
                {loginError && <p className="text-sm text-red-500">{loginError}</p>}
                <Button type="submit" className="w-full" disabled={loginLoading}>
                  {loginLoading ? 'ログイン中...' : 'ログイン'}
                </Button>
              </form>
            </TabsContent>

            {/* サインアップタブ */}
            <TabsContent value="signup">
              <form onSubmit={handleSignup} className="space-y-4 mt-4">
                <div className="space-y-2">
                  <label htmlFor="signup-name" className="text-sm font-medium">
                    名前
                  </label>
                  <Input
                    id="signup-name"
                    type="text"
                    placeholder="山田太郎"
                    value={signupName}
                    onChange={(e) => setSignupName(e.target.value)}
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label htmlFor="signup-affiliation" className="text-sm font-medium">
                    所属（任意）
                  </label>
                  <select
                    id="signup-affiliation"
                    value={signupAffiliation}
                    onChange={(e) =>
                      setSignupAffiliation(e.target.value as SignUpRequestAffiliation | '')
                    }
                    className="flex h-10 w-full rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
                  >
                    <option value="">選択してください</option>
                    <option value="開発">開発</option>
                    <option value="マーケティング">マーケティング</option>
                    <option value="組織管理">組織管理</option>
                  </select>
                </div>
                <div className="space-y-2">
                  <label htmlFor="signup-email" className="text-sm font-medium">
                    メールアドレス
                  </label>
                  <Input
                    id="signup-email"
                    type="email"
                    placeholder="example@email.com"
                    value={signupEmail}
                    onChange={(e) => setSignupEmail(e.target.value)}
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label htmlFor="signup-password" className="text-sm font-medium">
                    パスワード（8文字以上）
                  </label>
                  <Input
                    id="signup-password"
                    type="password"
                    placeholder="••••••••"
                    value={signupPassword}
                    onChange={(e) => setSignupPassword(e.target.value)}
                    required
                  />
                </div>
                <div className="space-y-2">
                  <label htmlFor="signup-confirm-password" className="text-sm font-medium">
                    パスワード確認
                  </label>
                  <Input
                    id="signup-confirm-password"
                    type="password"
                    placeholder="••••••••"
                    value={signupConfirmPassword}
                    onChange={(e) => setSignupConfirmPassword(e.target.value)}
                    required
                  />
                </div>
                {signupError && <p className="text-sm text-red-500">{signupError}</p>}
                <Button type="submit" className="w-full" disabled={signupLoading}>
                  {signupLoading ? '登録中...' : '新規登録'}
                </Button>
              </form>
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>
    </main>
  );
}
