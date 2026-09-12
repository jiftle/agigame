import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import {
  ModalForm,
  PageContainer,
  ProFormDigit,
  ProFormRadio,
  ProFormSelect,
  ProFormText,
  ProFormTextArea,
  ProTable,
} from '@ant-design/pro-components';
import { Access, useAccess } from '@/access';
import { Button, Popconfirm, Tabs, Tag } from 'antd';
import { useRef, useState } from 'react';

import { message } from '@/utils/antdApp';

import type { DictData, DictType } from '@/services/system';
import {
  createDictData,
  createDictType,
  deleteDictData,
  deleteDictType,
  getDictData,
  getDictType,
  queryDictDataList,
  queryDictTypeList,
  updateDictData,
  updateDictType,
} from '@/services/system';

function DictTypeTab() {
  const actionRef = useRef<ActionType | undefined>(undefined);
  const access = useAccess();
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<DictType | undefined>();

  const reload = () => actionRef.current?.reload();

  const columns: ProColumns<DictType>[] = [
    { title: 'ID', dataIndex: 'id', width: 70, search: false },
    { title: '字典名称', dataIndex: 'name' },
    { title: '字典类型', dataIndex: 'type' },
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
        <Access key="edit" accessible={access['system:dict:edit']}>
          <a
            onClick={async () => {
              setEditing(await getDictType(record.id as number));
              setOpen(true);
            }}
          >
            编辑
          </a>
        </Access>,
        <Access key="remove" accessible={access['system:dict:remove']}>
          <Popconfirm
            key="remove"
            title="确认删除该字典类型？"
            onConfirm={async () => {
              await deleteDictType([record.id as number]);
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
    <>
      <ProTable<DictType>
        headerTitle="字典类型"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        search={{ labelWidth: 'auto' }}
        request={async (params) => {
          const res = await queryDictTypeList({
            pageNum: params.current,
            pageSize: params.pageSize,
            name: params.name,
            type: params.type,
            status: params.status,
          });
          return { data: res.list || [], total: res.total || 0, success: true };
        }}
        toolBarRender={() => [
          <Access key="add" accessible={access['system:dict:add']}>
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
      <ModalForm<DictType>
        key={editing?.id ?? 'create'}
        title={editing ? '编辑字典类型' : '新增字典类型'}
        open={open}
        onOpenChange={setOpen}
        width={480}
        initialValues={editing || { status: 1 }}
        modalProps={{ destroyOnClose: true }}
        onFinish={async (values) => {
          if (editing) {
            await updateDictType({ ...values, id: editing.id });
            message.success('修改成功');
          } else {
            await createDictType(values);
            message.success('新增成功');
          }
          reload();
          return true;
        }}
      >
        <ProFormText
          name="name"
          label="字典名称"
          rules={[{ required: true, message: '请输入字典名称' }]}
        />
        <ProFormText
          name="type"
          label="字典类型"
          rules={[{ required: true, message: '请输入字典类型' }]}
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
    </>
  );
}

function DictDataTab() {
  const actionRef = useRef<ActionType | undefined>(undefined);
  const access = useAccess();
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<DictData | undefined>();
  const [typeOptions, setTypeOptions] = useState<{ label: string; value: string }[]>(
    [],
  );

  const reload = () => actionRef.current?.reload();

  const loadTypes = async () => {
    const res = await queryDictTypeList({ pageNum: 1, pageSize: 1000 });
    setTypeOptions((res.list || []).map((t) => ({ label: `${t.name}(${t.type})`, value: t.type as string })));
  };

  const columns: ProColumns<DictData>[] = [
    { title: 'ID', dataIndex: 'id', width: 70, search: false },
    { title: '排序', dataIndex: 'dictSort', search: false },
    { title: '字典标签', dataIndex: 'dictLabel' },
    { title: '字典键值', dataIndex: 'dictValue' },
    { title: '字典类型', dataIndex: 'dictType', valueType: 'select' },
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
        <Access key="edit" accessible={access['system:dict:edit']}>
          <a
            onClick={async () => {
              await loadTypes();
              setEditing(await getDictData(record.id as number));
              setOpen(true);
            }}
          >
            编辑
          </a>
        </Access>,
        <Access key="remove" accessible={access['system:dict:remove']}>
          <Popconfirm
            key="remove"
            title="确认删除该字典数据？"
            onConfirm={async () => {
              await deleteDictData([record.id as number]);
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
    <>
      <ProTable<DictData>
        headerTitle="字典数据"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        search={{ labelWidth: 'auto' }}
        request={async (params) => {
          const res = await queryDictDataList({
            pageNum: params.current,
            pageSize: params.pageSize,
            dictType: params.dictType,
            dictLabel: params.dictLabel,
          });
          return { data: res.list || [], total: res.total || 0, success: true };
        }}
        toolBarRender={() => [
          <Access key="add" accessible={access['system:dict:add']}>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={async () => {
                await loadTypes();
                setEditing(undefined);
                setOpen(true);
              }}
            >
              新增
            </Button>
          </Access>,
        ]}
      />
      <ModalForm<DictData>
        key={editing?.id ?? 'create'}
        title={editing ? '编辑字典数据' : '新增字典数据'}
        open={open}
        onOpenChange={setOpen}
        width={480}
        initialValues={editing || { dictSort: 0, status: 1, isDefault: 0 }}
        modalProps={{ destroyOnClose: true }}
        onFinish={async (values) => {
          if (editing) {
            await updateDictData({ ...values, id: editing.id });
            message.success('修改成功');
          } else {
            await createDictData(values);
            message.success('新增成功');
          }
          reload();
          return true;
        }}
      >
        <ProFormSelect
          name="dictType"
          label="字典类型"
          options={typeOptions}
          rules={[{ required: true, message: '请选择字典类型' }]}
        />
        <ProFormText
          name="dictLabel"
          label="字典标签"
          rules={[{ required: true, message: '请输入字典标签' }]}
        />
        <ProFormText
          name="dictValue"
          label="字典键值"
          rules={[{ required: true, message: '请输入字典键值' }]}
        />
        <ProFormDigit name="dictSort" label="排序" fieldProps={{ precision: 0 }} />
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
    </>
  );
}

export default function DictPage() {
  return (
    <PageContainer>
      <Tabs
        items={[
          { key: 'type', label: '字典类型', children: <DictTypeTab /> },
          { key: 'data', label: '字典数据', children: <DictDataTab /> },
        ]}
      />
    </PageContainer>
  );
}
