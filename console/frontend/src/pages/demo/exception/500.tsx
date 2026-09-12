import { Button, Result } from 'antd';

import { history } from '@/utils/history';

export default function Exception500() {
  return (
    <Result
      status="500"
      title="500"
      subTitle="抱歉，服务器开小差了，请稍后重试。"
      extra={
        <Button type="primary" onClick={() => history.push('/')}>
          返回首页
        </Button>
      }
    />
  );
}
