import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { Access, useAccess } from '@/access';
import { Button, Popconfirm, Tag } from 'antd';
import { useRef } from 'react';

import { message } from '@/utils/antdApp';

import type { OperLog } from '@/services/system';
import { clearOperLog, deleteOperLog, queryOperLogList } from '@/services/system';

export default function OperLogPage() {
  const actionRef = useRef<ActionType | undefined>(undefined);
  const access = useAccess();

  const columns: ProColumns<OperLog>[] = [
    { title: 'ID', dataIndex: 'id', width: 70, search: false },
    { title: '操作模块', dataIndex: 'title' },
    { title: '操作人', dataIndex: 'operName' },
    { title: '请求方式', dataIndex: 'requestMethod', search: false },
    { title: '请求地址', dataIndex: 'operUrl', search: false, ellipsis: true },
    { title: 'IP', dataIndex: 'operIp', search: false },
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
    { title: '耗时(ms)', dataIndex: 'cost', search: false },
    { title: '操作时间', dataIndex: 'operTime', valueType: 'dateTime', search: false },
    {
      title: '操作',
      valueType: 'option',
      width: 80,
      render: (_, record) => [
        <Access key="remove" accessible={access['system:operlog:remove']}>
          <Popconfirm
            key="remove"
            title="确认删除该日志？"
            onConfirm={async () => {
              await deleteOperLog([record.id as number]);
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
      <ProTable<OperLog>
        headerTitle="操作日志"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        search={{ labelWidth: 'auto' }}
        request={async (params) => {
          const res = await queryOperLogList({
            pageNum: params.current,
            pageSize: params.pageSize,
            title: params.title,
            operName: params.operName,
            status: params.status,
          });
          return { data: res.list || [], total: res.total || 0, success: true };
        }}
        toolBarRender={() => [
          <Access key="clear" accessible={access['system:operlog:remove']}>
            <Popconfirm
              key="clear"
              title="确认清空全部操作日志？"
              onConfirm={async () => {
                await clearOperLog();
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
