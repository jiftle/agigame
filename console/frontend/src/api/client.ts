import axios, { type AxiosRequestConfig } from 'axios';

import { useAuthStore } from '@/stores/auth';
import { message } from '@/utils/antdApp';
import { history } from '@/utils/history';

/** 后端统一响应包体 */
export interface ApiEnvelope<T = unknown> {
  code: number;
  message: string;
  data: T;
}

type InternalConfig = AxiosRequestConfig & {
  _skipAuthRefresh?: boolean;
  _retried?: boolean;
};

export type RequestOptions = InternalConfig;

const baseURL = import.meta.env.VITE_API_BASE || '/api/v1';

export const http = axios.create({ baseURL, timeout: 20000 });

// 请求拦截：统一注入 Bearer token
http.interceptors.request.use((config) => {
  const { token } = useAuthStore.getState();
  if (token) {
    config.headers = config.headers ?? {};
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

let refreshPromise: Promise<boolean> | null = null;

async function refreshAccessToken(): Promise<boolean> {
  const { refreshToken } = useAuthStore.getState();
  if (!refreshToken) {
    return false;
  }
  try {
    const res = await http.post<ApiEnvelope<{ token: string; refreshToken: string }>>(
      '/auth/refresh',
      { refreshToken },
      { _skipAuthRefresh: true } as AxiosRequestConfig,
    );
    const body = res.data;
    if (body?.code === 0 && body.data?.token) {
      useAuthStore.getState().setTokens(body.data.token, body.data.refreshToken);
      return true;
    }
    return false;
  } catch {
    return false;
  }
}

function redirectToLogin(): void {
  useAuthStore.getState().clear();
  const { pathname, search, hash } = window.location;
  if (pathname !== '/login') {
    history.replace(`/login?redirect=${encodeURIComponent(pathname + search + hash)}`);
  }
}

// 响应拦截：业务码判定 + 401/402 无感刷新
http.interceptors.response.use(
  (response) => {
    const config = response.config as InternalConfig;
    const body = response.data as ApiEnvelope | undefined;

    if (!body || typeof body.code !== 'number' || body.code === 0) {
      return response;
    }
    if (config._skipAuthRefresh) {
      return Promise.reject(new Error(body.message || '请求失败'));
    }

    const isAuthError = body.code === 401 || body.code === 402;
    if (isAuthError && !config._retried) {
      let initiated = false;
      if (!refreshPromise) {
        initiated = true;
        refreshPromise = refreshAccessToken().finally(() => {
          refreshPromise = null;
        });
      }
      return refreshPromise.then((refreshed) => {
        if (refreshed) {
          return http.request({ ...config, _retried: true } as AxiosRequestConfig);
        }
        if (initiated) {
          redirectToLogin();
          message.error('登录已过期，请重新登录');
        }
        return Promise.reject(new Error('登录已过期，请重新登录'));
      });
    }

    message.error(body.message || '请求失败');
    return Promise.reject(new Error(body.message || '请求失败'));
  },
  (error) => {
    message.error(error?.message || '网络异常，请稍后重试');
    return Promise.reject(error);
  },
);

/**
 * 与 services 约定：返回后端 response body（即 { code, message, data } 包体），
 * 调用方通过 res.data 取业务数据。
 */
async function request<T = unknown>(url: string, options?: RequestOptions): Promise<T> {
  const res = await http.request<ApiEnvelope<T>>({ url, ...options });
  return res.data as T;
}

export default request;
