import { Button, Result } from 'antd';
import type { ReactNode } from 'react';

import { useAccess } from '@/access';
import { history } from '@/utils/history';

/** 页面级权限：无权限时展示 403 */
export default function Permission({
  perm,
  children,
}: {
  perm: string;
  children: ReactNode;
}) {
  const access = useAccess();
  if (access[perm]) {
    return <>{children}</>;
  }
  return (
    <Result
      status="403"
      title="403"
      subTitle="抱歉，您没有访问该页面的权限。"
      extra={
        <Button type="primary" onClick={() => history.push('/')}>
          返回首页
        </Button>
      }
    />
  );
}
