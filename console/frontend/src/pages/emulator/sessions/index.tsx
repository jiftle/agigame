import type { ActionType, ProColumns, ProFormInstance } from '@ant-design/pro-components';
import { ModalForm, PageContainer, ProFormSelect, ProTable } from '@ant-design/pro-components';
import { Access, useAccess } from '@/access';
import { Button, Popconfirm, Tag } from 'antd';
import { useEffect, useRef, useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { message } from '@/utils/antdApp';

import type { EmuRomInfo, EmuSessionInfo, EmuStartInput } from '@/services/emulator';
import { queryRomList, querySessionList, startSession, stopSession } from '@/services/emulator';

export default function EmulatorSessionsPage() {
  const actionRef = useRef<ActionType | undefined>(undefined);
  const formRef = useRef<ProFormInstance<EmuStartInput> | undefined>(undefined);
  const access = useAccess();
  const navigate = useNavigate();
  const [roms, setRoms] = useState<EmuRomInfo[]>([]);

  useEffect(() => {
    queryRomList()
      .then(setRoms)
      .catch(() => undefined);
  }, []);

  const columns: ProColumns<EmuSessionInfo>[] = [
    { title: '会话ID', dataIndex: 'id', width: 120, search: false },
    { title: '平台', dataIndex: 'console', width: 70, search: false },
    { title: '卡带', dataIndex: 'cart', search: false },
    { title: '游戏', dataIndex: 'game', width: 80, search: false },
    { title: '模式', dataIndex: 'mode', width: 100, search: false },
    {
      title: '自动',
      dataIndex: 'auto',
      width: 80,
      search: false,
      render: (_, r) => (r.auto ? <Tag color="processing">Auto</Tag> : <Tag>手动</Tag>),
    },
    {
      title: '暂停',
      dataIndex: 'paused',
      width: 80,
      search: false,
      render: (_, r) => (r.paused ? <Tag color="warning">已暂停</Tag> : <Tag color="success">运行</Tag>),
    },
    { title: '帧数', dataIndex: 'frames', width: 100, search: false },
    { title: '创建时间', dataIndex: 'createdAt', valueType: 'dateTime', search: false },
    {
      title: '操作',
      valueType: 'option',
      width: 140,
      render: (_, record) => [
        <a key="open" onClick={() => navigate(`/emulator/sessions/${record.id}`)}>
          进入
        </a>,
        <Access key="stop" accessible={access['emu:session:control']}>
          <Popconfirm
            title="确认停止该会话？"
            onConfirm={async () => {
              await stopSession(record.id);
              message.success('已停止');
              actionRef.current?.reload();
            }}
          >
            <a style={{ color: '#ff4d4f' }}>停止</a>
          </Popconfirm>
        </Access>,
      ],
    },
  ];

  return (
    <PageContainer>
      <ProTable<EmuSessionInfo>
        headerTitle="模拟器会话"
        rowKey="id"
        actionRef={actionRef}
        columns={columns}
        search={false}
        options={false}
        request={async () => {
          const list = await querySessionList();
          return { data: list, total: list.length, success: true };
        }}
        toolBarRender={() => [
          <Access key="start" accessible={access['emu:session:control']}>
            <ModalForm<EmuStartInput>
              key="start"
              title="启动模拟器会话"
              trigger={<Button type="primary">启动会话</Button>}
              modalProps={{ destroyOnClose: true }}
              formRef={formRef}
              initialValues={{ console: 'gb', game: 'sml', mode: 'rules', palette: 'greyscale' }}
              onValuesChange={(changed) => {
                if (changed.rom) {
                  const rom = roms.find((r) => r.name === changed.rom);
                  if (rom) {
                    formRef.current?.setFieldsValue({ console: rom.console });
                  }
                }
              }}
              onFinish={async (values) => {
                const session = await startSession(values);
                message.success(`会话已启动：${session.id}`);
                actionRef.current?.reload();
                navigate(`/emulator/sessions/${session.id}`);
                return true;
              }}
            >
              <ProFormSelect
                name="rom"
                label="选择游戏"
                tooltip={`来自 ${'emulator.romDir'}，仅列出 .gb/.gbc/.gba`}
                placeholder={roms.length ? '请选择游戏' : 'roms 目录下暂无 ROM'}
                options={roms.map((r) => ({
                  label: `${r.title || r.name}（${r.console.toUpperCase()} · ${r.name}）`,
                  value: r.name,
                }))}
              />
              <ProFormSelect
                name="console"
                label="平台"
                tooltip="选择游戏后自动匹配；也可手动指定"
                options={[
                  { label: 'Game Boy / GBC', value: 'gb' },
                  { label: 'Game Boy Advance', value: 'gba' },
                ]}
              />
              <ProFormSelect
                name="game"
                label="游戏插件"
                tooltip="仅 GB 会话的 Agent 使用"
                options={[{ label: 'Super Mario Land', value: 'sml' }]}
              />
              <ProFormSelect
                name="mode"
                label="Agent 决策层"
                options={[
                  { label: '规则', value: 'rules' },
                  { label: '混合', value: 'hybrid' },
                  { label: 'LLM', value: 'llm' },
                  { label: '手动', value: 'manual' },
                ]}
              />
              <ProFormSelect
                name="palette"
                label="调色板"
                options={[
                  { label: '灰阶', value: 'greyscale' },
                  { label: '经典绿', value: 'original' },
                  { label: 'BGB 绿', value: 'bgb' },
                ]}
              />
            </ModalForm>
          </Access>,
        ]}
      />
    </PageContainer>
  );
}
