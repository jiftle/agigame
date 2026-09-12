import request from '@/api/client';

import type { ApiResult, PageResult } from './types';

/* ============================ 类型定义 ============================ */

export interface SysUser {
  id?: number;
  deptId?: number;
  username?: string;
  password?: string;
  nickname?: string;
  email?: string;
  phone?: string;
  sex?: number;
  avatar?: string;
  status?: number;
  loginIp?: string;
  loginDate?: string;
  remark?: string;
  createdAt?: string;
  roleIds?: number[];
}

export interface SysRole {
  id?: number;
  name?: string;
  code?: string;
  sort?: number;
  dataScope?: number;
  status?: number;
  remark?: string;
  menuIds?: number[];
  createdAt?: string;
}

export interface SysMenu {
  id?: number;
  parentId?: number;
  title?: string;
  name?: string;
  path?: string;
  component?: string;
  icon?: string;
  type?: string;
  perms?: string;
  sort?: number;
  visible?: number;
  status?: number;
  redirect?: string;
  isFrame?: number;
  isCache?: number;
  children?: SysMenu[];
}

export interface SysDept {
  id?: number;
  parentId?: number;
  name?: string;
  leader?: string;
  phone?: string;
  email?: string;
  sort?: number;
  status?: number;
  remark?: string;
  children?: SysDept[];
}

export interface DictType {
  id?: number;
  name?: string;
  type?: string;
  status?: number;
  remark?: string;
}

export interface DictData {
  id?: number;
  dictSort?: number;
  dictLabel?: string;
  dictValue?: string;
  dictType?: string;
  isDefault?: number;
  status?: number;
  remark?: string;
}

export interface SysConfig {
  id?: number;
  configName?: string;
  configKey?: string;
  configValue?: string;
  configType?: number;
  remark?: string;
}

export interface LoginLog {
  id?: number;
  username?: string;
  ip?: string;
  location?: string;
  browser?: string;
  os?: string;
  status?: number;
  msg?: string;
  loginTime?: string;
}

export interface OperLog {
  id?: number;
  title?: string;
  businessType?: string;
  requestMethod?: string;
  operName?: string;
  operUrl?: string;
  operIp?: string;
  operParam?: string;
  status?: number;
  errorMsg?: string;
  cost?: number;
  operTime?: string;
}

/* ============================ 用户 ============================ */

export async function queryUserList(params: any): Promise<PageResult<SysUser>> {
  const res = await request<ApiResult<PageResult<SysUser>>>('/system/user/list', { params });
  return res.data;
}
export async function getUser(id: number): Promise<{ user: SysUser; roleIds: number[] }> {
  const res = await request<ApiResult<{ user: SysUser; roleIds: number[] }>>(`/system/user/${id}`);
  return res.data;
}
export async function createUser(data: SysUser): Promise<void> {
  await request('/system/user', { method: 'POST', data });
}
export async function updateUser(data: SysUser): Promise<void> {
  await request('/system/user', { method: 'PUT', data });
}
export async function deleteUser(ids: number[]): Promise<void> {
  await request('/system/user', { method: 'DELETE', data: { ids } });
}
export async function resetUserPwd(id: number, password?: string): Promise<void> {
  await request('/system/user/resetPwd', { method: 'PUT', data: { id, password } });
}
export async function changeUserStatus(id: number, status: number): Promise<void> {
  await request('/system/user/status', { method: 'PUT', data: { id, status } });
}

/* ============================ 角色 ============================ */

export async function queryRoleList(params: any): Promise<PageResult<SysRole>> {
  const res = await request<ApiResult<PageResult<SysRole>>>('/system/role/list', { params });
  return res.data;
}
export async function getRole(id: number): Promise<{ role: SysRole; menuIds: number[] }> {
  const res = await request<ApiResult<{ role: SysRole; menuIds: number[] }>>(`/system/role/${id}`);
  return res.data;
}
export async function createRole(data: SysRole): Promise<void> {
  await request('/system/role', { method: 'POST', data });
}
export async function updateRole(data: SysRole): Promise<void> {
  await request('/system/role', { method: 'PUT', data });
}
export async function deleteRole(ids: number[]): Promise<void> {
  await request('/system/role', { method: 'DELETE', data: { ids } });
}

/* ============================ 菜单 ============================ */

export async function queryMenuList(params?: any): Promise<SysMenu[]> {
  const res = await request<ApiResult<{ list: SysMenu[] }>>('/system/menu/list', { params });
  return res.data.list || [];
}
export async function getMenu(id: number): Promise<SysMenu> {
  const res = await request<ApiResult<{ menu: SysMenu }>>(`/system/menu/${id}`);
  return res.data.menu;
}
export async function createMenu(data: SysMenu): Promise<void> {
  await request('/system/menu', { method: 'POST', data });
}
export async function updateMenu(data: SysMenu): Promise<void> {
  await request('/system/menu', { method: 'PUT', data });
}
export async function deleteMenu(id: number): Promise<void> {
  await request(`/system/menu/${id}`, { method: 'DELETE' });
}

