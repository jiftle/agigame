import {
  CheckCircleOutlined,
  DownloadOutlined,
  PlusOutlined,
  SyncOutlined,
  UserOutlined,
} from '@ant-design/icons';
import { Avatar, Badge, Button, Divider, Flex, Space, Tag, Typography, theme } from 'antd';

import DemoBlock from '../DemoBlock';

const { Title, Text, Paragraph, Link } = Typography;

export default function GeneralSection() {
  const { token } = theme.useToken();

  return (
    <>
      <DemoBlock title="Button 按钮">
        <Space wrap>
          <Button type="primary">Primary</Button>
          <Button>Default</Button>
          <Button type="dashed">Dashed</Button>
          <Button type="text">Text</Button>
          <Button type="link">Link</Button>
          <Button type="primary" danger>
            Danger
          </Button>
          <Button disabled>Disabled</Button>
          <Button type="primary" loading>
            Loading
          </Button>
          <Button type="primary" icon={<PlusOutlined />}>
            图标
          </Button>
          <Button type="primary" shape="circle" icon={<DownloadOutlined />} />
          <Button icon={<SyncOutlined spin />} />
          <Button type="primary" icon={<CheckCircleOutlined />} />
        </Space>
        <Space wrap>
          <Button size="large" type="primary">
            Large
          </Button>
          <Button type="primary">Middle</Button>
          <Button size="small" type="primary">
            Small
          </Button>
        </Space>
      </DemoBlock>

      <DemoBlock title="Typography 排版">
        <Title level={3} style={{ margin: 0 }}>
          标题 Title
        </Title>
        <Paragraph style={{ margin: 0 }}>
          段落文本，
          <Text strong>加粗</Text>、<Text type="secondary">次要</Text>、
          <Text type="danger">危险</Text>、<Text mark>标记</Text>、
          <Text code>code</Text>、<Text keyboard>K</Text>。
          <Link href="https://ant.design" target="_blank">
            链接
          </Link>
        </Paragraph>
      </DemoBlock>

      <DemoBlock title="Tag / Badge / Avatar">
        <Space wrap>
          <Tag>默认</Tag>
          <Tag color="processing">进行中</Tag>
          <Tag color="success">成功</Tag>
          <Tag color="warning">警告</Tag>
          <Tag color="error">错误</Tag>
          <Tag color="magenta">自定义</Tag>
          <Tag variant="filled">无边框</Tag>
        </Space>
        <Space size="large" wrap>
          <Badge count={5}>
            <Avatar shape="square" icon={<UserOutlined />} />
          </Badge>
          <Badge count={0} showZero>
            <Avatar shape="square" icon={<UserOutlined />} />
          </Badge>
          <Badge dot>
            <Avatar shape="square" icon={<UserOutlined />} />
          </Badge>
          <Badge status="success" text="成功" />
          <Badge status="processing" text="处理中" />
          <Badge status="warning" text="警告" />
          <Badge status="error" text="失败" />
          <Badge status="default" text="默认" />
          <Avatar icon={<UserOutlined />} />
          <Avatar style={{ backgroundColor: token.colorPrimary }}>AB</Avatar>
        </Space>
      </DemoBlock>

      <DemoBlock title="Divider / Flex / Space">
        <Flex gap="middle" wrap>
          <Button>Flex A</Button>
          <Button>Flex B</Button>
          <Button>Flex C</Button>
        </Flex>
        <Divider plain>分割线</Divider>
        <Space separator={<Divider orientation="vertical" />}>
          <Text>一</Text>
          <Text>二</Text>
          <Text>三</Text>
        </Space>
      </DemoBlock>
    </>
  );
}
