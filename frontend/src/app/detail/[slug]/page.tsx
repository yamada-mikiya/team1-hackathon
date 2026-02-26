'use client';

import { ArrowLeft, Calendar, ExternalLink, User } from 'lucide-react';
import Image from 'next/image';
import Link from 'next/link';
import { useParams } from 'next/navigation';
import ReactMarkdown, { type Components } from 'react-markdown';
import rehypeRaw from 'rehype-raw';
import remarkGfm from 'remark-gfm';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent } from '@/components/ui/card';
import { useArticle } from '@/hooks/useArticle';

export default function ArticleDetailPage() {
  const params = useParams();
  const slug = params.slug as string;

  const { data: article, isLoading, error } = useArticle(slug);

  const formatDate = (dateString?: string) => {
    if (!dateString) return '';
    return new Date(dateString).toLocaleDateString('ja-JP', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  const getDepartmentLabel = (dept?: string) => {
    switch (dept) {
      case 'MKT':
        return 'マーケティング';
      case 'Dev':
        return '開発';
      case 'Ops':
        return '組織管理';
      default:
        return dept;
    }
  };

  if (isLoading) {
    return (
      <main className="min-h-screen bg-slate-50">
        <div className="container mx-auto px-4 py-12 max-w-4xl">
          <div className="text-center py-20 text-muted-foreground">Loading...</div>
        </div>
      </main>
    );
  }

  if (error || !article) {
    return (
      <main className="min-h-screen bg-slate-50">
        <div className="container mx-auto px-4 py-12 max-w-4xl">
          <div className="text-center py-20 text-red-500">
            記事が見つかりませんでした。
            <div className="mt-4">
              <Link
                href="/"
                className="text-primary hover:underline inline-flex items-center gap-2"
              >
                <ArrowLeft className="h-4 w-4" />
                トップページに戻る
              </Link>
            </div>
          </div>
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
          <ArrowLeft className="h-4 w-4" />
          <span className="text-sm font-medium">記事一覧に戻る</span>
        </Link>

        <article>
          {/* Article Header */}
          <header className="mb-8">
            <div className="flex items-center gap-3 mb-4">
              <Badge variant="outline" className="text-xs font-medium border-slate-200">
                {getDepartmentLabel(article.department)}
              </Badge>
              <Badge variant="outline" className="text-xs font-medium border-slate-200">
                {article.article_type}
              </Badge>
            </div>

            <h1 className="text-4xl font-bold text-slate-800 mb-6 leading-tight">
              {article.title}
            </h1>

            {/* Author and Date Info */}
            <div className="flex items-center gap-6 text-sm text-slate-600">
              <Link
                href={`/users/${article.author?.id}`}
                className="flex items-center gap-2 hover:underline"
              >
                <Avatar className="h-10 w-10 border-2 border-slate-100">
                  <AvatarImage src={article.author?.icon_url} alt={article.author?.name} />
                  <AvatarFallback>{article.author?.name?.slice(0, 2)}</AvatarFallback>
                </Avatar>
                <span className="font-medium text-slate-700">{article.author?.name}</span>
              </Link>
              <div className="flex items-center gap-2">
                <Calendar className="h-4 w-4 text-slate-400" />
                <time className="text-slate-600">{formatDate(article.created_at)}</time>
              </div>
            </div>
          </header>

          {/* Thumbnail Image */}
          {article.thumbnail_url && (
            <div className="relative w-full aspect-video overflow-hidden rounded-lg mb-8 bg-slate-100 shadow-md">
              <Image
                src={article.thumbnail_url}
                alt={article.title || 'Article thumbnail'}
                fill
                className="object-cover"
                priority
              />
            </div>
          )}

          {/* Article Content */}
          <Card className="bg-white shadow-sm border-slate-200">
            <CardContent className="p-8 md:p-12">
              {article.article_type === 'external' && article.external_url ? (
                <div className="text-center py-12">
                  <ExternalLink className="h-16 w-16 mx-auto mb-4 text-slate-400" />
                  <h2 className="text-xl font-semibold mb-4 text-slate-700">
                    この記事は外部リンクです
                  </h2>
                  <p className="text-slate-600 mb-6">以下のリンクから記事をご覧ください。</p>
                  <a
                    href={article.external_url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-2 px-6 py-3 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition-colors font-medium"
                  >
                    <ExternalLink className="h-5 w-5" />
                    外部リンクを開く
                  </a>
                  <p className="text-xs text-slate-500 mt-4">{article.external_url}</p>
                </div>
              ) : (
                <div className="prose prose-slate max-w-none prose-headings:font-bold prose-headings:text-slate-800 prose-headings:mb-4 prose-headings:mt-8 first:prose-headings:mt-0 prose-h1:text-3xl prose-h1:border-b prose-h1:border-slate-200 prose-h1:pb-2 prose-h2:text-2xl prose-h2:border-b prose-h2:border-slate-200 prose-h2:pb-2 prose-h3:text-xl prose-h4:text-lg prose-p:text-slate-700 prose-p:leading-relaxed prose-p:mb-4 prose-a:text-primary prose-a:font-medium prose-a:no-underline hover:prose-a:underline prose-strong:text-slate-900 prose-strong:font-bold prose-em:italic prose-code:text-sm prose-code:bg-slate-100 prose-code:text-slate-800 prose-code:px-1.5 prose-code:py-0.5 prose-code:rounded prose-code:font-mono prose-code:before:content-[''] prose-code:after:content-[''] prose-pre:bg-slate-900 prose-pre:text-slate-50 prose-pre:p-4 prose-pre:rounded-lg prose-pre:overflow-x-auto prose-img:rounded-lg prose-img:shadow-md prose-img:my-6 prose-blockquote:border-l-4 prose-blockquote:border-primary prose-blockquote:pl-4 prose-blockquote:italic prose-blockquote:text-slate-600 prose-ul:list-disc prose-ul:my-4 prose-ul:pl-6 prose-ol:list-decimal prose-ol:my-4 prose-ol:pl-6 prose-li:text-slate-700 prose-li:my-2 prose-hr:border-slate-200 prose-hr:my-8 prose-table:border-collapse prose-table:my-6 prose-table:w-full prose-th:border prose-th:border-slate-300 prose-th:bg-slate-100 prose-th:p-3 prose-th:text-left prose-th:font-semibold prose-td:border prose-td:border-slate-300 prose-td:p-3">
                  {article.content ? (
                    <ReactMarkdown
                      remarkPlugins={[remarkGfm]}
                      rehypePlugins={[rehypeRaw]}
                      components={
                        {
                          a: ({ node, ...props }) => (
                            <a {...props} target="_blank" rel="noopener noreferrer" />
                          ),
                          h1: ({ node, ...props }) => (
                            <h1
                              className="text-3xl font-bold text-slate-800 border-b border-slate-200 pb-2 mb-4 mt-8 first:mt-0"
                              {...props}
                            />
                          ),
                          h2: ({ node, ...props }) => (
                            <h2
                              className="text-2xl font-bold text-slate-800 border-b border-slate-200 pb-2 mb-4 mt-8"
                              {...props}
                            />
                          ),
                          h3: ({ node, ...props }) => (
                            <h3 className="text-xl font-bold text-slate-800 mb-3 mt-6" {...props} />
                          ),
                          h4: ({ node, ...props }) => (
                            <h4 className="text-lg font-bold text-slate-800 mb-3 mt-6" {...props} />
                          ),
                          h5: ({ node, ...props }) => (
                            <h5
                              className="text-base font-bold text-slate-800 mb-2 mt-4"
                              {...props}
                            />
                          ),
                          h6: ({ node, ...props }) => (
                            <h6 className="text-sm font-bold text-slate-800 mb-2 mt-4" {...props} />
                          ),
                          p: ({ node, ...props }) => (
                            <p className="text-slate-700 leading-relaxed mb-4" {...props} />
                          ),
                          ul: ({ node, ...props }) => (
                            <ul className="list-disc pl-6 my-4 space-y-2" {...props} />
                          ),
                          ol: ({ node, ...props }) => (
                            <ol className="list-decimal pl-6 my-4 space-y-2" {...props} />
                          ),
                          li: ({ node, ...props }) => <li className="text-slate-700" {...props} />,
                          blockquote: ({ node, ...props }) => (
                            <blockquote
                              className="border-l-4 border-primary pl-4 italic text-slate-600 my-4"
                              {...props}
                            />
                          ),
                          code: ({ node, className, children, ...props }) => {
                            const isInline = !className;
                            return isInline ? (
                              <code
                                className="bg-slate-100 text-slate-800 px-1.5 py-0.5 rounded text-sm font-mono"
                                {...props}
                              >
                                {children}
                              </code>
                            ) : (
                              <code className={className} {...props}>
                                {children}
                              </code>
                            );
                          },
                          pre: ({ node, ...props }) => (
                            <pre
                              className="bg-slate-900 text-slate-50 p-4 rounded-lg overflow-x-auto my-6"
                              {...props}
                            />
                          ),
                          img: ({ node, ...props }) => (
                            <img
                              className="rounded-lg shadow-md my-6"
                              {...props}
                              alt={props.alt || ''}
                            />
                          ),
                          hr: ({ node, ...props }) => (
                            <hr className="border-slate-200 my-8" {...props} />
                          ),
                          table: ({ node, ...props }) => (
                            <div className="overflow-x-auto my-6">
                              <table className="border-collapse w-full" {...props} />
                            </div>
                          ),
                          th: ({ node, ...props }) => (
                            <th
                              className="border border-slate-300 bg-slate-100 p-3 text-left font-semibold"
                              {...props}
                            />
                          ),
                          td: ({ node, ...props }) => (
                            <td className="border border-slate-300 p-3" {...props} />
                          ),
                        } as Components
                      }
                    >
                      {article.content}
                    </ReactMarkdown>
                  ) : (
                    <p className="text-slate-500 text-center py-8">記事の内容がありません。</p>
                  )}
                </div>
              )}
            </CardContent>
          </Card>

          {/* Article Footer */}
          <footer className="mt-8 pt-6 border-t border-slate-200">
            <div className="flex items-center justify-between text-sm text-slate-500">
              <div>
                <span>最終更新: {formatDate(article.updated_at)}</span>
              </div>
              <div className="flex items-center gap-2">
                <User className="h-4 w-4" />
                <span>著者: {article.author?.name}</span>
              </div>
            </div>
          </footer>
        </article>
      </div>
    </main>
  );
}
