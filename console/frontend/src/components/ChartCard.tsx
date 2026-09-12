import { Card, Empty, Skeleton } from 'antd';
import type { ReactNode } from 'react';

export interface ChartCardProps {
  title: ReactNode;
  extra?: ReactNode;
  loading?: boolean;
  empty?: boolean;
  height?: number;
  children: ReactNode;
}

/**
 * 图表卡片：统一标题、loading、空态与高度
 */
export default function ChartCard({
  title,
  extra,
  loading,
  empty,
  height = 300,
  children,
}: ChartCardProps) {
  const showPlaceholder = loading || empty;
  return (
    <Card title={title} extra={extra} style={{ height: '100%' }}>
      {showPlaceholder ? (
        <div style={{ minHeight: height, display: 'flex', alignItems: 'center' }}>
          <div style={{ width: '100%' }}>
            {loading ? (
              <Skeleton active paragraph={{ rows: 6 }} />
            ) : (
              <Empty
                image={Empty.PRESENTED_IMAGE_SIMPLE}
                description="暂无数据"
              />
            )}
          </div>
        </div>
      ) : (
        children
      )}
    </Card>
  );
}
