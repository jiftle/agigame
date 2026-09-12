import { createContext, useContext, type ReactNode } from 'react';

export type AccessMap = Record<string, boolean>;

const ALL_PERMS = [
  'system:user:list',
  'system:user:add',
  'system:user:edit',
  'system:user:remove',
  'system:user:resetPwd',
  'system:role:list',
  'system:role:add',
  'system:role:edit',
  'system:role:remove',
  'system:menu:list',
  'system:menu:add',
  'system:menu:edit',
  'system:menu:remove',
  'system:dept:list',
  'system:dept:add',
  'system:dept:edit',
  'system:dept:remove',
  'system:dict:list',
  'system:dict:add',
  'system:dict:edit',
  'system:dict:remove',
  'system:config:list',
  'system:config:add',
  'system:config:edit',
  'system:config:remove',
  'system:operlog:list',
  'system:operlog:remove',
  'system:loginlog:list',
  'system:loginlog:remove',
];

/** 由后端返回的权限标识生成前端访问开关 */
export function buildAccess(perms: string[] = [], isSuper = false): AccessMap {
  const has = (perm: string) =>
    isSuper || perms.includes('*:*:*') || perms.includes(perm);

  const result: AccessMap = {
    canSystem: has('system:user:list') || has('system:role:list'),
  };
  ALL_PERMS.forEach((perm) => {
    result[perm] = has(perm);
  });
  return result;
}

const AccessContext = createContext<AccessMap>({});

export function AccessProvider({
  access,
  children,
}: {
  access: AccessMap;
  children: ReactNode;
}) {
  return <AccessContext.Provider value={access}>{children}</AccessContext.Provider>;
}

export function useAccess(): AccessMap {
  return useContext(AccessContext);
}

export function Access({
  accessible,
  fallback = null,
  children,
}: {
  accessible?: boolean;
  fallback?: ReactNode;
  children?: ReactNode;
}) {
  return <>{accessible ? children : fallback}</>;
}
