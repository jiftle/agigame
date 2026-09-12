import { PageContainer, ProCard } from '@ant-design/pro-components';
import { Alert, Tabs } from 'antd';

import DataDisplaySection from './sections/DataDisplaySection';
import DataEntrySection from './sections/DataEntrySection';
import FeedbackSection from './sections/FeedbackSection';
import GeneralSection from './sections/GeneralSection';
import NavigationSection from './sections/NavigationSection';
import ProSection from './sections/ProSection';

export default function DemoOverview() {
  return (
    <PageContainer
      title="组件总览"
      subTitle="验证 antd 6 / pro-components 在当前主题与明暗模式下的渲染，仅开发环境可见"
    >
      <Alert
        type="info"
        showIcon
        style={{ marginBottom: 16 }}
        title="验证提示"
        description="点击右上角头像 → 主题设置，切换明暗与主色后，逐块检查组件渲染与配色是否正常。"
      />
      <ProCard>
        <Tabs
          items={[
            { key: 'general', label: '通用', children: <GeneralSection /> },
            { key: 'data-entry', label: '数据录入', children: <DataEntrySection /> },
            { key: 'data-display', label: '数据展示', children: <DataDisplaySection /> },
            { key: 'feedback', label: '反馈', children: <FeedbackSection /> },
            { key: 'navigation', label: '导航', children: <NavigationSection /> },
            { key: 'pro', label: 'Pro 组件', children: <ProSection /> },
          ]}
        />
      </ProCard>
    </PageContainer>
  );
}
