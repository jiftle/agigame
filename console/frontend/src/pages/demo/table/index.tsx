import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { Button, Popconfirm, Tag } from 'antd';
import type { Key } from 'react';
import { useRef, useState } from 'react';

import { message } from '@/utils/antdApp';

interface DemoRow {
  id: number;
  name: string;
  owner: string;
  status: number;
  amount: number;
  createdAt: string;
}

const owners = ['张三', '李四', '王五'];

const genData = (): DemoRow[] =>
  Array.from({ length: 46 }).map((_, i) => ({
    id: i + 1,
    name: `项目 ${i + 1}`,
    owner: owners[i % owners.length],
    status: i % 4 === 0 ? 0 : 1,
    amount: (i + 1) * 1000,
    createdAt: `2026-09-${String((i % 28) + 1).padStart(2, '0')}`,
  }));

export default function DemoTable() {
  const actionRef = useRef<ActionType | undefined>(undefined);
  const dataRef = useRef<DemoRow[]>(genData());
  const [selected, setSelected] = useState<Key[]>([]);

  const handleDelete = (ids: Key[]) => {
    dataRef.current = dataRef.current.filter((r) => !ids.includes(r.id));
    setSelected([]);
    actionRef.current?.reload();
    message.success('删除成功');
  };

  const columns: ProColumns<DemoRow>[] = [
    { title: 'ID', dataIndex: 'id', width: 70, search: false },
    { title: '项目名称', dataIndex: 'name', ellipsis: true },
    {
      title: '负责人',
      dataIndex: 'owner',
      valueType: 'select',
      valueEnum: {
        张三: { text: '张三' },
        李四: { text: '李四' },
        王五: { text: '王五' },
      },
    },
    { title: '预算', dataIndex: 'amount', valueType: 'money', search: false },
    {
      title: '状态',
      dataIndex: 'status',
      valueType: 'select',
      valueEnum: {
        1: { text: '进行中' },
        0: { text: '已归档' },
      },
      render: (_, record) =>
        record.status === 1 ? (
          <Tag color="processing">进行中</Tag>
        ) : (
          <Tag>已归档</Tag>
        ),
    },
    {
      title: '创建时间',
      dataIndex: 'createdAt',
      valueType: 'date',
      search: false,
    },
    {
      title: '操作',
      valueType: 'option',
      width: 90,
      render: (_, record) => [
        <Popconfirm
          key="delete"
          title="确认删除该记录？"
          onConfirm={() => handleDelete([record.id])}
        >
          <a style={{ color: '#ff4d4f' }}>删除</a>
        </Popconfirm>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<DemoRow>
        headerTitle="查询表格"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        search={{ labelWidth: 'auto' }}
        rowSelection={{ selectedRowKeys: selected, onChange: setSelected }}
        request={async (params) => {
          const { current = 1, pageSize = 10 } = params;
          const name = params.name as string | undefined;
          const owner = params.owner as string | undefined;
          const status = params.status as string | number | undefined;

          let list = dataRef.current;
          if (name) {
            list = list.filter((r) => r.name.includes(name));
          }
          if (owner) {
            list = list.filter((r) => r.owner === owner);
          }
          if (status !== undefined && status !== '') {
            list = list.filter((r) => String(r.status) === String(status));
          }
          const start = (current - 1) * pageSize;
          return {
            data: list.slice(start, start + pageSize),
            total: list.length,
            success: true,
          };
        }}
        toolBarRender={() => [
          <Button
            key="add"
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => message.info('演示页面：新增')}
          >
            新增
          </Button>,
          <Popconfirm
            key="batch"
            title={`确认删除选中的 ${selected.length} 项？`}
            disabled={!selected.length}
            onConfirm={() => handleDelete(selected)}
          >
            <Button danger disabled={!selected.length}>
              批量删除
            </Button>
          </Popconfirm>,
        ]}
      />
    </PageContainer>
  );
}
