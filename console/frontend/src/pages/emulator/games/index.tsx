import type { ActionType, ProColumns } from '@ant-design/pro-components';
import { PageContainer, ProTable } from '@ant-design/pro-components';
import { Access, useAccess } from '@/access';
import { Button, Tag } from 'antd';
import { useRef } from 'react';
import { useNavigate } from 'react-router-dom';

import { message } from '@/utils/antdApp';

import type { EmuRomInfo } from '@/services/emulator';
import { queryRomList, startSession } from '@/services/emulator';

function formatSize(bytes: number): string {
  if (!bytes) return '-';
  const mb = bytes / (1024 * 1024);
  if (mb >= 1) return `${mb.toFixed(1)} MB`;
  return `${Math.max(1, Math.round(bytes / 1024))} KB`;
}

export default function EmulatorGamesPage() {
  const actionRef = useRef<ActionType | undefined>(undefined);
  const access = useAccess();
  const navigate = useNavigate();

  const launch = async (rom: EmuRomInfo) => {
    try {
      const session = await startSession({
        rom: rom.name,
        console: rom.console,
        mode: 'rules',
        palette: 'greyscale',
      });
      message.success(`已启动：${rom.title || rom.name}`);
      navigate(`/emulator/sessions/${session.id}`);
    } catch {
      // 错误提示已由请求拦截器统一处理
    }
  };

  const columns: ProColumns<EmuRomInfo>[] = [
    {
      title: '游戏',
      dataIndex: 'title',
      render: (_, r) => r.title || r.name,
    },
    {
      title: '平台',
      dataIndex: 'console',
      width: 100,
      render: (_, r) => (
        <Tag color={r.console === 'gba' ? 'purple' : 'blue'}>{r.console.toUpperCase()}</Tag>
      ),
    },
    { title: '文件名', dataIndex: 'name', ellipsis: true },
    {
      title: '大小',
      dataIndex: 'size',
      width: 110,
      render: (_, r) => formatSize(r.size),
    },
    {
      title: '操作',
      valueType: 'option',
      width: 120,
      render: (_, record) => [
        <Access key="launch" accessible={access['emu:session:control']}>
          <a key="launch" onClick={() => launch(record)}>
            启动游戏
          </a>
        </Access>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<EmuRomInfo>
        headerTitle="游戏"
        rowKey="name"
        actionRef={actionRef}
        columns={columns}
        search={false}
        options={false}
        request={async () => {
          const list = await queryRomList();
          return { data: list, total: list.length, success: true };
        }}
        toolBarRender={() => [
          <Button key="refresh" onClick={() => actionRef.current?.reload()}>
            刷新
          </Button>,
        ]}
      />
    </PageContainer>
  );
}
