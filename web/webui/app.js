"use strict";

const canvas = document.getElementById("screen");
const ctx = canvas.getContext("2d");

const el = {
  conn: document.getElementById("conn"),
  cart: document.getElementById("cart"),
  fps: document.getElementById("fps"),
  overlay: document.getElementById("overlay"),
  overlayMsg: document.getElementById("overlayMsg"),
  log: document.getElementById("log"),
  btnAuto: document.getElementById("btnAuto"),
  selMode: document.getElementById("selMode"),
  btnReset: document.getElementById("btnReset"),
  btnPause: document.getElementById("btnPause"),
  btnMute: document.getElementById("btnMute"),
  selPalette: document.getElementById("selPalette"),
  sFrame: document.getElementById("sFrame"),
  sPC: document.getElementById("sPC"),
  sSP: document.getElementById("sSP"),
  sAF: document.getElementById("sAF"),
  sBC: document.getElementById("sBC"),
  sDE: document.getElementById("sDE"),
  sHL: document.getElementById("sHL"),
  sLY: document.getElementById("sLY"),
  sLCDC: document.getElementById("sLCDC"),
  sIFIE: document.getElementById("sIFIE"),
  aMode: document.getElementById("aMode"),
  aProgress: document.getElementById("aProgress"),
  aLives: document.getElementById("aLives"),
  aCoins: document.getElementById("aCoins"),
  aState: document.getElementById("aState"),
  aEnemy: document.getElementById("aEnemy"),
  aDeaths: document.getElementById("aDeaths"),
  aReward: document.getElementById("aReward"),
  aDecision: document.getElementById("aDecision"),
};

const KEYMAP = {
  ArrowUp: "Up",
  ArrowDown: "Down",
  ArrowLeft: "Left",
  ArrowRight: "Right",
  z: "A",
  Z: "A",
  x: "B",
  X: "B",
  Enter: "Start",
  Backspace: "Select",
  q: "L",
  Q: "L",
  w: "R",
  W: "R",
};

let ws = null;
let reconnectDelay = 1000;
let held = new Set();          // buttons currently held
let manual = true;             // false = auto (agent)
let paused = false;
let lastFrameAt = 0;
let lastTick = 0;
let frameRate = 0;

// --- WebSocket ----------------------------------------------------------------

function connect() {
  const proto = location.protocol === "https:" ? "wss" : "ws";
  ws = new WebSocket(`${proto}://${location.host}/ws`);
  ws.binaryType = "arraybuffer";

  ws.onopen = () => {
    el.conn.textContent = "已连接";
    el.conn.classList.add("ok");
    el.conn.classList.remove("off");
    reconnectDelay = 1000;
  };

  ws.onclose = () => {
    el.conn.textContent = "未连接";
    el.conn.classList.remove("ok");
    el.conn.classList.add("off");
    setTimeout(connect, reconnectDelay);
    reconnectDelay = Math.min(reconnectDelay * 2, 10000);
  };

  ws.onmessage = (e) => {
    // 二进制帧：magic(4)|width(2)|height(2)|tick(8)|RGBA，直接 putImageData。
    if (e.data instanceof ArrayBuffer) {
      drawBinaryFrame(e.data);
      return;
    }
    let msg;
    try {
      msg = JSON.parse(e.data);
    } catch (_) {
      return;
    }
    switch (msg.type) {
      case "hello":
        el.cart.textContent = msg.cart || "-";
        applyScreenSize(msg.width, msg.height);
        log("info", `已连接，ROM: ${msg.cart || "-"} console=${msg.console || "gb"} fps=${msg.fps}`);
        break;
      case "frame":
        drawFrame(msg.img, msg.tick);
        break;
      case "audio":
        playPcm(msg.pcm, msg.rate);
        break;
      case "state":
        updateState(msg.state);
        updateAgent(msg.agent || null);
        break;
      case "log":
        log(msg.level || "info", msg.msg || "");
        break;
    }
  };
}

function send(obj) {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(obj));
  }
}

// --- Rendering --------------------------------------------------------------

let pendingFrame = null;
let frameRaf = null;

const BINARY_MAGIC = [0x41, 0x47, 0x46, 0x4d]; // 'AGFM'

