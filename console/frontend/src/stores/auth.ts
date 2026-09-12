import { create } from 'zustand';

const TOKEN_KEY = 'adminbase_token';
const REFRESH_TOKEN_KEY = 'adminbase_refresh_token';

interface AuthState {
  token: string;
  refreshToken: string;
  setTokens: (token: string, refreshToken: string) => void;
  clear: () => void;
}

function read(key: string): string {
  try {
    return localStorage.getItem(key) || '';
  } catch {
    return '';
  }
}

export const useAuthStore = create<AuthState>((set) => ({
  token: read(TOKEN_KEY),
  refreshToken: read(REFRESH_TOKEN_KEY),
  setTokens: (token, refreshToken) => {
    try {
      localStorage.setItem(TOKEN_KEY, token);
      localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
    } catch {
      // ignore storage errors
    }
    set({ token, refreshToken });
  },
  clear: () => {
    try {
      localStorage.removeItem(TOKEN_KEY);
      localStorage.removeItem(REFRESH_TOKEN_KEY);
    } catch {
      // ignore storage errors
    }
    set({ token: '', refreshToken: '' });
  },
}));

export const getToken = () => useAuthStore.getState().token;
export const getRefreshToken = () => useAuthStore.getState().refreshToken;
export const clearToken = () => useAuthStore.getState().clear();
