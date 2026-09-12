import type { ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router-dom';

import GlobalLoading from '@/loading';
import { useCurrentUser } from '@/hooks/useCurrentUser';
import { useAuthStore } from '@/stores/auth';

/** 登录守卫：无 token 或用户信息拉取失败时跳转登录 */
export default function RequireAuth({ children }: { children: ReactNode }) {
  const token = useAuthStore((state) => state.token);
  const location = useLocation();
  const { data, isLoading, isError } = useCurrentUser();

  if (!token) {
    const redirect = encodeURIComponent(location.pathname + location.search);
    return <Navigate to={`/login?redirect=${redirect}`} replace />;
  }
  if (isLoading) {
    return <GlobalLoading />;
  }
  if (isError || !data) {
    return <Navigate to="/login" replace />;
  }
  return <>{children}</>;
}
