import { DownOutlined } from '@ant-design/icons';
import {
  Breadcrumb,
  Button,
  Dropdown,
  Pagination,
  Segmented,
  Space,
  Steps,
  Tabs,
  theme,
} from 'antd';

import { message } from '@/utils/antdApp';

import DemoBlock from '../DemoBlock';

export default function NavigationSection() {
  const { token } = theme.useToken();

  return (
    <>
      <DemoBlock title="Tabs 标签页 / Segmented 分段控制器">
        <Tabs
          defaultActiveKey="1"
          items={[
            { key: '1', label: '标签一', children: '标签一的内容' },
            { key: '2', label: '标签二', children: '标签二的内容' },
            { key: '3', label: '标签三', children: '标签三的内容', disabled: true },
          ]}
        />
        <Tabs
          type="card"
          defaultActiveKey="1"
          items={[
            { key: '1', label: '卡片一', children: '卡片式标签页' },
            { key: '2', label: '卡片二', children: '卡片式标签页' },
          ]}
        />
        <Segmented options={['日报', '周报', '月报']} />
      </DemoBlock>

      <DemoBlock title="Steps 步骤条">
        <Steps
          current={1}
          items={[
            { title: '已完成', description: '创建项目' },
            { title: '进行中', description: '提交审核' },
            { title: '待处理', description: '发布上线' },
          ]}
        />
        <Steps
          orientation="vertical"
          size="small"
          current={1}
          items={[
            { title: '第一步' },
            { title: '第二步' },
            { title: '第三步' },
          ]}
        />
      </DemoBlock>

      <DemoBlock title="Breadcrumb 面包屑 / Dropdown 下拉">
        <Breadcrumb
          items={[
            { title: '首页' },
            { title: '系统管理' },
            { title: '用户管理' },
          ]}
        />
        <Dropdown
          menu={{
            items: [
              { key: 'edit', label: '编辑' },
              { key: 'copy', label: '复制' },
              { type: 'divider' },
              { key: 'remove', label: '删除', danger: true },
            ],
            onClick: ({ key }) => message.info(`点击：${key}`),
          }}
        >
          <Button>
            下拉菜单 <DownOutlined />
          </Button>
        </Dropdown>
      </DemoBlock>

      <DemoBlock title="Pagination 分页">
        <Space orientation="vertical">
          <Pagination defaultCurrent={1} total={85} showSizeChanger showQuickJumper />
          <Pagination simple defaultCurrent={1} total={50} />
          <span style={{ color: token.colorTextTertiary }}>
            共 85 条 · 当前每页 10 条
          </span>
        </Space>
      </DemoBlock>
    </>
  );
}
