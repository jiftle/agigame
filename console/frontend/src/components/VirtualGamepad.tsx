import type { CSSProperties, PointerEvent as ReactPointerEvent, MouseEvent as ReactMouseEvent } from 'react';

export interface VirtualGamepadProps {
  /** 按下某键（值为 A/B/Start/Select/Up/Down/Left/Right/L/R）。 */
  onPress: (button: string) => void;
  /** 松开某键。 */
  onRelease: (button: string) => void;
  /** Agent 自动模式下禁用手柄。 */
  disabled?: boolean;
  /** GBA 等带肩键的平台显示 L/R。 */
  showShoulders?: boolean;
}

const baseBtn: CSSProperties = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  background: '#232834',
  border: '1px solid #3a4150',
  color: '#d7dce4',
  borderRadius: 8,
  fontSize: 15,
  fontWeight: 600,
  userSelect: 'none',
  WebkitUserSelect: 'none',
  touchAction: 'none',
  cursor: 'pointer',
  padding: 0,
};

const dpadCell: CSSProperties = { ...baseBtn, width: 52, height: 52, fontSize: 18 };
const roundBtn: CSSProperties = { ...baseBtn, width: 56, height: 56, borderRadius: 28 };
const pillBtn: CSSProperties = { ...baseBtn, height: 30, minWidth: 64, borderRadius: 15, fontSize: 12 };
const shoulderBtn: CSSProperties = { ...baseBtn, height: 30, width: 96, borderRadius: 10, fontSize: 13 };

export default function VirtualGamepad({
  onPress,
  onRelease,
  disabled = false,
  showShoulders = false,
}: VirtualGamepadProps) {
  const bind = (button: string) => ({
    onPointerDown: (e: ReactPointerEvent<HTMLButtonElement>) => {
      e.preventDefault();
      if (disabled) return;
      e.currentTarget.setPointerCapture(e.pointerId);
      onPress(button);
    },
    onPointerUp: (e: ReactPointerEvent<HTMLButtonElement>) => {
      e.preventDefault();
      if (disabled) return;
      onRelease(button);
    },
    onPointerCancel: () => {
      if (!disabled) onRelease(button);
    },
    onContextMenu: (e: ReactMouseEvent) => e.preventDefault(),
  });

  const key = (
    label: string,
    button: string,
    style: CSSProperties,
    extra?: CSSProperties,
  ) => (
    <button
      type="button"
      disabled={disabled}
      {...bind(button)}
      style={{ ...style, ...(disabled ? { opacity: 0.4, cursor: 'not-allowed' } : {}), ...extra }}
    >
      {label}
    </button>
  );

  return (
    <div style={{ marginTop: 12, opacity: disabled ? 0.7 : 1 }}>
      {showShoulders && (
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 10 }}>
          {key('L', 'L', shoulderBtn)}
          {key('R', 'R', shoulderBtn)}
        </div>
      )}

      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', gap: 16 }}>
        {/* 十字键 */}
        <div
          style={{
            display: 'grid',
            gridTemplateColumns: 'repeat(3, 52px)',
            gridTemplateRows: 'repeat(3, 52px)',
            gap: 4,
          }}
        >
          <span />
          {key('▲', 'Up', dpadCell)}
          <span />
          {key('◀', 'Left', dpadCell)}
          <span />
          {key('▶', 'Right', dpadCell)}
          <span />
          {key('▼', 'Down', dpadCell)}
          <span />
        </div>

        {/* A/B + Start/Select */}
        <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 10 }}>
          <div style={{ display: 'flex', gap: 12, alignItems: 'center' }}>
            {key('B', 'B', roundBtn)}
            {key('A', 'A', roundBtn)}
          </div>
          <div style={{ display: 'flex', gap: 10 }}>
            {key('Select', 'Select', pillBtn)}
            {key('Start', 'Start', pillBtn)}
          </div>
        </div>
      </div>
    </div>
  );
}