function drawBinaryFrame(buf) {
  if (buf.byteLength < 16) return;
  const bytes = new Uint8Array(buf);
  for (let i = 0; i < 4; i++) {
    if (bytes[i] !== BINARY_MAGIC[i]) return;
  }
  const view = new DataView(buf);
  const width = view.getUint16(4, true);
  const height = view.getUint16(6, true);
  const tick = Number(view.getBigUint64(8, true));
  if (width <= 0 || height <= 0 || buf.byteLength < 16 + width * height * 4) return;
  pendingFrame = { rgba: new Uint8ClampedArray(buf, 16, width * height * 4), width, height };
  lastTick = tick;
  if (frameRaf === null) {
    frameRaf = requestAnimationFrame(flushFrame);
  }
}

function drawFrame(imgBase64, tick) {
  pendingFrame = { imgBase64 };
  if (frameRaf === null) {
    frameRaf = requestAnimationFrame(flushFrame);
  }
  lastTick = tick;
}

function flushFrame() {
  frameRaf = null;
  const frame = pendingFrame;
  pendingFrame = null;
  if (!frame) return;

  if (frame.rgba) {
    if (canvas.width !== frame.width) canvas.width = frame.width;
    if (canvas.height !== frame.height) canvas.height = frame.height;
    const imgData = new ImageData(frame.rgba, frame.width, frame.height);
    ctx.putImageData(imgData, 0, 0);
  } else if (frame.imgBase64) {
    const img = new Image();
    img.onload = () => {
      ctx.imageSmoothingEnabled = false;
      ctx.drawImage(img, 0, 0, canvas.width, canvas.height);
    };
    img.src = "data:image/png;base64," + frame.imgBase64;
  }

  const now = performance.now();
  if (lastFrameAt) {
    const dt = (now - lastFrameAt) / 1000;
    frameRate = Math.round(1 / dt);
  }
  lastFrameAt = now;
}

function applyScreenSize(w, h) {
  if (!w || !h) return;
  canvas.width = w;
  canvas.height = h;
  const scale = Math.max(1, Math.floor(576 / h));
  canvas.style.width = w * scale + "px";
  canvas.style.height = h * scale + "px";
}

function updateState(s) {
  el.sFrame.textContent = s.frame;
  el.sPC.textContent = hex4(s.cpu.pc);
  el.sSP.textContent = hex4(s.cpu.sp);
  el.sAF.textContent = hex4(s.cpu.af);
  el.sBC.textContent = hex4(s.cpu.bc);
  el.sDE.textContent = hex4(s.cpu.de);
  el.sHL.textContent = hex4(s.cpu.hl);
  el.sLY.textContent = hex2(s.ly);
  el.sLCDC.textContent = bits(s.lcdc);
  el.sIFIE.textContent = hex2(s.if) + " / " + hex2(s.ie);
  el.fps.textContent = frameRate + " fps · tick " + lastTick;

  const wasPaused = paused;
  paused = !!s.paused;
  if (paused !== wasPaused) {
    el.overlayMsg.textContent = paused ? "PAUSED" : "";
    el.overlay.classList.toggle("show", paused);
  }
}

function hex4(v) { return "0x" + (v >>> 0).toString(16).toUpperCase().padStart(4, "0"); }
function hex2(v) { return "0x" + (v & 0xff).toString(16).toUpperCase().padStart(2, "0"); }
function bits(v) { return (v & 0xff).toString(2).padStart(8, "0"); }

function updateAgent(a) {
  const show = a ? "on" : "off";
  document.body.classList.remove("agent-off", "agent-on");
  document.body.classList.add("agent-" + show);
  if (!a) {
    el.aMode.textContent = "手动";
    el.aProgress.textContent = el.aLives.textContent = el.aCoins.textContent =
      el.aState.textContent = el.aEnemy.textContent = el.aDeaths.textContent =
      el.aReward.textContent = el.aDecision.textContent = "-";
    return;
  }
  el.aMode.textContent = a.mode || "-";
  el.aProgress.textContent = "camera=" + (a.progress !== undefined ? a.progress : "-") +
    " / max=" + (a.maxProgress !== undefined ? a.maxProgress : "-");
  el.aLives.textContent = a.lives;
  el.aCoins.textContent = a.coins;
  el.aState.textContent = (a.state || "-") + (a.hard ? " · hard" : "") + (a.demo ? " · demo" : "");
  if (a.enemyDX !== undefined && a.enemyDX !== null) {
    el.aEnemy.textContent = "dx=" + a.enemyDX;
  } else {
    el.aEnemy.textContent = "无";
  }
  el.aDeaths.textContent = a.deaths !== undefined ? a.deaths : "-";
  el.aReward.textContent = a.rewardTotal !== undefined ? a.rewardTotal.toFixed(1) : "-";
  el.aDecision.textContent = a.decisionRationale || "-";
}

