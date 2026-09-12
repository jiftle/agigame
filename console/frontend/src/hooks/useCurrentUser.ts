import { useQuery } from '@tanstack/react-query';

import { getUserInfo, type UserInfoResult } from '@/services/auth';
import { useAuthStore } from '@/stores/auth';

/** 当前登录用户信息（角色 / 权限 / 菜单） */
export function useCurrentUser() {
  const token = useAuthStore((state) => state.token);
  return useQuery<UserInfoResult>({
    queryKey: ['auth', 'user'],
    queryFn: getUserInfo,
    enabled: !!token,
    retry: false,
    staleTime: 5 * 60 * 1000,
  });
}
