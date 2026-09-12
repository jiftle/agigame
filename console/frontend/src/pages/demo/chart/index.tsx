import { Area, Column, Line, Pie } from '@ant-design/plots';
import { PageContainer } from '@ant-design/pro-components';
import { Col, Row } from 'antd';

import ChartCard from '@/components/ChartCard';
import { useChartTheme } from '@/hooks/useChartTheme';

const days = ['09-04', '09-05', '09-06', '09-07', '09-08', '09-09', '09-10'];
const trendValue = [32, 45, 38, 61, 55, 72, 68];
const trendData = days.map((date, i) => ({ date, value: trendValue[i] }));

const months = ['一月', '二月', '三月', '四月', '五月', '六月'];
const monthValue = [120, 200, 150, 80, 70, 110];
const columnData = months.map((month, i) => ({ month, value: monthValue[i] }));

const pieData = [
  { type: '超级管理员', value: 1 },
  { type: '运营', value: 3 },
  { type: '研发', value: 5 },
  { type: '访客', value: 2 },
];

export default function DemoChart() {
  const { g2Theme, colorPrimary } = useChartTheme();

  return (
    <PageContainer title="图表示例">
      <Row gutter={[16, 16]}>
        <Col xs={24} lg={12}>
          <ChartCard title="折线图">
            <Line
              data={trendData}
              xField="date"
              yField="value"
              height={300}
              theme={g2Theme}
              style={{ stroke: colorPrimary }}
            />
          </ChartCard>
        </Col>
        <Col xs={24} lg={12}>
          <ChartCard title="面积图">
            <Area
              data={trendData}
              xField="date"
              yField="value"
              shapeField="smooth"
              height={300}
              theme={g2Theme}
              style={{ fill: colorPrimary, fillOpacity: 0.2 }}
            />
          </ChartCard>
        </Col>
        <Col xs={24} lg={12}>
          <ChartCard title="柱状图">
            <Column
              data={columnData}
              xField="month"
              yField="value"
              height={300}
              theme={g2Theme}
              style={{ fill: colorPrimary }}
            />
          </ChartCard>
        </Col>
        <Col xs={24} lg={12}>
          <ChartCard title="饼图">
            <Pie
              data={pieData}
              angleField="value"
              colorField="type"
              height={300}
              theme={g2Theme}
              innerRadius={0.5}
              radius={0.8}
            />
          </ChartCard>
        </Col>
      </Row>
    </PageContainer>
  );
}
