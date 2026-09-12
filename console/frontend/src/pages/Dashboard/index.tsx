import { Area, Pie } from '@ant-design/plots';
import { PageContainer, ProCard, StatisticCard } from '@ant-design/pro-components';
import { Col, Row, theme } from 'antd';

import ChartCard from '@/components/ChartCard';
import { useChartTheme } from '@/hooks/useChartTheme';
import { useCurrentUser } from '@/hooks/useCurrentUser';

const { Statistic } = StatisticCard;

const loginTrend = [
  { date: '09-04', count: 32 },
  { date: '09-05', count: 45 },
  { date: '09-06', count: 38 },
  { date: '09-07', count: 61 },
  { date: '09-08', count: 55 },
  { date: '09-09', count: 72 },
  { date: '09-10', count: 68 },
];

const roleDistribution = [
  { type: '超级管理员', value: 1 },
  { type: '运营', value: 3 },
  { type: '研发', value: 5 },
  { type: '访客', value: 2 },
];

export default function Dashboard() {
  const { data } = useCurrentUser();
  const user = data?.user;
  const roles = data?.roles ?? [];
  const perms = data?.perms ?? [];
  const menus = data?.menus ?? [];
  const isSuper = perms.includes('*:*:*');
  const { g2Theme, colorPrimary } = useChartTheme();
  const { token } = theme.useToken();

  return (
    <PageContainer title="仪表盘">
      <ProCard gutter={16} wrap>
        <StatisticCard
          colSpan={{ xs: 24, sm: 12, md: 6 }}
          statistic={{ title: '当前用户', value: user?.nickname || '-', suffix: '' }}
        />
        <StatisticCard
          colSpan={{ xs: 24, sm: 12, md: 6 }}
          statistic={{ title: '角色数', value: roles.length }}
        />
        <StatisticCard
          colSpan={{ xs: 24, sm: 12, md: 6 }}
          statistic={{
            title: '权限数',
            value: isSuper ? '全部' : perms.length,
          }}
        />
        <StatisticCard
          colSpan={{ xs: 24, sm: 12, md: 6 }}
          statistic={{ title: '菜单数', value: menus.length }}
        />
      </ProCard>

      <Row gutter={[16, 16]} style={{ marginTop: 16 }}>
        <Col xs={24} lg={16}>
          <ChartCard
            title="近 7 天登录趋势"
            extra={<span style={{ color: token.colorTextTertiary }}>示例数据</span>}
          >
            <Area
              data={loginTrend}
              xField="date"
              yField="count"
              shapeField="smooth"
              height={300}
              theme={g2Theme}
              style={{ fill: colorPrimary, fillOpacity: 0.2 }}
              axis={{
                x: { title: false },
                y: { title: false },
              }}
            />
          </ChartCard>
        </Col>
        <Col xs={24} lg={8}>
          <ChartCard
            title="用户角色分布"
            extra={<span style={{ color: token.colorTextTertiary }}>示例数据</span>}
          >
            <Pie
              data={roleDistribution}
              angleField="value"
              colorField="type"
              height={300}
              theme={g2Theme}
              radius={0.8}
              innerRadius={0.5}
              legend={false}
              label={{ text: 'type', position: 'outside' }}
            />
          </ChartCard>
        </Col>
      </Row>

      <ProCard title="欢迎使用 AdminBase" style={{ marginTop: 16 }}>
        <Statistic title="技术栈" value="Ant Design Pro + GoFrame + SQLite" />
      </ProCard>
    </PageContainer>
  );
}
