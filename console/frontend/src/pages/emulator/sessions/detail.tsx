import { PageContainer } from '@ant-design/pro-components';
import { Badge, Button, Card, Descriptions, Select, Space, Switch, Tag, Typography } from 'antd';
import { useEffect, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';

import { message } from '@/utils/antdApp';

import { useEmulatorSocket } from '@/hooks/useEmulatorSocket';
import { getSession, stopSession, type EmuSessionInfo } from '@/services/emulator';

const KEYMAP: Record<string, string> = {
  ArrowUp: 'Up',
  ArrowDown: 'Down',
  ArrowLeft: 'Left',
  ArrowRight: 'Right',
  z: 'A',
  Z: 'A',
  x: 'B',
  X: 'B',
  Enter: 'Start',
  Backspace: 'Select',
  q: 'L',
  Q: 'L',
  w: 'R',
  W: 'R',
};

export default function EmulatorSessionDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const canvasRef = useRef<HTMLCanvasElement | null>(null);
  const autoRef = useRef(false);
  const heldRef = useRef<Set<string>>(new Set());

  const [info, setInfo] = useState<EmuSessionInfo | null>(null);
  const [screen, setScreen] = useState<{ w: number; h: number }>({ w: 160, h: 144 });
  const [auto, setAuto] = useState(false);
  const [mode, setMode] = useState('rules');
  const [palette, setPalette] = useState('greyscale');
  const [paused, setPaused] = useState(false);

  autoRef.current = auto;

  const drawFrame = (dataUrl: string) => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext('2d');
    if (!ctx) return;
    const img = new Image();
    img.onload = () => {
      ctx.imageSmoothingEnabled = false;
      ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
    };
    img.src = dataUrl;
  };

  const socket = useEmulatorSocket(id, drawFrame);

  useEffect(() => {
    if (!id) return;
    getSession(id)
      .then((s) => {
        setInfo(s);
        setAuto(s.auto);
        setMode(s.mode);
        setPaused(s.paused);
        if (s.width && s.height) setScreen({ w: s.width, h: s.height });
      })
      .catch(() => undefined);
  }, [id]);

  // 手动模式键盘输入
  useEffect(() => {
    const down = (e: KeyboardEvent) => {
      const btn = KEYMAP[e.key];
      if (!btn || autoRef.current) return;
      if (['ArrowUp', 'ArrowDown', 'ArrowLeft', 'ArrowRight', 'Backspace', ' '].includes(e.key)) {
        e.preventDefault();
      }
      if (heldRef.current.has(btn)) return;
      heldRef.current.add(btn);
      socket.sendKeys([btn], []);
    };
    const up = (e: KeyboardEvent) => {
      const btn = KEYMAP[e.key];
      if (!btn) return;
      if (heldRef.current.has(btn)) {
        heldRef.current.delete(btn);
        socket.sendKeys([], [btn]);
      }
    };
    window.addEventListener('keydown', down);
    window.addEventListener('keyup', up);
    return () => {
      window.removeEventListener('keydown', down);
      window.removeEventListener('keyup', up);
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [socket.sendKeys]);

  const toggleAuto = (next: boolean) => {
    setAuto(next);
    socket.sendConfig({ auto: next, mode });
    message.info(next ? '已切换到 Agent 自动' : '已切换到手动');
  };

  const changeMode = (next: string) => {
    setMode(next);
    socket.sendConfig({ auto, mode: next });
  };

  const changePalette = (next: string) => {
    setPalette(next);
    socket.sendConfig({ palette: next });
  };

  const control = (action: 'reset' | 'pause' | 'resume') => {
    socket.sendControl(action);
    if (action === 'pause') setPaused(true);
    if (action === 'resume') setPaused(false);
  };

  const agent = socket.agent || {};
  const num = (v: unknown, digits = 0): string =>
    typeof v === 'number' ? v.toFixed(digits) : v === undefined || v === null ? '-' : String(v);

  return (
    <PageContainer
      title={`会话 ${id ?? '-'}`}
      onBack={() => navigate('/emulator/sessions')}
      extra={[
        <Badge
          key="conn"
          status={socket.connected ? 'processing' : 'default'}
          text={socket.connected ? '实时连接' : '未连接'}
        />,
      ]}
    >
      <Space align="start" size="large" wrap>
        <Card title="画面" styles={{ body: { padding: 12 } }}>
          <canvas
            ref={canvasRef}
            width={screen.w * 3}
            height={screen.h * 3}
            style={{
              imageRendering: 'pixelated',
              background: '#000',
              border: '2px solid #222',
              maxWidth: '100%',
            }}
          />
        </Card>

        <Space direction="vertical" size="middle" style={{ minWidth: 360 }}>
          <Card title="控制">
            <Space direction="vertical" size="small" style={{ width: '100%' }}>
              <Space>
                <span>Agent 自动</span>
                <Switch checked={auto} onChange={toggleAuto} />
                <Select
                  value={mode}
                  style={{ width: 110 }}
                  onChange={changeMode}
                  options={[
                    { label: '规则', value: 'rules' },
                    { label: '混合', value: 'hybrid' },
                    { label: 'LLM', value: 'llm' },
                    { label: '手动', value: 'manual' },
                  ]}
                />
              </Space>
              <Space>
                <span>调色板</span>
                <Select
                  value={palette}
                  style={{ width: 130 }}
                  onChange={changePalette}
                  options={[
                    { label: '灰阶', value: 'greyscale' },
                    { label: '经典绿', value: 'original' },
                    { label: 'BGB 绿', value: 'bgb' },
                  ]}
                />
              </Space>
              <Space>
                <Button onClick={() => control('reset')}>Reset</Button>
                {paused ? (
                  <Button type="primary" onClick={() => control('resume')}>
                    Resume
                  </Button>
                ) : (
                  <Button onClick={() => control('pause')}>Pause</Button>
                )}
                <Button
                  danger
                  onClick={async () => {
                    if (id) await stopSession(id);
                    navigate('/emulator/sessions');
                  }}
                >
                  停止会话
                </Button>
              </Space>
              <Typography.Text type="secondary">
                手动模式按键：方向键 · Z=A · X=B · Enter=Start · Backspace=Select · Q=L · W=R
              </Typography.Text>
            </Space>
          </Card>

          <Card title="Agent">
            <Descriptions column={2} size="small">
              <Descriptions.Item label="决策层">{String(agent.mode ?? mode)}</Descriptions.Item>
              <Descriptions.Item label="进度">{num(agent.maxProgress)}</Descriptions.Item>
              <Descriptions.Item label="生命">{num(agent.lives)}</Descriptions.Item>
              <Descriptions.Item label="金币">{num(agent.coins)}</Descriptions.Item>
              <Descriptions.Item label="状态">{String(agent.state ?? '-')}</Descriptions.Item>
              <Descriptions.Item label="敌人距离">{num(agent.enemyDX)}</Descriptions.Item>
              <Descriptions.Item label="死亡">{num(agent.deaths)}</Descriptions.Item>
              <Descriptions.Item label="奖励">{num(agent.rewardTotal, 1)}</Descriptions.Item>
              <Descriptions.Item label="当前状态">
                {paused ? <Tag color="warning">暂停</Tag> : <Tag color="success">运行</Tag>}
              </Descriptions.Item>
              <Descriptions.Item label="卡带">{info?.cart ?? '-'}</Descriptions.Item>
              <Descriptions.Item label="LLM 决策" span={2}>
                {String(agent.decisionRationale ?? '-')}
              </Descriptions.Item>
            </Descriptions>
          </Card>
        </Space>

        <Card title="日志" style={{ width: 380, height: 560, overflow: 'auto' }}>
          <div style={{ fontFamily: 'monospace', fontSize: 12, lineHeight: 1.6 }}>
            {socket.logs.length === 0 ? (
              <Typography.Text type="secondary">暂无日志</Typography.Text>
            ) : (
              socket.logs.map((l, i) => (
                <div key={i}>
                  [{l.time}] {l.msg}
                </div>
              ))
            )}
          </div>
        </Card>
      </Space>
    </PageContainer>
  );
}