/* ============================ 部门 ============================ */

export async function queryDeptList(params?: any): Promise<SysDept[]> {
  const res = await request<ApiResult<{ list: SysDept[] }>>('/system/dept/list', { params });
  return res.data.list || [];
}
export async function getDept(id: number): Promise<SysDept> {
  const res = await request<ApiResult<{ dept: SysDept }>>(`/system/dept/${id}`);
  return res.data.dept;
}
export async function createDept(data: SysDept): Promise<void> {
  await request('/system/dept', { method: 'POST', data });
}
export async function updateDept(data: SysDept): Promise<void> {
  await request('/system/dept', { method: 'PUT', data });
}
export async function deleteDept(id: number): Promise<void> {
  await request(`/system/dept/${id}`, { method: 'DELETE' });
}

/* ============================ 字典 ============================ */

export async function queryDictTypeList(params: any): Promise<PageResult<DictType>> {
  const res = await request<ApiResult<PageResult<DictType>>>('/system/dict/type/list', { params });
  return res.data;
}
export async function getDictType(id: number): Promise<DictType> {
  const res = await request<ApiResult<{ dictType: DictType }>>(`/system/dict/type/${id}`);
  return res.data.dictType;
}
export async function createDictType(data: DictType): Promise<void> {
  await request('/system/dict/type', { method: 'POST', data });
}
export async function updateDictType(data: DictType): Promise<void> {
  await request('/system/dict/type', { method: 'PUT', data });
}
export async function deleteDictType(ids: number[]): Promise<void> {
  await request('/system/dict/type', { method: 'DELETE', data: { ids } });
}

export async function queryDictDataList(params: any): Promise<PageResult<DictData>> {
  const res = await request<ApiResult<PageResult<DictData>>>('/system/dict/data/list', { params });
  return res.data;
}
export async function getDictDataByType(dictType: string): Promise<DictData[]> {
  const res = await request<ApiResult<{ list: DictData[] }>>(`/system/dict/data/type/${dictType}`);
  return res.data.list || [];
}
export async function getDictData(id: number): Promise<DictData> {
  const res = await request<ApiResult<{ dictData: DictData }>>(`/system/dict/data/${id}`);
  return res.data.dictData;
}
export async function createDictData(data: DictData): Promise<void> {
  await request('/system/dict/data', { method: 'POST', data });
}
export async function updateDictData(data: DictData): Promise<void> {
  await request('/system/dict/data', { method: 'PUT', data });
}
export async function deleteDictData(ids: number[]): Promise<void> {
  await request('/system/dict/data', { method: 'DELETE', data: { ids } });
}

/* ============================ 参数 ============================ */

export async function queryConfigList(params: any): Promise<PageResult<SysConfig>> {
  const res = await request<ApiResult<PageResult<SysConfig>>>('/system/config/list', { params });
  return res.data;
}
export async function getConfig(id: number): Promise<SysConfig> {
  const res = await request<ApiResult<{ config: SysConfig }>>(`/system/config/${id}`);
  return res.data.config;
}
export async function createConfig(data: SysConfig): Promise<void> {
  await request('/system/config', { method: 'POST', data });
}
export async function updateConfig(data: SysConfig): Promise<void> {
  await request('/system/config', { method: 'PUT', data });
}
export async function deleteConfig(ids: number[]): Promise<void> {
  await request('/system/config', { method: 'DELETE', data: { ids } });
}

/* ============================ 日志 ============================ */

export async function queryOperLogList(params: any): Promise<PageResult<OperLog>> {
  const res = await request<ApiResult<PageResult<OperLog>>>('/system/log/oper/list', { params });
  return res.data;
}
export async function deleteOperLog(ids: number[]): Promise<void> {
  await request('/system/log/oper', { method: 'DELETE', data: { ids } });
}
export async function clearOperLog(): Promise<void> {
  await request('/system/log/oper/clear', { method: 'DELETE' });
}
export async function queryLoginLogList(params: any): Promise<PageResult<LoginLog>> {
  const res = await request<ApiResult<PageResult<LoginLog>>>('/system/log/login/list', { params });
  return res.data;
}
export async function deleteLoginLog(ids: number[]): Promise<void> {
  await request('/system/log/login', { method: 'DELETE', data: { ids } });
}
export async function clearLoginLog(): Promise<void> {
  await request('/system/log/login/clear', { method: 'DELETE' });
}
