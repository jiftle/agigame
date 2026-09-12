import { useCallback, useEffect, useRef, useState } from 'react';

import { getToken } from '@/stores/auth';

import type { EmuConfigInput } from '@/services/emulator';

export interface EmuLogEntry {
  time: string;
  level: string;
  msg: string;
}

export interface EmuSocketHandle {
  connected: boolean;
  agent: Record<string, unknown> | null;
  logs: EmuLogEntry[];
  sendKeys: (pressed: string[], released: string[]) => void;
  sendControl: (action: string) => void;
  sendConfig: (cfg: EmuConfigInput) => void;
}

type FrameSink = (dataUrl: string) => void;

/** 建立会话 WebSocket：接收 frame/state/log，下发 keys/control/config。 */
export function useEmulatorSocket(
  sessionId: string | undefined,
  onFrame: FrameSink,
): EmuSocketHandle {
  const wsRef = useRef<WebSocket | null>(null);
  const onFrameRef = useRef<FrameSink>(onFrame);
  onFrameRef.current = onFrame;

  const [connected, setConnected] = useState(false);
  const [agent, setAgent] = useState<Record<string, unknown> | null>(null);
  const [logs, setLogs] = useState<EmuLogEntry[]>([]);

  useEffect(() => {
    if (!sessionId) return;
    const proto = window.location.protocol === 'https:' ? 'wss' : 'ws';
    const token = getToken();
    const url = `${proto}://${window.location.host}/ws/emu/${sessionId}/stream?token=${encodeURIComponent(token)}`;
    const ws = new WebSocket(url);
    wsRef.current = ws;

    ws.onopen = () => setConnected(true);
    ws.onclose = () => setConnected(false);
    ws.onerror = () => setConnected(false);
    ws.onmessage = (ev) => {
      let msg: { type?: string } & Record<string, unknown>;
      try {
        msg = JSON.parse(ev.data as string);
      } catch {
        return;
      }
      switch (msg.type) {
        case 'frame':
          if (typeof msg.img === 'string') {
            onFrameRef.current(`data:image/png;base64,${msg.img}`);
          }
          break;
        case 'state':
          setAgent((msg.agent as Record<string, unknown>) || null);
          break;
        case 'log':
          setLogs((prev) => {
            const next = [
              ...prev,
              {
                time: new Date().toLocaleTimeString(),
                level: String(msg.level || 'info'),
                msg: String(msg.msg || ''),
              },
            ];
            return next.length > 300 ? next.slice(next.length - 300) : next;
          });
          break;
        default:
          break;
      }
    };

    return () => {
      ws.close();
      wsRef.current = null;
      setConnected(false);
    };
  }, [sessionId]);

  const send = useCallback((payload: Record<string, unknown>) => {
    const ws = wsRef.current;
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(payload));
    }
  }, []);

  const sendKeys = useCallback(
    (pressed: string[], released: string[]) => {
      send({ type: 'keys', pressed, released });
    },
    [send],
  );
  const sendControl = useCallback((action: string) => send({ type: 'control', action }), [send]);
  const sendConfig = useCallback(
    (cfg: EmuConfigInput) =>
      send({ type: 'config', auto: cfg.auto, agent: cfg.mode, palette: cfg.palette }),
    [send],
  );

  return { connected, agent, logs, sendKeys, sendControl, sendConfig };
}
