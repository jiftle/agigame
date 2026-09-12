import request from '@/api/client';

import type { ApiResult, CurrentUser, MenuNode } from './types';

export interface LoginParams {
  username: string;
  password: string;
}

export interface LoginResult {
  token: string;
  refreshToken: string;
  expiresIn: number;
  tokenType: string;
}

export interface UserInfoResult {
  user: CurrentUser;
  roles: string[];
  perms: string[];
  menus: MenuNode[];
}

export async function login(body: LoginParams): Promise<LoginResult> {
  const res = await request<ApiResult<LoginResult>>('/auth/login', {
    method: 'POST',
    data: body,
  });
  return res.data;
}

export async function refreshToken(token: string): Promise<LoginResult> {
  const res = await request<ApiResult<LoginResult>>('/auth/refresh', {
    method: 'POST',
    data: { refreshToken: token },
  });
  return res.data;
}

export async function getUserInfo(): Promise<UserInfoResult> {
  const res = await request<ApiResult<UserInfoResult>>('/auth/user-info', {
    method: 'GET',
  });
  return res.data;
}

export async function logout(refreshToken?: string): Promise<void> {
  await request('/auth/logout', { method: 'POST', data: { refreshToken } });
}

export async function changePassword(data: {
  oldPassword: string;
  newPassword: string;
}): Promise<void> {
  await request('/auth/password', { method: 'PUT', data });
}

export interface UpdateProfileParams {
  nickname: string;
  email?: string;
  phone?: string;
  avatar?: string;
}

export async function updateProfile(data: UpdateProfileParams): Promise<void> {
  await request('/auth/profile', { method: 'PUT', data });
}
