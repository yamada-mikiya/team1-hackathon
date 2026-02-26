import axios, { type AxiosRequestConfig } from 'axios';

export const AXIOS_INSTANCE = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8080',
  withCredentials: true, // Cookieを自動的に送信（開発環境用）
  headers: {
    'Content-Type': 'application/json',
  },
});

// リクエストインターセプター: localStorageからトークンを取得してAuthorizationヘッダーに設定
// 3rd-party cookieがブロックされる環境でも動作するようにする
AXIOS_INSTANCE.interceptors.request.use(
  (config) => {
    if (typeof window !== 'undefined') {
      const token = localStorage.getItem('auth_token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error) => Promise.reject(error),
);

// レスポンスインターセプター: 401エラーで認証ページにリダイレクト
// ただし、/api/auth/meへのリクエストは除外（認証状態確認のため）
AXIOS_INSTANCE.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // /api/auth/meへのリクエストは401でもリダイレクトしない（未ログイン状態）
      const isAuthMeRequest = error.config?.url?.includes('/api/auth/me');

      if (
        !isAuthMeRequest &&
        typeof window !== 'undefined' &&
        !window.location.pathname.includes('/auth')
      ) {
        window.location.href = '/auth';
      }
    }
    return Promise.reject(error);
  },
);

export const customInstance = <T>(
  config: AxiosRequestConfig,
  options?: AxiosRequestConfig,
): Promise<T> => {
  const promise = AXIOS_INSTANCE({
    ...config,
    ...options,
  }).then(({ data }) => data);

  return promise;
};

export default AXIOS_INSTANCE;
