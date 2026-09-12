import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import {
  ModalForm,
  PageContainer,
  ProFormDigit,
  ProFormRadio,
  ProFormText,
  ProFormTreeSelect,
  ProTable,
} from '@ant-design/pro-components';
import { Access, useAccess } from '@/access';
import { Button, Popconfirm, Tag } from 'antd';
import { useRef, useState } from 'react';

import { message } from '@/utils/antdApp';

import type { SysMenu } from '@/services/system';
import {
  createMenu,
  deleteMenu,
  getMenu,
  queryMenuList,
  updateMenu,
} from '@/services/system';

const toTreeData = (list: SysMenu[]): any[] =>
  (list || []).map((m) => ({
    title: m.title,
    value: m.id,
    children: m.children ? toTreeData(m.children) : undefined,
  }));

const typeEnum: Record<string, { text: string; color: string }> = {
  M: { text: '目录', color: 'blue' },
  C: { text: '菜单', color: 'green' },
  F: { text: '按钮', color: 'orange' },
};

export default function MenuPage() {
  const actionRef = useRef<ActionType | undefined>(undefined);
  const access = useAccess();
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<SysMenu | undefined>();
  const [tree, setTree] = useState<SysMenu[]>([]);

  const reload = () => actionRef.current?.reload();

  const loadTree = async () => {
    const list = await queryMenuList();
    setTree(list);
    return list;
  };

  const handleEdit = async (record: SysMenu) => {
    const menu = await getMenu(record.id as number);
    setEditing(menu);
    setOpen(true);
  };

  const columns: ProColumns<SysMenu>[] = [
    { title: '菜单名称', dataIndex: 'title' },
    { title: '图标', dataIndex: 'icon', search: false },
    {
      title: '类型',
      dataIndex: 'type',
      search: false,
      render: (_, record) => {
        const t = typeEnum[record.type as string];
        return t ? <Tag color={t.color}>{t.text}</Tag> : record.type;
      },
    },
    { title: '路由地址', dataIndex: 'path', search: false },
    { title: '权限标识', dataIndex: 'perms', search: false },
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
    {
      title: '操作',
      valueType: 'option',
      width: 140,
      render: (_, record) => [
        <Access key="edit" accessible={access['system:menu:edit']}>
          <a onClick={() => handleEdit(record)}>编辑</a>
        </Access>,
        <Access key="remove" accessible={access['system:menu:remove']}>
          <Popconfirm
            key="remove"
            title="确认删除该菜单？"
            onConfirm={async () => {
              await deleteMenu(record.id as number);
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
      <ProTable<SysMenu>
        headerTitle="菜单列表"
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
          <Access key="add" accessible={access['system:menu:add']}>
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

      <ModalForm<SysMenu>
        key={editing?.id ?? 'create'}
        title={editing ? '编辑菜单' : '新增菜单'}
        open={open}
        onOpenChange={setOpen}
        width={560}
        initialValues={
          editing || {
            parentId: 0,
            type: 'C',
            sort: 0,
            visible: 1,
            status: 1,
            isFrame: 0,
            isCache: 1,
          }
        }
        modalProps={{ destroyOnClose: true }}
        onFinish={async (values) => {
          if (editing) {
            await updateMenu({ ...values, id: editing.id });
            message.success('修改成功');
          } else {
            await createMenu(values);
            message.success('新增成功');
          }
          reload();
          return true;
        }}
      >
        <ProFormTreeSelect
          name="parentId"
          label="上级菜单"
          fieldProps={{
            treeData: [{ title: '根目录', value: 0, children: toTreeData(tree) }],
            treeDefaultExpandAll: true,
          }}
        />
        <ProFormRadio.Group
          name="type"
          label="菜单类型"
          options={[
            { label: '目录', value: 'M' },
            { label: '菜单', value: 'C' },
            { label: '按钮', value: 'F' },
          ]}
        />
        <ProFormText
          name="title"
          label="菜单标题"
          rules={[{ required: true, message: '请输入菜单标题' }]}
        />
        <ProFormText name="name" label="路由名称" />
        <ProFormText name="path" label="路由地址" />
        <ProFormText name="component" label="组件路径" />
        <ProFormText name="icon" label="图标" />
        <ProFormText name="perms" label="权限标识" />
        <ProFormDigit name="sort" label="排序" fieldProps={{ precision: 0 }} />
        <ProFormRadio.Group
          name="status"
          label="状态"
          options={[
            { label: '正常', value: 1 },
            { label: '停用', value: 0 },
          ]}
        />
        <ProFormRadio.Group
          name="visible"
          label="显示"
          options={[
            { label: '显示', value: 1 },
            { label: '隐藏', value: 0 },
          ]}
        />
      </ModalForm>
    </PageContainer>
  );
}
