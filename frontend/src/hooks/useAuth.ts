'use client';

import { create } from 'zustand';
import { postApiAuthLogin, postApiAuthSignup } from '@/generated/api/認証-auth/認証-auth';
import type { SignUpRequestAffiliation, UserResponse } from '@/generated/models';
import { AXIOS_INSTANCE } from '@/lib/api-client';

interface AuthState {
  user: UserResponse | null;
  isAuthenticated: boolean;
  login: (email: string, password: string) => Promise<void>;
  signup: (
    email: string,
    password: string,
    name: string,
    affiliation?: SignUpRequestAffiliation,
  ) => Promise<void>;
  logout: () => Promise<void>;
  initialize: () => Promise<void>;
  setUser: (user: UserResponse | null) => void;
}

export const useAuth = create<AuthState>()((set) => ({
  user: null,
  isAuthenticated: false,

  initialize: async () => {
    try {
      // バックエンドの/api/auth/meエンドポイントを呼び出して認証状態を確認
      const response = await AXIOS_INSTANCE.get<UserResponse>('/api/auth/me');
      set({ user: response.data, isAuthenticated: true });
    } catch (_error) {
      // 認証されていない場合
      set({ user: null, isAuthenticated: false });
    }
  },

  setUser: (user) => {
    set({ user, isAuthenticated: !!user });
  },

  login: async (email, password) => {
    try {
      const response = await postApiAuthLogin({ email, password });
      // バックエンドがCookieを設定するが、3rd-party cookieブロック対策としてlocalStorageにも保存
      if (response.token) {
        localStorage.setItem('auth_token', response.token);
      }
      set({ user: response.user, isAuthenticated: true });
    } catch (error: unknown) {
      const message =
        error &&
        typeof error === 'object' &&
        'response' in error &&
        error.response &&
        typeof error.response === 'object' &&
        'data' in error.response &&
        error.response.data &&
        typeof error.response.data === 'object' &&
        'error' in error.response.data
          ? String(error.response.data.error)
          : 'ログインに失敗しました';
      throw new Error(message);
    }
  },

  signup: async (email, password, name, affiliation) => {
    try {
      const response = await postApiAuthSignup({
        email,
        password,
        name,
        affiliation: affiliation || undefined,
      });
      // バックエンドがCookieを設定するが、3rd-party cookieブロック対策としてlocalStorageにも保存
      if (response.token) {
        localStorage.setItem('auth_token', response.token);
      }
      set({ user: response.user, isAuthenticated: true });
    } catch (error: unknown) {
      const message =
        error &&
        typeof error === 'object' &&
        'response' in error &&
        error.response &&
        typeof error.response === 'object' &&
        'data' in error.response &&
        error.response.data &&
        typeof error.response.data === 'object' &&
        'error' in error.response.data
          ? String(error.response.data.error)
          : 'サインアップに失敗しました';
      throw new Error(message);
    }
  },

  logout: async () => {
    try {
      // バックエンドの/api/auth/logoutエンドポイントを呼び出してCookieを削除
      await AXIOS_INSTANCE.post('/api/auth/logout');
      // localStorageからトークンを削除
      localStorage.removeItem('auth_token');
      set({ user: null, isAuthenticated: false });
    } catch (_error) {
      // エラーでも状態をクリア
      localStorage.removeItem('auth_token');
      set({ user: null, isAuthenticated: false });
    }
  },
}));
