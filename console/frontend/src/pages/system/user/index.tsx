import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import {
  ModalForm,
  PageContainer,
  ProFormRadio,
  ProFormSelect,
  ProFormText,
  ProFormTextArea,
  ProFormTreeSelect,
  ProTable,
} from '@ant-design/pro-components';
import { Access, useAccess } from '@/access';
import { Button, Popconfirm, Switch } from 'antd';
import { useEffect, useRef, useState } from 'react';

import { message } from '@/utils/antdApp';

import type { SysDept, SysRole, SysUser } from '@/services/system';
import {
  changeUserStatus,
  createUser,
  deleteUser,
  getUser,
  queryDeptList,
  queryRoleList,
  queryUserList,
  resetUserPwd,
  updateUser,
} from '@/services/system';

const toTreeData = (list: SysDept[]): any[] =>
  (list || []).map((d) => ({
    title: d.name,
    value: d.id,
    children: d.children ? toTreeData(d.children) : undefined,
  }));

const flattenDept = (list: SysDept[], map: Record<number, string> = {}) => {
  (list || []).forEach((d) => {
    if (d.id !== undefined) map[d.id] = d.name || '';
    if (d.children) flattenDept(d.children, map);
  });
  return map;
};

export default function UserPage() {
  const actionRef = useRef<ActionType | undefined>(undefined);
  const access = useAccess();
  const [deptTree, setDeptTree] = useState<any[]>([]);
  const [deptMap, setDeptMap] = useState<Record<number, string>>({});
  const [roles, setRoles] = useState<SysRole[]>([]);
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<SysUser | undefined>();

  useEffect(() => {
    queryDeptList().then((list) => {
      setDeptTree(toTreeData(list));
      setDeptMap(flattenDept(list));
    });
    queryRoleList({ pageNum: 1, pageSize: 100 }).then((res) =>
      setRoles(res.list || []),
    );
  }, []);

  const reload = () => actionRef.current?.reload();

  const handleEdit = async (record: SysUser) => {
    const detail = await getUser(record.id as number);
    setEditing({ ...detail.user, roleIds: detail.roleIds });
    setOpen(true);
  };

  const handleDelete = async (record: SysUser) => {
    await deleteUser([record.id as number]);
    message.success('删除成功');
    reload();
  };

  const columns: ProColumns<SysUser>[] = [
    { title: 'ID', dataIndex: 'id', width: 70, search: false },
    { title: '用户名', dataIndex: 'username' },
    { title: '昵称', dataIndex: 'nickname' },
    {
      title: '部门',
      dataIndex: 'deptId',
      search: false,
      render: (_, record) => deptMap[record.deptId as number] || '-',
    },
    { title: '手机号', dataIndex: 'phone', search: false },
    {
      title: '状态',
      dataIndex: 'status',
      valueType: 'select',
      valueEnum: {
        1: { text: '正常', status: 'Success' },
        0: { text: '停用', status: 'Error' },
      },
      render: (_, record) => (
        <Switch
          checked={record.status === 1}
          disabled={!access['system:user:edit']}
          onChange={async (checked) => {
            await changeUserStatus(record.id as number, checked ? 1 : 0);
            message.success('修改成功');
            reload();
          }}
        />
      ),
    },
    { title: '创建时间', dataIndex: 'createdAt', valueType: 'dateTime', search: false },
    {
      title: '操作',
      valueType: 'option',
      width: 200,
      render: (_, record) => [
        <Access key="edit" accessible={access['system:user:edit']}>
          <a onClick={() => handleEdit(record)}>编辑</a>
        </Access>,
        <Access key="reset" accessible={access['system:user:resetPwd']}>
          <Popconfirm
            key="reset"
            title="重置该用户密码为 123456 ？"
            onConfirm={async () => {
              await resetUserPwd(record.id as number);
              message.success('重置成功');
            }}
          >
            <a>重置密码</a>
          </Popconfirm>
        </Access>,
        <Access key="remove" accessible={access['system:user:remove']}>
          <Popconfirm
            key="remove"
            title="确认删除该用户？"
            onConfirm={() => handleDelete(record)}
          >
            <a style={{ color: '#ff4d4f' }}>删除</a>
          </Popconfirm>
        </Access>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<SysUser>
        headerTitle="用户列表"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        search={{ labelWidth: 'auto' }}
        request={async (params) => {
          const res = await queryUserList({
            pageNum: params.current,
            pageSize: params.pageSize,
            username: params.username,
            status: params.status,
            deptId: params.deptId,
          });
          return { data: res.list || [], total: res.total || 0, success: true };
        }}
        form={{ syncToUrl: false }}
        toolBarRender={() => [
          <Access key="add" accessible={access['system:user:add']}>
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

      <ModalForm<SysUser>
        key={editing?.id ?? 'create'}
        title={editing ? '编辑用户' : '新增用户'}
        open={open}
        onOpenChange={setOpen}
        width={520}
        initialValues={
          editing || { status: 1, sex: 0, deptId: 1, roleIds: [] }
        }
        modalProps={{ destroyOnClose: true }}
        onFinish={async (values) => {
          if (editing) {
            await updateUser({ ...values, id: editing.id });
            message.success('修改成功');
          } else {
            await createUser(values);
            message.success('新增成功');
          }
          reload();
          return true;
        }}
      >
        <ProFormText
          name="username"
          label="用户名"
          rules={[{ required: true, message: '请输入用户名' }]}
        />
        <ProFormText
          name="nickname"
          label="昵称"
          rules={[{ required: true, message: '请输入昵称' }]}
        />
        {!editing && (
          <ProFormText
            name="password"
            label="密码"
            placeholder="留空则使用默认密码 123456"
          />
        )}
        <ProFormTreeSelect
          name="deptId"
          label="部门"
          fieldProps={{ treeData: deptTree, treeDefaultExpandAll: true }}
        />
        <ProFormSelect
          name="roleIds"
          label="角色"
          mode="multiple"
          options={roles.map((r) => ({ label: r.name, value: r.id }))}
        />
        <ProFormText name="phone" label="手机号" />
        <ProFormText name="email" label="邮箱" />
        <ProFormRadio.Group
          name="sex"
          label="性别"
          options={[
            { label: '男', value: 1 },
            { label: '女', value: 2 },
            { label: '未知', value: 0 },
          ]}
        />
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
