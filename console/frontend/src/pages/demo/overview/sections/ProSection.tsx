import type { ProColumns } from '@ant-design/pro-components';
import {
  ProCard,
  ProDescriptions,
  ProForm,
  ProFormSelect,
  ProFormText,
  ProList,
  ProTable,
  StatisticCard,
} from '@ant-design/pro-components';
import { Tag } from 'antd';

import { message } from '@/utils/antdApp';

import DemoBlock from '../DemoBlock';

interface Item {
  id: number;
  name: string;
  status: number;
}

const data: Item[] = [
  { id: 1, name: '项目一', status: 1 },
  { id: 2, name: '项目二', status: 0 },
];

const columns: ProColumns<Item>[] = [
  { title: 'ID', dataIndex: 'id', width: 60, search: false },
  { title: '名称', dataIndex: 'name' },
  {
    title: '状态',
    dataIndex: 'status',
    valueType: 'select',
    valueEnum: {
      1: { text: '正常', status: 'Success' },
      0: { text: '停用', status: 'Error' },
    },
  },
];

export default function ProSection() {
  return (
    <>
      <DemoBlock title="ProCard / StatisticCard">
        <ProCard gutter={16} wrap>
          <ProCard colSpan={{ xs: 24, md: 12 }}>普通 ProCard 内容</ProCard>
          <ProCard colSpan={{ xs: 24, md: 12 }}>普通 ProCard 内容</ProCard>
        </ProCard>
        <ProCard gutter={16} wrap>
          <StatisticCard
            colSpan={{ xs: 24, sm: 12, md: 8 }}
            statistic={{ title: '用户总数', value: 1280, suffix: '人' }}
          />
          <StatisticCard
            colSpan={{ xs: 24, sm: 12, md: 8 }}
            statistic={{ title: '角色数', value: 6 }}
          />
          <StatisticCard
            colSpan={{ xs: 24, sm: 12, md: 8 }}
            statistic={{ title: '权限数', value: '全部' }}
          />
        </ProCard>
      </DemoBlock>

      <DemoBlock title="ProDescriptions">
        <ProDescriptions
          column={2}
          dataSource={{ name: 'AdminBase', owner: '张三', status: '进行中' }}
          columns={[
            { title: '项目', dataIndex: 'name' },
            { title: '负责人', dataIndex: 'owner' },
            {
              title: '状态',
              dataIndex: 'status',
              render: (_, record) => <Tag color="processing">{record.status}</Tag>,
            },
          ]}
        />
      </DemoBlock>

      <DemoBlock title="ProList">
        <ProList<{ title: string; subTitle: string }>
          rowKey="title"
          dataSource={[
            { title: '规则一', subTitle: '遵循项目惯例' },
            { title: '规则二', subTitle: '最小改动' },
            { title: '规则三', subTitle: '类型安全' },
          ]}
          metas={{
            title: { dataIndex: 'title' },
            subTitle: { dataIndex: 'subTitle' },
          }}
        />
      </DemoBlock>

      <DemoBlock title="ProTable">
        <ProTable<Item>
          rowKey="id"
          headerTitle="示例列表"
          size="small"
          columns={columns}
          search={false}
          options={false}
          pagination={false}
          request={async () => ({ data, total: data.length, success: true })}
        />
      </DemoBlock>

      <DemoBlock title="ProForm">
        <ProForm
          style={{ maxWidth: 480 }}
          onFinish={async () => {
            message.success('提交成功');
            return true;
          }}
        >
          <ProFormText
            name="name"
            label="名称"
            rules={[{ required: true, message: '请输入名称' }]}
          />
          <ProFormSelect
            name="type"
            label="类型"
            options={[
              { label: '类型 A', value: 'a' },
              { label: '类型 B', value: 'b' },
            ]}
          />
        </ProForm>
      </DemoBlock>
    </>
  );
}
