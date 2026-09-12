/** 基于 Web Audio API 的 PCM（立体声 s16le）排队播放器。 */
export class PcmPlayer {
  private ctx: AudioContext | null = null;
  private nextTime = 0;
  muted = false;

  private ensureContext(): AudioContext | null {
    if (!this.ctx) {
      const Ctor =
        window.AudioContext || (window as unknown as { webkitAudioContext?: typeof AudioContext }).webkitAudioContext;
      if (!Ctor) return null;
      this.ctx = new Ctor();
    }
    if (this.ctx.state === 'suspended') {
      void this.ctx.resume();
    }
    return this.ctx;
  }

  /** 在用户手势里调用以解锁音频（浏览器自动播放策略）。 */
  resume(): void {
    this.ensureContext();
  }

  playPcm(base64: string, rate: number): void {
    if (this.muted || !base64) return;
    const ctx = this.ensureContext();
    if (!ctx) return;

    const bin = atob(base64);
    const bytes = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; i++) bytes[i] = bin.charCodeAt(i);

    const samples = bytes.length >> 1;
    const frames = samples >> 1;
    if (frames <= 0) return;

    const i16 = new Int16Array(bytes.buffer, bytes.byteOffset, samples);
    const buffer = ctx.createBuffer(2, frames, rate || 32768);
    const left = buffer.getChannelData(0);
    const right = buffer.getChannelData(1);
    for (let i = 0; i < frames; i++) {
      left[i] = i16[2 * i] / 32768;
      right[i] = i16[2 * i + 1] / 32768;
    }

    const src = ctx.createBufferSource();
    src.buffer = buffer;
    src.connect(ctx.destination);

    const now = ctx.currentTime;
    if (this.nextTime < now + 0.02) {
      this.nextTime = now + 0.08; // 小抖动缓冲
    }
    src.start(this.nextTime);
    this.nextTime += buffer.duration;
  }
}
