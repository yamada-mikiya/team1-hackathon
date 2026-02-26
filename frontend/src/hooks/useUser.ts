'use client';

import { useQuery } from '@tanstack/react-query';
import type { ArticleResponse } from '@/generated/models';
import { AXIOS_INSTANCE } from '@/lib/api-client';

export interface UserDetail {
  id: number;
  name: string;
  affiliation?: string | null;
  icon_url?: string | null;
  portfolio_key?: string | null;
  created_at: string;
  articles: ArticleResponse[];
}

interface UseUserOptions {
  userId?: number;
  portfolioKey?: string;
}

/**
 * ユーザー詳細情報を取得するカスタムフック
 * @param userId - 取得するユーザーのID
 * @param portfolioKey - ポートフォリオキー（未認証でアクセスする場合に必要）
 */
export function useUser({ userId, portfolioKey }: UseUserOptions) {
  const { data, error, isLoading } = useQuery({
    queryKey: ['user', userId, portfolioKey],
    queryFn: async () => {
      if (!userId) return null;
      const params = portfolioKey ? { portfolio_key: portfolioKey } : {};
      const response = await AXIOS_INSTANCE.get<UserDetail>(`/api/users/${userId}`, { params });
      return response.data;
    },
    enabled: !!userId,
  });

  return {
    user: data,
    isLoading,
    error,
  };
}
