'use client';

import { ArrowLeft, Building2, Calendar, Lock, Share2 } from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { use } from 'react';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardFooter, CardHeader } from '@/components/ui/card';
import type { ArticleResponse } from '@/generated/models';
import { useAuth } from '@/hooks/useAuth';
import { useUser } from '@/hooks/useUser';

interface UserPageProps {
  params: Promise<{
    id: string;
  }>;
}

export default function UserPage({ params }: UserPageProps) {
  const { id } = use(params);
  const userId = parseInt(id, 10);
  const searchParams = useSearchParams();
  const portfolioKey = searchParams.get('portfolio_key') || undefined;

  const { user, isLoading, error } = useUser({
    userId,
    portfolioKey,
  });
  const { isAuthenticated, user: authUser } = useAuth();

  const formatDate = (dateString?: string) => {
    if (!dateString) return '';
    return new Date(dateString).toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  // ポートフォリオリンクをコピー
  const copyPortfolioLink = () => {
    if (!user?.portfolio_key) {
      // portfolio_keyがない場合はダミーキーで生成（バックエンドで検証されるため問題なし）
      const url = `${window.location.origin}/users/${userId}?portfolio_key=PLACEHOLDER`;
      navigator.clipboard.writeText(url);
      alert(
        'ポートフォリオリンクをコピーしました！\n注意: データベースにポートフォリオキーが設定されていない可能性があります。',
      );
      return;
    }
    const url = `${window.location.origin}/users/${userId}?portfolio_key=${user.portfolio_key}`;
    navigator.clipboard.writeText(url);
    alert('ポートフォリオリンクをコピーしました！');
  };

  // 認証状態またはportfolio_keyに応じて記事をフィルタリング
  const showInternalArticles = isAuthenticated || !!portfolioKey;
  const visibleArticles =
    user?.articles?.filter((article) => {
      if (showInternalArticles) {
        return article.status === 'public' || article.status === 'internal';
      }
      return article.status === 'public';
    }) || [];

  if (isLoading) {
    return (
      <main className="min-h-screen bg-slate-50">
        <div className="container mx-auto px-4 py-8 max-w-4xl">
          <div className="animate-pulse">
            <div className="h-8 bg-slate-200 rounded w-1/4 mb-8" />
            <div className="bg-white rounded-lg shadow p-8">
              <div className="flex items-start gap-6">
                <div className="w-32 h-32 bg-slate-200 rounded-full" />
                <div className="flex-1 space-y-4">
                  <div className="h-8 bg-slate-200 rounded w-1/2" />
                  <div className="h-4 bg-slate-200 rounded w-1/3" />
                  <div className="h-4 bg-slate-200 rounded w-1/4" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </main>
    );
  }

  if (error || !user) {
    return (
      <main className="min-h-screen bg-slate-50">
        <div className="container mx-auto px-4 py-8 max-w-4xl">
          <Link
            href="/"
            className="inline-flex items-center gap-2 text-slate-600 hover:text-primary transition-colors mb-8"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            一覧に戻る
          </Link>
          <Card>
            <CardContent className="py-12 text-center">
              <p className="text-slate-500">ユーザーが見つかりませんでした</p>
            </CardContent>
          </Card>
        </div>
      </main>
    );
  }

  return (
    <main className="min-h-screen bg-slate-50">
      <div className="container mx-auto px-4 py-8 max-w-4xl">
        <Link
          href="/"
          className="inline-flex items-center gap-2 text-slate-600 hover:text-primary transition-colors mb-8"
        >
          <ArrowLeft className="mr-2 h-4 w-4" />
          一覧に戻る
        </Link>

        <Card>
          <CardHeader>
            <div className="flex items-start gap-6">
              <Avatar className="w-32 h-32">
                <AvatarImage src={user.icon_url || undefined} alt={user.name} />
                <AvatarFallback className="bg-primary text-primary-foreground text-3xl">
                  {user.name?.slice(0, 2).toUpperCase() || 'U'}
                </AvatarFallback>
              </Avatar>

              <div className="flex-1">
                <h1 className="text-3xl font-bold mb-4">{user.name}</h1>

                <div className="space-y-2 text-slate-600">
                  {user.affiliation && (
                    <div className="flex items-center gap-2">
                      <Building2 className="h-4 w-4" />
                      <span>{user.affiliation}</span>
                    </div>
                  )}

                  <div className="flex items-center gap-2">
                    <Calendar className="h-4 w-4" />
                    <span>登録日: {formatDate(user.created_at)}</span>
                  </div>
                </div>
              </div>
            </div>
          </CardHeader>

          <CardContent>
            <div className="border-t pt-6 mb-8">
              <div className="flex justify-between items-start mb-4">
                <h2 className="text-xl font-semibold">この著者について</h2>

                {/* ポートフォリオモードコントロール */}
                <div className="flex gap-2">
                  {isAuthenticated && authUser?.id === userId && (
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={copyPortfolioLink}
                      className="text-xs"
                    >
                      <Share2 className="h-3 w-3 mr-1" />
                      ポートフォリオリンクをコピー
                    </Button>
                  )}
                </div>
              </div>

              {portfolioKey && !isAuthenticated && (
                <div className="mb-4 p-3 bg-amber-50 border border-amber-200 rounded-md">
                  <p className="text-sm text-amber-800">
                    <Lock className="h-4 w-4 inline mr-1" />
                    ポートフォリオモード: 内部公開記事も表示しています
                  </p>
                </div>
              )}

              <p className="text-slate-600">
                {user.affiliation
                  ? `${user.affiliation}に所属しています。`
                  : '所属情報はありません。'}
              </p>
            </div>

            {/* 記事一覧 */}
            <div className="border-t pt-6">
              <h2 className="text-xl font-semibold mb-6">
                投稿記事 ({visibleArticles.length}件)
                {!showInternalArticles &&
                  user.articles &&
                  user.articles.length > visibleArticles.length && (
                    <span className="text-sm font-normal text-slate-500 ml-2">
                      （外部公開のみ）
                    </span>
                  )}
              </h2>
              {visibleArticles.length > 0 ? (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {visibleArticles.map((article) => (
                    <ArticleCard
                      key={article.id}
                      article={article}
                      formatDate={formatDate}
                      isAuthenticated={isAuthenticated}
                    />
                  ))}
                </div>
              ) : (
                <div className="text-center py-8">
                  <p className="text-slate-500 mb-4">
                    {!showInternalArticles && user.articles && user.articles.length > 0
                      ? '外部公開記事がありません'
                      : 'まだ記事がありません'}
                  </p>
                  {!showInternalArticles && user.articles && user.articles.length > 0 && (
                    <div className="flex flex-col items-center gap-3">
                      <div className="flex items-center gap-2 text-sm text-slate-600">
                        <Lock className="h-4 w-4" />
                        <span>内部公開記事を見るにはログインが必要です</span>
                      </div>
                      <Link href="/auth">
                        <Button variant="default" size="sm">
                          ログイン
                        </Button>
                      </Link>
                    </div>
                  )}
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </main>
  );
}

function ArticleCard({
  article,
  formatDate,
  isAuthenticated,
}: {
  article: ArticleResponse;
  formatDate: (date?: string) => string;
  isAuthenticated: boolean;
}) {
  const isExternal = article.article_type === 'external' && article.external_url;
  const href = isExternal ? article.external_url! : `/detail/${article.slug}`;

  const handleCardClick = () => {
    if (isExternal && article.external_url) {
      window.open(article.external_url, '_blank', 'noopener,noreferrer');
    } else {
      window.location.href = href;
    }
  };

  // 記事のステータスに応じたバッジの色とテキストを決定
  const getStatusBadge = () => {
    if (article.status === 'internal') {
      return (
        <Badge
          variant="secondary"
          className="text-[10px] bg-amber-100 text-amber-800 border-amber-200"
        >
          <Lock className="h-3 w-3 mr-1" />
          内部公開
        </Badge>
      );
    }
    if (article.status === 'public') {
      return (
        <Badge
          variant="default"
          className="text-[10px] bg-green-100 text-green-800 border-green-200"
        >
          外部公開
        </Badge>
      );
    }
    return null;
  };

  return (
    <Card
      className="group hover:shadow-lg transition-all duration-300 cursor-pointer overflow-hidden"
      onClick={handleCardClick}
    >
      {article.thumbnail_url && (
        <div className="relative w-full aspect-video overflow-hidden bg-slate-100">
          <Image
            src={article.thumbnail_url}
            alt={article.title || 'Article thumbnail'}
            fill
            className="object-cover group-hover:scale-105 transition-transform duration-500"
          />
        </div>
      )}

      <CardHeader className="pb-2">
        <div className="flex items-center gap-2 mb-2 flex-wrap">
          <Badge variant="outline" className="text-[10px]">
            {article.department}
          </Badge>
          <Badge variant="outline" className="text-[10px]">
            {article.article_type}
          </Badge>
          {getStatusBadge()}
        </div>
        <h3 className="text-base font-bold leading-snug group-hover:text-primary transition-colors line-clamp-2">
          {article.title}
        </h3>
      </CardHeader>

      <CardFooter className="pt-0 text-xs text-slate-400">
        {formatDate(article.created_at)}
      </CardFooter>
    </Card>
  );
}
