import { DashboardOutlined, ExperimentOutlined, SettingOutlined, VideoCameraOutlined } from '@ant-design/icons';
import type { ReactNode } from 'react';

export interface MenuItemConfig {
  path: string;
  name: string;
  icon?: ReactNode;
  access?: string;
  children?: MenuItemConfig[];
}

const systemMenu: MenuItemConfig[] = [
  { path: '/system/user', name: '用户管理', access: 'system:user:list' },
  { path: '/system/role', name: '角色管理', access: 'system:role:list' },
  { path: '/system/menu', name: '菜单管理', access: 'system:menu:list' },
  { path: '/system/dept', name: '部门管理', access: 'system:dept:list' },
  { path: '/system/dict', name: '字典管理', access: 'system:dict:list' },
  { path: '/system/config', name: '参数配置', access: 'system:config:list' },
  { path: '/system/oper-log', name: '操作日志', access: 'system:operlog:list' },
  { path: '/system/login-log', name: '登录日志', access: 'system:loginlog:list' },
];

const demoMenu: MenuItemConfig[] = [
  { path: '/demo/form', name: '表单页' },
  { path: '/demo/table', name: '查询表格' },
  { path: '/demo/detail', name: '详情页' },
  { path: '/demo/result', name: '结果页' },
  { path: '/demo/chart', name: '图表示例' },
  { path: '/demo/overview', name: '组件总览' },
  { path: '/demo/status', name: '状态与反馈' },
];

export const menuConfig: MenuItemConfig[] = [
  { path: '/dashboard', name: '仪表盘', icon: <DashboardOutlined /> },
  {
    path: '/emulator',
    name: '模拟器',
    icon: <VideoCameraOutlined />,
    children: [{ path: '/emulator/sessions', name: '会话管理', access: 'emu:session:list' }],
  },
  { path: '/system', name: '系统管理', icon: <SettingOutlined />, children: systemMenu },
  ...(import.meta.env.DEV
    ? [{ path: '/demo', name: '示例页面', icon: <ExperimentOutlined />, children: demoMenu }]
    : []),
];
