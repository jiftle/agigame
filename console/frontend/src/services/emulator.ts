import request from '@/api/client';

import type { ApiResult } from './types';

export interface EmuSessionInfo {
  id: string;
  console: string;
  cart: string;
  game: string;
  mode: string;
  auto: boolean;
  paused: boolean;
  frames: number;
  width: number;
  height: number;
  createdAt: string;
}

export interface EmuStartInput {
  rom?: string;
  console?: string;
  game?: string;
  mode?: string;
  palette?: string;
}

export interface EmuConfigInput {
  auto?: boolean;
  mode?: string;
  palette?: string;
}

export interface EmuRomInfo {
  name: string;
  title: string;
  console: string;
  ext: string;
  size: number;
}

export async function queryRomList(): Promise<EmuRomInfo[]> {
  const res = await request<ApiResult<{ list: EmuRomInfo[] }>>('/emu/rom/list');
  return res.data.list || [];
}

export async function querySessionList(): Promise<EmuSessionInfo[]> {
  const res = await request<ApiResult<{ list: EmuSessionInfo[] }>>('/emu/session/list');
  return res.data.list || [];
}

export async function startSession(data: EmuStartInput): Promise<EmuSessionInfo> {
  const res = await request<ApiResult<{ session: EmuSessionInfo }>>('/emu/session/start', {
    method: 'POST',
    data,
  });
  return res.data.session;
}

export async function getSession(id: string): Promise<EmuSessionInfo> {
  const res = await request<ApiResult<{ session: EmuSessionInfo }>>(`/emu/session/${id}`);
  return res.data.session;
}

export async function stopSession(id: string): Promise<void> {
  await request(`/emu/session/${id}/stop`, { method: 'POST' });
}

export async function controlSession(id: string, action: string): Promise<void> {
  await request(`/emu/session/${id}/control`, { method: 'POST', data: { action } });
}

export async function configSession(id: string, data: EmuConfigInput): Promise<EmuSessionInfo> {
  const res = await request<ApiResult<{ session: EmuSessionInfo }>>(
    `/emu/session/${id}/config`,
    { method: 'POST', data },
  );
  return res.data.session;
}
