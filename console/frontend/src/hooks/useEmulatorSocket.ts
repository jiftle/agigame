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

type FrameSink = (frame: FrameData) => void;

/** 一帧画面：优先使用原始 RGBA，否则回退到 PNG dataURL。 */
export interface FrameData {
  rgba?: Uint8ClampedArray<ArrayBuffer>;
  width?: number;
  height?: number;
  dataUrl?: string;
}
type AudioSink = (pcmBase64: string, rate: number) => void;

const BINARY_MAGIC = [0x41, 0x47, 0x46, 0x4d]; // 'AGFM'

/** 解析二进制帧：magic(4)|width(2)|height(2)|tick(8)|RGBA。 */
function parseBinaryFrame(buf: ArrayBuffer): FrameData | null {
  if (buf.byteLength < 16) return null;
  const bytes = new Uint8Array(buf);
  for (let i = 0; i < 4; i++) {
    if (bytes[i] !== BINARY_MAGIC[i]) return null;
  }
  const view = new DataView(buf);
  const width = view.getUint16(4, true);
  const height = view.getUint16(6, true);
  const expected = 16 + width * height * 4;
  if (width <= 0 || height <= 0 || buf.byteLength < expected) return null;
  return { rgba: new Uint8ClampedArray(buf, 16, width * height * 4), width, height };
}

/** 会话 WS 地址：优先直连 VITE_WS_TARGET（dev 绕过 Vite 代理），否则同源。 */
function buildStreamURL(sessionId: string): string {
  const path = `/ws/emu/${sessionId}/stream?token=${encodeURIComponent(getToken())}`;
  const target = import.meta.env.VITE_WS_TARGET as string | undefined;
  if (target) {
    const u = new URL(target, window.location.href);
    const proto = u.protocol === 'https:' ? 'wss' : 'ws';
    return `${proto}://${u.host}${path}`;
  }
  const proto = window.location.protocol === 'https:' ? 'wss' : 'ws';
  return `${proto}://${window.location.host}${path}`;
}

/** 建立会话 WebSocket：接收 frame/state/log/audio，下发 keys/control/config。 */
export function useEmulatorSocket(
  sessionId: string | undefined,
  onFrame: FrameSink,
  onAudio?: AudioSink,
): EmuSocketHandle {
  const wsRef = useRef<WebSocket | null>(null);
  const onFrameRef = useRef<FrameSink>(onFrame);
  onFrameRef.current = onFrame;
  const onAudioRef = useRef<AudioSink | undefined>(onAudio);
  onAudioRef.current = onAudio;

  const [connected, setConnected] = useState(false);
  const [agent, setAgent] = useState<Record<string, unknown> | null>(null);
  const [logs, setLogs] = useState<EmuLogEntry[]>([]);

  useEffect(() => {
    if (!sessionId) return;
    const ws = new WebSocket(buildStreamURL(sessionId));
    ws.binaryType = 'arraybuffer';
    wsRef.current = ws;

    ws.onopen = () => setConnected(true);
    ws.onclose = () => setConnected(false);
    ws.onerror = () => setConnected(false);
    ws.onmessage = (ev) => {
      // 二进制帧：直接解析为 RGBA，无 JSON/base64/PNG 解码。
      if (ev.data instanceof ArrayBuffer) {
        const frame = parseBinaryFrame(ev.data);
        if (frame) onFrameRef.current(frame);
        return;
      }

      let msg: { type?: string } & Record<string, unknown>;
      try {
        msg = JSON.parse(ev.data as string);
      } catch {
        return;
      }
      switch (msg.type) {
        case 'frame':
          if (typeof msg.img === 'string') {
            onFrameRef.current({ dataUrl: `data:image/png;base64,${msg.img}` });
          }
          break;
        case 'state':
          setAgent((msg.agent as Record<string, unknown>) || null);
          break;
        case 'audio':
          onAudioRef.current?.(
            String(msg.pcm || ''),
            typeof msg.rate === 'number' ? msg.rate : 32768,
          );
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
