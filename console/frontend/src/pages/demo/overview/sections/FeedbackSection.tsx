import { Button, Drawer, Empty, Popconfirm, Progress, Result, Skeleton, Space, Spin, Alert } from 'antd';
import { useState } from 'react';

import { message, modal, notification } from '@/utils/antdApp';

import DemoBlock from '../DemoBlock';

export default function FeedbackSection() {
  const [drawerOpen, setDrawerOpen] = useState(false);

  return (
    <>
      <DemoBlock title="Alert 警告提示">
        <Alert type="success" showIcon title="成功提示" />
        <Alert type="info" showIcon title="信息提示" description="附带说明文字的信息提示。" />
        <Alert type="warning" showIcon title="警告提示" />
        <Alert type="error" showIcon title="错误提示" closable />
      </DemoBlock>

      <DemoBlock title="Progress 进度条 / Spin / Skeleton">
        <Space size="large" wrap>
          <Progress type="line" percent={60} style={{ width: 240 }} />
          <Progress type="circle" percent={75} size={80} />
          <Progress type="dashboard" percent={40} size={80} />
        </Space>
        <Space size="large" wrap>
          <Progress percent={100} status="success" style={{ width: 200 }} />
          <Progress percent={70} status="exception" style={{ width: 200 }} />
        </Space>
        <Space size="large" wrap>
          <Spin />
          <Spin size="large" />
          <Spin description="加载中..." />
        </Space>
        <Skeleton active paragraph={{ rows: 2 }} />
      </DemoBlock>

      <DemoBlock title="Empty 空状态">
        <Space size="large" wrap>
          <Empty description="暂无数据" />
          <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description="简洁空态" />
        </Space>
      </DemoBlock>

      <DemoBlock title="Result 结果页">
        <Space size="large" wrap>
          <Result status="success" title="提交成功" subTitle="操作已完成。" style={{ padding: 0 }} />
          <Result status="warning" title="存在风险" style={{ padding: 0 }} />
          <Result status="error" title="提交失败" style={{ padding: 0 }} />
        </Space>
      </DemoBlock>

      <DemoBlock title="全局反馈（message / notification / modal）">
        <Space wrap>
          <Button onClick={() => message.success('操作成功')}>message.success</Button>
          <Button onClick={() => message.error('操作失败')}>message.error</Button>
          <Button onClick={() => message.loading('加载中...', 1)}>message.loading</Button>
          <Button
            onClick={() =>
              notification.success({ message: '通知标题', description: '通知内容' })
            }
          >
            notification
          </Button>
          <Button
            type="primary"
            onClick={() =>
              modal.confirm({
                title: '确认操作',
                content: '这是由 AntdApp 上下文弹出的确认框。',
                onOk: () => message.success('已确认'),
              })
            }
          >
            modal.confirm
          </Button>
        </Space>
      </DemoBlock>

      <DemoBlock title="Drawer / Popconfirm">
        <Space wrap>
          <Button onClick={() => setDrawerOpen(true)}>打开 Drawer</Button>
          <Popconfirm title="确认删除？" okText="确认" cancelText="取消" onConfirm={() => message.success('已删除')}>
            <Button danger>Popconfirm</Button>
          </Popconfirm>
        </Space>
        <Drawer
          title="抽屉标题"
          open={drawerOpen}
          onClose={() => setDrawerOpen(false)}
          placement="right"
        >
          <p>抽屉内容，用于承载表单或详情。</p>
        </Drawer>
      </DemoBlock>
    </>
  );
}
