import { PageContainer } from '@ant-design/pro-components';
import {
  Alert,
  Button,
  Card,
  Descriptions,
  Empty,
  Result,
  Segmented,
  Skeleton,
  Space,
  Tag,
} from 'antd';
import { useState } from 'react';

import { Access, useAccess } from '@/access';
import ErrorBoundary from '@/components/ErrorBoundary';
import { message } from '@/utils/antdApp';

import DemoBlock from '../overview/DemoBlock';

type Phase = 'loading' | 'empty' | 'error' | 'success';

const PERMS = [
  'system:user:add',
  'system:user:edit',
  'system:user:remove',
  'system:role:list',
];

function AsyncStateDemo() {
  const [phase, setPhase] = useState<Phase>('success');

  const load = (next: Phase) => {
    setPhase('loading');
    window.setTimeout(() => setPhase(next), 600);
  };

  return (
    <Space orientation="vertical" style={{ width: '100%' }}>
      <Segmented<Phase>
        value={phase}
        onChange={(value) => load(value)}
        options={[
          { label: '加载中', value: 'loading' },
          { label: '空数据', value: 'empty' },
          { label: '加载失败', value: 'error' },
          { label: '正常', value: 'success' },
        ]}
      />
      <Card size="small">
        {phase === 'loading' && <Skeleton active paragraph={{ rows: 3 }} />}
        {phase === 'empty' && <Empty description="暂无数据" />}
        {phase === 'error' && (
          <Result
            status="error"
            title="加载失败"
            subTitle="请检查网络后重试。"
            extra={
              <Button type="primary" onClick={() => load('success')}>
                重试
              </Button>
            }
          />
        )}
        {phase === 'success' && (
          <Result status="success" title="加载成功" subTitle="数据已成功加载。" />
        )}
      </Card>
    </Space>
  );
}

function AccessDemo() {
  const access = useAccess();

  return (
    <Space orientation="vertical" style={{ width: '100%' }}>
      <Descriptions
        column={{ xs: 1, sm: 2 }}
        bordered
        size="small"
        items={PERMS.map((perm) => ({
          key: perm,
          label: perm,
          children: access[perm] ? (
            <Tag color="success">有权限</Tag>
          ) : (
            <Tag>无权限</Tag>
          ),
        }))}
      />
      <Space wrap>
        <Access
          accessible={access['system:user:add']}
          fallback={<Button disabled>新增用户（无权限）</Button>}
        >
          <Button type="primary">新增用户（有权限）</Button>
        </Access>
        <Access accessible={access['system:user:edit']}>
          <Button>编辑用户</Button>
        </Access>
      </Space>
      <Alert
        type="info"
        showIcon
        title="页面级权限在路由层用 guard('system:xxx:list', <Page />)（内部为 <Permission perm>），按钮级用上面的 useAccess + <Access>。"
      />
    </Space>
  );
}

function Crash({ crash }: { crash: boolean }) {
  if (crash) {
    throw new Error('示例：页面渲染异常');
  }
  return <Alert type="success" showIcon title="组件运行正常" />;
}

function BoundaryDemo() {
  const [crash, setCrash] = useState(false);

  return (
    <Space orientation="vertical" style={{ width: '100%' }}>
      <Space>
        <Button danger onClick={() => setCrash(true)}>
          触发渲染异常
        </Button>
        <Button onClick={() => setCrash(false)}>重置</Button>
      </Space>
      <ErrorBoundary key={String(crash)}>
        <Crash crash={crash} />
      </ErrorBoundary>
      <Alert
        type="warning"
        showIcon
        title="ErrorBoundary 仅捕获渲染期异常；全局边界位于 App 根，此处用局部边界演示，避免影响整站。"
      />
    </Space>
  );
}

export default function DemoStatus() {
  return (
    <PageContainer
      title="状态与反馈"
      subTitle="验证异步三态、权限控制与异常边界，仅开发环境可见"
    >
      <DemoBlock title="异步三态（loading / empty / error / success）">
        <AsyncStateDemo />
      </DemoBlock>
      <DemoBlock title="权限演示（useAccess / Access）">
        <AccessDemo />
      </DemoBlock>
      <DemoBlock title="异常边界（ErrorBoundary）">
        <BoundaryDemo />
      </DemoBlock>
      <DemoBlock title="无权限页面预览">
        <Result
          status="403"
          title="403"
          subTitle="抱歉，你无权访问该页面（可由路由 Permission 触发）。"
          extra={
            <Button type="primary" onClick={() => message.info('此处仅作展示')}>
              返回首页
            </Button>
          }
        />
      </DemoBlock>
    </PageContainer>
  );
}
