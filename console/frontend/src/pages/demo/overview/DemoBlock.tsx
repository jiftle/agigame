import { Card, Space } from 'antd';
import type { ReactNode } from 'react';

interface Props {
  title: string;
  children: ReactNode;
}

/** 组件总览页的区块容器：统一标题与间距 */
export default function DemoBlock({ title, children }: Props) {
  return (
    <Card title={title} size="small" style={{ marginBottom: 16 }}>
      <Space orientation="vertical" size="middle" style={{ width: '100%' }}>
        {children}
      </Space>
    </Card>
  );
}
