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
    let msg;
    try {
      msg = JSON.parse(e.data);
    } catch (_) {
      return;
    }
    switch (msg.type) {
      case "hello":
        el.cart.textContent = msg.cart || "-";
        log("info", `已连接，ROM: ${msg.cart || "-"} fps=${msg.fps}`);
        break;
      case "frame":
        drawFrame(msg.img, msg.tick);
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

function drawFrame(imgBase64, tick) {
  const img = new Image();
  img.onload = () => {
    ctx.imageSmoothingEnabled = false;
    ctx.drawImage(img, 0, 0, 160, 144);
  };
  img.src = "data:image/png;base64," + imgBase64;

  const now = performance.now();
  if (lastFrameAt) {
    const dt = (now - lastFrameAt) / 1000;
    frameRate = Math.round(1 / dt);
  }
  lastFrameAt = now;
  lastTick = tick;
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

connect();