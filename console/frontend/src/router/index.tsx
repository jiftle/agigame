import { lazy, Suspense, type ReactNode } from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';

import BasicLayout from '@/layouts/BasicLayout';
import GlobalLoading from '@/loading';

import Permission from './Permission';
import RequireAuth from './RequireAuth';

const Login = lazy(() => import('@/pages/Login'));
const Dashboard = lazy(() => import('@/pages/Dashboard'));
const AccountPage = lazy(() => import('@/pages/account'));
const NotFound = lazy(() => import('@/pages/404'));

const UserPage = lazy(() => import('@/pages/system/user'));
const RolePage = lazy(() => import('@/pages/system/role'));
const MenuPage = lazy(() => import('@/pages/system/menu'));
const DeptPage = lazy(() => import('@/pages/system/dept'));
const DictPage = lazy(() => import('@/pages/system/dict'));
const ConfigPage = lazy(() => import('@/pages/system/config'));
const OperLogPage = lazy(() => import('@/pages/system/oper-log'));
const LoginLogPage = lazy(() => import('@/pages/system/login-log'));

const DemoForm = lazy(() => import('@/pages/demo/form'));
const DemoTable = lazy(() => import('@/pages/demo/table'));
const DemoDetail = lazy(() => import('@/pages/demo/detail'));
const DemoResult = lazy(() => import('@/pages/demo/result'));
const DemoChart = lazy(() => import('@/pages/demo/chart'));
const DemoOverview = lazy(() => import('@/pages/demo/overview'));
const DemoStatus = lazy(() => import('@/pages/demo/status'));
const Demo403 = lazy(() => import('@/pages/demo/exception/403'));
const Demo500 = lazy(() => import('@/pages/demo/exception/500'));

const guard = (perm: string, node: ReactNode) => (
  <Permission perm={perm}>{node}</Permission>
);

const demoRoutes = import.meta.env.DEV ? (
  <>
    <Route path="/demo" element={<Navigate to="/demo/form" replace />} />
    <Route path="/demo/form" element={<DemoForm />} />
    <Route path="/demo/table" element={<DemoTable />} />
    <Route path="/demo/detail" element={<DemoDetail />} />
    <Route path="/demo/result" element={<DemoResult />} />
    <Route path="/demo/chart" element={<DemoChart />} />
    <Route path="/demo/overview" element={<DemoOverview />} />
    <Route path="/demo/status" element={<DemoStatus />} />
    <Route path="/demo/exception/403" element={<Demo403 />} />
    <Route path="/demo/exception/500" element={<Demo500 />} />
  </>
) : null;

export default function AppRoutes() {
  return (
    <Suspense fallback={<GlobalLoading />}>
      <Routes>
        <Route path="/login" element={<Login />} />

        <Route
          element={
            <RequireAuth>
              <BasicLayout />
            </RequireAuth>
          }
        >
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/account" element={<AccountPage />} />

          <Route path="/system" element={<Navigate to="/system/user" replace />} />
          <Route path="/system/user" element={guard('system:user:list', <UserPage />)} />
          <Route path="/system/role" element={guard('system:role:list', <RolePage />)} />
          <Route path="/system/menu" element={guard('system:menu:list', <MenuPage />)} />
          <Route path="/system/dept" element={guard('system:dept:list', <DeptPage />)} />
          <Route path="/system/dict" element={guard('system:dict:list', <DictPage />)} />
          <Route path="/system/config" element={guard('system:config:list', <ConfigPage />)} />
          <Route path="/system/oper-log" element={guard('system:operlog:list', <OperLogPage />)} />
          <Route path="/system/login-log" element={guard('system:loginlog:list', <LoginLogPage />)} />

          {demoRoutes}
        </Route>

        <Route path="*" element={<NotFound />} />
      </Routes>
    </Suspense>
  );
}
