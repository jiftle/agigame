import { QrcodeOutlined } from '@ant-design/icons';
import {
  Button,
  Carousel,
  Collapse,
  Descriptions,
  Image,
  Listy,
  Popover,
  QRCode,
  Space,
  Statistic,
  Table,
  Tag,
  Timeline,
  Tooltip,
  Tree,
  Watermark,
  theme,
} from 'antd';
import type { TableProps } from 'antd';

import DemoBlock from '../DemoBlock';

interface Row {
  key: number;
  name: string;
  age: number;
  status: string;
}

const columns: TableProps<Row>['columns'] = [
  { title: '姓名', dataIndex: 'name' },
  { title: '年龄', dataIndex: 'age' },
  {
    title: '状态',
    dataIndex: 'status',
    render: (value: string) => <Tag color="processing">{value}</Tag>,
  },
];

const dataSource: Row[] = [
  { key: 1, name: '张三', age: 28, status: '在职' },
  { key: 2, name: '李四', age: 32, status: '在职' },
  { key: 3, name: '王五', age: 24, status: '离职' },
];

const treeData = [
  {
    title: '系统管理',
    key: 'system',
    children: [
      { title: '用户管理', key: 'user' },
      { title: '角色管理', key: 'role' },
    ],
  },
];

export default function DataDisplaySection() {
  const { token } = theme.useToken();

  return (
    <>
      <DemoBlock title="Table 表格 / Tree 树">
        <Table<Row>
          rowKey="key"
          columns={columns}
          dataSource={dataSource}
          pagination={false}
          size="small"
        />
        <Tree treeData={treeData} defaultExpandAll selectable={false} />
      </DemoBlock>

      <DemoBlock title="Descriptions / Statistic">
        <Descriptions column={{ xs: 1, sm: 2 }} bordered size="small">
          <Descriptions.Item label="项目名称">AdminBase</Descriptions.Item>
          <Descriptions.Item label="负责人">张三</Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag color="processing">进行中</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="预算">￥100,000</Descriptions.Item>
        </Descriptions>
        <Space size="large" wrap>
          <Statistic title="活跃用户" value={1128} />
          <Statistic title="转化率" value={11.28} suffix="%" />
          <Statistic title="金额" value={9280} prefix="￥" />
          <Statistic
            title="同比增长"
            value={12.3}
            suffix="%"
            styles={{ content: { color: token.colorSuccess } }}
          />
        </Space>
      </DemoBlock>

      <DemoBlock title="Listy / Timeline / Collapse">
        <Listy<{ id: number; text: string }>
          rowKey="id"
          height={150}
          items={[
            { id: 1, text: '规则一：遵循项目惯例' },
            { id: 2, text: '规则二：最小改动' },
            { id: 3, text: '规则三：类型安全' },
          ]}
          itemRender={(item) => item.text}
        />
        <Timeline
          items={[
            { color: 'green', children: '创建成功' },
            { color: 'blue', children: '提交审核' },
            { color: 'gray', children: '等待发布' },
          ]}
        />
        <Collapse
          items={[
            {
              key: '1',
              label: '折叠面板标题',
              children: <p style={{ margin: 0 }}>折叠面板内容，用于分组展示信息。</p>,
            },
          ]}
        />
      </DemoBlock>

      <DemoBlock title="Carousel / Image">
        <Carousel autoplay style={{ maxWidth: 480 }}>
          {['一', '二', '三'].map((label) => (
            <div key={label}>
              <div
                style={{
                  height: 120,
                  lineHeight: '120px',
                  textAlign: 'center',
                  color: token.colorPrimary,
                  background: token.colorPrimaryBg,
                  borderRadius: token.borderRadius,
                }}
              >
                轮播 {label}
              </div>
            </div>
          ))}
        </Carousel>
        <Image src="/logo.svg" width={96} alt="AdminBase" />
      </DemoBlock>

      <DemoBlock title="Tooltip / Popover / QRCode / Watermark">
        <Space wrap>
          <Tooltip title="这是提示">
            <Button>Tooltip</Button>
          </Tooltip>
          <Popover title="标题" content="这是气泡卡片内容">
            <Button>Popover</Button>
          </Popover>
          <Popover
            content={
              <QRCode value="https://ant.design" size={120} bordered={false} />
            }
          >
            <Button icon={<QrcodeOutlined />}>二维码</Button>
          </Popover>
        </Space>
        <Watermark content="AdminBase">
          <div
            style={{
              height: 120,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
              border: `1px dashed ${token.colorBorder}`,
              borderRadius: token.borderRadius,
            }}
          >
            水印区域
          </div>
        </Watermark>
      </DemoBlock>
    </>
  );
}
