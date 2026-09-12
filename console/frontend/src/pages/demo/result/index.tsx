import { PageContainer, ProCard } from '@ant-design/pro-components';
import { Button, Result } from 'antd';

import { history } from '@/utils/history';

export default function DemoResult() {
  return (
    <PageContainer title="结果页">
      <ProCard gutter={[16, 16]} wrap>
        <ProCard colSpan={{ xs: 24, md: 12 }}>
          <Result
            status="success"
            title="提交成功"
            subTitle="项目已创建，正在等待审核。"
            extra={<Button type="primary">返回列表</Button>}
          />
        </ProCard>
        <ProCard colSpan={{ xs: 24, md: 12 }}>
          <Result
            status="error"
            title="提交失败"
            subTitle="请检查表单内容后重试。"
            extra={<Button>重试</Button>}
          />
        </ProCard>
        <ProCard colSpan={{ xs: 24, md: 12 }}>
          <Result
            status="warning"
            title="存在风险"
            subTitle="该操作可能影响线上数据，请谨慎执行。"
            extra={<Button onClick={() => history.push('/demo/exception/500')}>查看 500 页</Button>}
          />
        </ProCard>
        <ProCard colSpan={{ xs: 24, md: 12 }}>
          <Result
            status="403"
            title="403"
            subTitle="抱歉，你无权访问该页面。"
            extra={<Button onClick={() => history.push('/demo/exception/403')}>查看 403 页</Button>}
          />
        </ProCard>
      </ProCard>
    </PageContainer>
  );
}
