import { PlusOutlined } from '@ant-design/icons';
import type { ActionType, ProColumns } from '@ant-design/pro-components';
import {
  ModalForm,
  PageContainer,
  ProFormRadio,
  ProFormText,
  ProFormTextArea,
  ProTable,
} from '@ant-design/pro-components';
import { Access, useAccess } from '@/access';
import { Button, Popconfirm } from 'antd';
import { useRef, useState } from 'react';

import { message } from '@/utils/antdApp';

import type { SysConfig } from '@/services/system';
import {
  createConfig,
  deleteConfig,
  getConfig,
  queryConfigList,
  updateConfig,
} from '@/services/system';

export default function ConfigPage() {
  const actionRef = useRef<ActionType | undefined>(undefined);
  const access = useAccess();
  const [open, setOpen] = useState(false);
  const [editing, setEditing] = useState<SysConfig | undefined>();

  const reload = () => actionRef.current?.reload();

  const columns: ProColumns<SysConfig>[] = [
    { title: 'ID', dataIndex: 'id', width: 70, search: false },
    { title: '参数名称', dataIndex: 'configName' },
    { title: '参数键名', dataIndex: 'configKey' },
    { title: '参数键值', dataIndex: 'configValue', search: false },
    {
      title: '系统内置',
      dataIndex: 'configType',
      search: false,
      render: (_, record) => (record.configType === 0 ? '是' : '否'),
    },
    { title: '备注', dataIndex: 'remark', search: false, ellipsis: true },
    {
      title: '操作',
      valueType: 'option',
      width: 140,
      render: (_, record) => [
        <Access key="edit" accessible={access['system:config:edit']}>
          <a
            onClick={async () => {
              setEditing(await getConfig(record.id as number));
              setOpen(true);
            }}
          >
            编辑
          </a>
        </Access>,
        <Access key="remove" accessible={access['system:config:remove']}>
          <Popconfirm
            key="remove"
            title="确认删除该参数？"
            onConfirm={async () => {
              await deleteConfig([record.id as number]);
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
      <ProTable<SysConfig>
        headerTitle="参数列表"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        search={{ labelWidth: 'auto' }}
        request={async (params) => {
          const res = await queryConfigList({
            pageNum: params.current,
            pageSize: params.pageSize,
            configName: params.configName,
            configKey: params.configKey,
          });
          return { data: res.list || [], total: res.total || 0, success: true };
        }}
        toolBarRender={() => [
          <Access key="add" accessible={access['system:config:add']}>
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
      <ModalForm<SysConfig>
        key={editing?.id ?? 'create'}
        title={editing ? '编辑参数' : '新增参数'}
        open={open}
        onOpenChange={setOpen}
        width={480}
        initialValues={editing || { configType: 0 }}
        modalProps={{ destroyOnClose: true }}
        onFinish={async (values) => {
          if (editing) {
            await updateConfig({ ...values, id: editing.id });
            message.success('修改成功');
          } else {
            await createConfig(values);
            message.success('新增成功');
          }
          reload();
          return true;
        }}
      >
        <ProFormText
          name="configName"
          label="参数名称"
          rules={[{ required: true, message: '请输入参数名称' }]}
        />
        <ProFormText
          name="configKey"
          label="参数键名"
          rules={[{ required: true, message: '请输入参数键名' }]}
        />
        <ProFormText
          name="configValue"
          label="参数键值"
          rules={[{ required: true, message: '请输入参数键值' }]}
        />
        <ProFormRadio.Group
          name="configType"
          label="系统内置"
          options={[
            { label: '是', value: 0 },
            { label: '否', value: 1 },
          ]}
        />
        <ProFormTextArea name="remark" label="备注" />
      </ModalForm>
    </PageContainer>
  );
}
