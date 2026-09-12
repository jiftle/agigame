import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import {
  ModalForm,
  PageContainer,
  ProFormDigit,
  ProFormRadio,
  ProFormText,
  ProFormTextArea,
  ProFormTreeSelect,
  ProTable,
} from '@ant-design/pro-components';
import { Access, useAccess } from '@/access';
import { Button, Popconfirm, Tag } from 'antd';
import { useRef, useState } from 'react';

import { message } from '@/utils/antdApp';

import type { SysDept } from '@/services/system';
import {
  createDept,
  deleteDept,
  getDept,
  queryDeptList,
  updateDept,
} from '@/services/system';

const toTreeData = (list: SysDept[]): any[] =>
  (list || []).map((d) => ({
    title: d.name,
    value: d.id,
    children: d.children ? toTreeData(d.children) : undefined,
  }));

export default function DeptPage() {
  const actionRef = useRef<ActionType | undefined>(undefined);
  const access = useAccess();
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<SysDept | undefined>();
  const [tree, setTree] = useState<SysDept[]>([]);

  const reload = () => actionRef.current?.reload();

  const loadTree = async () => {
    const list = await queryDeptList();
    setTree(list);
    return list;
  };

  const handleEdit = async (record: SysDept) => {
    const dept = await getDept(record.id as number);
    setEditing(dept);
    setOpen(true);
  };

  const columns: ProColumns<SysDept>[] = [
    { title: '部门名称', dataIndex: 'name' },
    { title: '负责人', dataIndex: 'leader', search: false },
    { title: '联系电话', dataIndex: 'phone', search: false },
    { title: '排序', dataIndex: 'sort', search: false },
    {
      title: '状态',
      dataIndex: 'status',
      search: false,
      render: (_, record) =>
        record.status === 1 ? (
          <Tag color="success">正常</Tag>
        ) : (
          <Tag color="error">停用</Tag>
        ),
    },
    { title: '创建时间', dataIndex: 'createdAt', valueType: 'dateTime', search: false },
    {
      title: '操作',
      valueType: 'option',
      width: 140,
      render: (_, record) => [
        <Access key="edit" accessible={access['system:dept:edit']}>
          <a onClick={() => handleEdit(record)}>编辑</a>
        </Access>,
        <Access key="remove" accessible={access['system:dept:remove']}>
          <Popconfirm
            key="remove"
            title="确认删除该部门？"
            onConfirm={async () => {
              await deleteDept(record.id as number);
              message.success('删除成功');
              reload();
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
      <ProTable<SysDept>
        headerTitle="部门列表"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        pagination={false}
        search={false}
        request={async () => {
          const list = await loadTree();
          return { data: list, success: true };
        }}
        toolBarRender={() => [
          <Access key="add" accessible={access['system:dept:add']}>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => {
                setEditing(undefined);
                setOpen(true);
              }}
            >
              新增
            </Button>
          </Access>,
        ]}
      />

      <ModalForm<SysDept>
        key={editing?.id ?? 'create'}
        title={editing ? '编辑部门' : '新增部门'}
        open={open}
        onOpenChange={setOpen}
        width={520}
        initialValues={editing || { parentId: 0, status: 1, sort: 0 }}
        modalProps={{ destroyOnClose: true }}
        onFinish={async (values) => {
          if (editing) {
            await updateDept({ ...values, id: editing.id });
            message.success('修改成功');
          } else {
            await createDept(values);
            message.success('新增成功');
          }
          reload();
          return true;
        }}
      >
        <ProFormTreeSelect
          name="parentId"
          label="上级部门"
          fieldProps={{
            treeData: [{ title: '顶级部门', value: 0, children: toTreeData(tree) }],
            treeDefaultExpandAll: true,
          }}
        />
        <ProFormText
          name="name"
          label="部门名称"
          rules={[{ required: true, message: '请输入部门名称' }]}
        />
        <ProFormText name="leader" label="负责人" />
        <ProFormText name="phone" label="联系电话" />
        <ProFormText name="email" label="邮箱" />
        <ProFormDigit name="sort" label="排序" fieldProps={{ precision: 0 }} />
        <ProFormRadio.Group
          name="status"
          label="状态"
          options={[
            { label: '正常', value: 1 },
            { label: '停用', value: 0 },
          ]}
        />
        <ProFormTextArea name="remark" label="备注" />
      </ModalForm>
    </PageContainer>
  );
}
