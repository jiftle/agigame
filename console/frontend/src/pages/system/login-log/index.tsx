import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { Access, useAccess } from '@/access';
import { Button, Popconfirm, Tag } from 'antd';
import { useRef } from 'react';

import { message } from '@/utils/antdApp';

import type { LoginLog } from '@/services/system';
import { clearLoginLog, deleteLoginLog, queryLoginLogList } from '@/services/system';

export default function LoginLogPage() {
  const actionRef = useRef<ActionType | undefined>(undefined);
  const access = useAccess();

  const columns: ProColumns<LoginLog>[] = [
    { title: 'ID', dataIndex: 'id', width: 70, search: false },
    { title: '用户名', dataIndex: 'username' },
    { title: 'IP', dataIndex: 'ip', search: false },
    { title: '浏览器', dataIndex: 'browser', search: false },
    { title: '操作系统', dataIndex: 'os', search: false },
    {
      title: '状态',
      dataIndex: 'status',
      valueType: 'select',
      valueEnum: {
        1: { text: '成功', status: 'Success' },
        0: { text: '失败', status: 'Error' },
      },
      render: (_, record) =>
        record.status === 1 ? (
          <Tag color="success">成功</Tag>
        ) : (
          <Tag color="error">失败</Tag>
        ),
    },
    { title: '描述', dataIndex: 'msg', search: false },
    { title: '登录时间', dataIndex: 'loginTime', valueType: 'dateTime', search: false },
    {
      title: '操作',
      valueType: 'option',
      width: 80,
      render: (_, record) => [
        <Access key="remove" accessible={access['system:loginlog:remove']}>
          <Popconfirm
            key="remove"
            title="确认删除该日志？"
            onConfirm={async () => {
              await deleteLoginLog([record.id as number]);
              message.success('删除成功');
              actionRef.current?.reload();
            }}
          >
            <a style={{ color: '#ff4d4f' }}>删除</a>
          </Popconfirm>
        </Access>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<LoginLog>
        headerTitle="登录日志"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        search={{ labelWidth: 'auto' }}
        request={async (params) => {
          const res = await queryLoginLogList({
            pageNum: params.current,
            pageSize: params.pageSize,
            username: params.username,
            status: params.status,
          });
          return { data: res.list || [], total: res.total || 0, success: true };
        }}
        toolBarRender={() => [
          <Access key="clear" accessible={access['system:loginlog:remove']}>
            <Popconfirm
              key="clear"
              title="确认清空全部登录日志？"
              onConfirm={async () => {
                await clearLoginLog();
                message.success('清空成功');
                actionRef.current?.reload();
              }}
            >
              <Button danger>清空</Button>
            </Popconfirm>
          </Access>,
        ]}
      />
    </PageContainer>
  );
}