function log(level, text) {
  const div = document.createElement("div");
  div.className = level;
  div.textContent = `[${new Date().toLocaleTimeString()}] ${text}`;
  el.log.appendChild(div);
  while (el.log.children.length > 200) el.log.removeChild(el.log.firstChild);
  el.log.scrollTop = el.log.scrollHeight;
}

// --- Audio -------------------------------------------------------------------

let audioCtx = null;
let audioNextTime = 0;
let muted = false;

function ensureAudio() {
  if (!audioCtx) {
    const Ctor = window.AudioContext || window.webkitAudioContext;
    if (!Ctor) return null;
    audioCtx = new Ctor();
  }
  if (audioCtx.state === "suspended") audioCtx.resume();
  return audioCtx;
}

function base64ToBytes(b64) {
  const bin = atob(b64);
  const bytes = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; i++) bytes[i] = bin.charCodeAt(i);
  return bytes;
}

function playPcm(b64, rate) {
  if (muted || !b64) return;
  const ctx = ensureAudio();
  if (!ctx) return;
  const bytes = base64ToBytes(b64);
  const samples = bytes.length >> 1;
  const frames = samples >> 1;
  if (frames <= 0) return;
  const i16 = new Int16Array(bytes.buffer, bytes.byteOffset, samples);
  const buf = ctx.createBuffer(2, frames, rate || 32768);
  const l = buf.getChannelData(0);
  const r = buf.getChannelData(1);
  for (let i = 0; i < frames; i++) {
    l[i] = i16[2 * i] / 32768;
    r[i] = i16[2 * i + 1] / 32768;
  }
  const src = ctx.createBufferSource();
  src.buffer = buf;
  src.connect(ctx.destination);
  const now = ctx.currentTime;
  if (audioNextTime < now + 0.02) audioNextTime = now + 0.08;
  src.start(audioNextTime);
  audioNextTime += buf.duration;
}

// --- Input -------------------------------------------------------------------

function buttonName(key) {
  return KEYMAP[key];
}

function press(btn) {
  if (held.has(btn)) return;
  held.add(btn);
  send({ type: "keys", pressed: [btn], released: [] });
}

function release(btn) {
  if (!held.has(btn)) return;
  held.delete(btn);
  send({ type: "keys", pressed: [], released: [btn] });
}

window.addEventListener("keydown", (e) => {
  const btn = buttonName(e.key);
  if (!btn) return;
  if (["ArrowUp", "ArrowDown", "ArrowLeft", "ArrowRight", "Backspace", " "].includes(e.key)) {
    e.preventDefault();
  }
  if (manual) press(btn);
});

window.addEventListener("keyup", (e) => {
  const btn = buttonName(e.key);
  if (!btn) return;
  release(btn);
});

window.addEventListener("blur", () => {
  [...held].forEach(release);
});

// --- Controls ----------------------------------------------------------------

el.btnReset.addEventListener("click", () => {
  held.forEach(release);
  send({ type: "control", action: "reset" });
  log("info", "请求 Reset");
});

el.btnPause.addEventListener("click", () => {
  const action = paused ? "resume" : "pause";
  send({ type: "control", action });
  log("info", paused ? "继续" : "暂停");
});

el.btnMute.addEventListener("click", () => {
  muted = !muted;
  el.btnMute.textContent = muted ? "Unmute" : "Mute";
  el.btnMute.classList.toggle("primary", muted);
  ensureAudio();
});

el.btnAuto.addEventListener("click", () => {
  manual = !manual;
  el.btnAuto.textContent = manual ? "Auto" : "Manual";
  el.btnAuto.classList.toggle("primary", !manual);
  if (!manual) {
    const mode = el.selMode.value;
    held.forEach(release);
    send({ type: "config", auto: true, agent: mode });
    log("info", `已切换 Auto（Agent 模式: ${mode}）`);
  } else {
    send({ type: "config", auto: false, agent: "manual" });
    log("info", "已切换 Manual，键盘接管");
  }
});

el.selPalette.addEventListener("change", () => {
  send({ type: "config", palette: el.selPalette.value });
  log("info", `调色板 -> ${el.selPalette.selectedOptions[0].textContent}`);
});

connect();