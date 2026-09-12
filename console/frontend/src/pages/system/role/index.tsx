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
import { Button, Popconfirm } from 'antd';
import { useEffect, useRef, useState } from 'react';

import { message } from '@/utils/antdApp';

import type { SysMenu, SysRole } from '@/services/system';
import {
  createRole,
  deleteRole,
  getRole,
  queryMenuList,
  queryRoleList,
  updateRole,
} from '@/services/system';

const toTreeData = (list: SysMenu[]): any[] =>
  (list || []).map((m) => ({
    title: m.title,
    key: m.id,
    value: m.id,
    children: m.children ? toTreeData(m.children) : undefined,
  }));

export default function RolePage() {
  const actionRef = useRef<ActionType | undefined>(undefined);
  const access = useAccess();
  const [menuTree, setMenuTree] = useState<any[]>([]);
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<SysRole | undefined>();

  useEffect(() => {
    queryMenuList().then((list) => setMenuTree(toTreeData(list)));
  }, []);

  const reload = () => actionRef.current?.reload();

  const handleEdit = async (record: SysRole) => {
    const detail = await getRole(record.id as number);
    setEditing({ ...detail.role, menuIds: detail.menuIds });
    setOpen(true);
  };

  const columns: ProColumns<SysRole>[] = [
    { title: 'ID', dataIndex: 'id', width: 70, search: false },
    { title: '角色名称', dataIndex: 'name' },
    { title: '角色编码', dataIndex: 'code' },
    { title: '排序', dataIndex: 'sort', search: false },
    {
      title: '状态',
      dataIndex: 'status',
      valueType: 'select',
      valueEnum: {
        1: { text: '正常', status: 'Success' },
        0: { text: '停用', status: 'Error' },
      },
    },
    { title: '备注', dataIndex: 'remark', search: false, ellipsis: true },
    {
      title: '操作',
      valueType: 'option',
      width: 140,
      render: (_, record) => [
        <Access key="edit" accessible={access['system:role:edit']}>
          <a onClick={() => handleEdit(record)}>编辑</a>
        </Access>,
        <Access key="remove" accessible={access['system:role:remove']}>
          <Popconfirm
            key="remove"
            title="确认删除该角色？"
            onConfirm={async () => {
              await deleteRole([record.id as number]);
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
      <ProTable<SysRole>
        headerTitle="角色列表"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        search={{ labelWidth: 'auto' }}
        request={async (params) => {
          const res = await queryRoleList({
            pageNum: params.current,
            pageSize: params.pageSize,
            name: params.name,
            code: params.code,
            status: params.status,
          });
          return { data: res.list || [], total: res.total || 0, success: true };
        }}
        toolBarRender={() => [
          <Access key="add" accessible={access['system:role:add']}>
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

      <ModalForm<SysRole>
        key={editing?.id ?? 'create'}
        title={editing ? '编辑角色' : '新增角色'}
        open={open}
        onOpenChange={setOpen}
        width={560}
        initialValues={editing || { status: 1, sort: 0, dataScope: 1, menuIds: [] }}
        modalProps={{ destroyOnClose: true }}
        onFinish={async (values) => {
          if (editing) {
            await updateRole({ ...values, id: editing.id });
            message.success('修改成功');
          } else {
            await createRole(values);
            message.success('新增成功');
          }
          reload();
          return true;
        }}
      >
        <ProFormText
          name="name"
          label="角色名称"
          rules={[{ required: true, message: '请输入角色名称' }]}
        />
        <ProFormText
          name="code"
          label="角色编码"
          rules={[{ required: true, message: '请输入角色编码' }]}
        />
        <ProFormDigit name="sort" label="排序" fieldProps={{ precision: 0 }} />
        <ProFormRadio.Group
          name="status"
          label="状态"
          options={[
            { label: '正常', value: 1 },
            { label: '停用', value: 0 },
          ]}
        />
        <ProFormTreeSelect
          name="menuIds"
          label="菜单权限"
          fieldProps={{
            treeData: menuTree,
            multiple: true,
            treeCheckable: true,
            showCheckedStrategy: 'SHOW_PARENT',
            treeDefaultExpandAll: true,
          }}
        />
        <ProFormTextArea name="remark" label="备注" />
      </ModalForm>
    </PageContainer>
  );
}
