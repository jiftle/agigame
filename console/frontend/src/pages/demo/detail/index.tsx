import { PageContainer } from '@ant-design/pro-components';
import { Card, Col, Descriptions, Row, Steps, Tag, Timeline } from 'antd';

const stepsItems = [
  { title: '创建项目', description: '2026-09-04 10:20' },
  { title: '提交审核', description: '2026-09-05 14:30' },
  { title: '审核通过', description: '进行中' },
  { title: '上线发布', description: '待开始' },
];

const timelineItems = [
  { color: 'blue', children: '张三 创建了项目' },
  { color: 'green', children: '李四 提交了审核' },
  { color: 'gray', children: '王五 更新了参数配置' },
];

export default function DemoDetail() {
  return (
    <PageContainer title="详情页">
      <Row gutter={[16, 16]}>
        <Col xs={24} lg={16}>
          <Card title="项目信息" variant="borderless">
            <Descriptions column={{ xs: 1, sm: 2 }} bordered size="small">
              <Descriptions.Item label="项目名称">AdminBase</Descriptions.Item>
              <Descriptions.Item label="负责人">张三</Descriptions.Item>
              <Descriptions.Item label="状态">
                <Tag color="processing">进行中</Tag>
              </Descriptions.Item>
              <Descriptions.Item label="预算">￥100,000</Descriptions.Item>
              <Descriptions.Item label="启动日期">2026-09-04</Descriptions.Item>
              <Descriptions.Item label="优先级">中</Descriptions.Item>
              <Descriptions.Item label="描述" span={2}>
                基于 GoFrame + Ant Design Pro 的通用后台管理基座演示项目。
              </Descriptions.Item>
            </Descriptions>
          </Card>
        </Col>
        <Col xs={24} lg={8}>
          <Card title="流程进度" variant="borderless">
            <Steps direction="vertical" current={2} items={stepsItems} />
          </Card>
          <Card title="操作日志" variant="borderless" style={{ marginTop: 16 }}>
            <Timeline items={timelineItems} />
          </Card>
        </Col>
      </Row>
    </PageContainer>
  );
}
