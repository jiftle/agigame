# GoBoy-LLM

AI 玩 Game Boy 实验平台。仓库分三大块：**模拟器核心 / Web 模拟器前端 / 控制台**。

```
agigame/
├── emulator/           模拟器核心（纯 Go 引擎）
│   ├── core/           CPU/PPU/APU/MMU/Cart + 回调注入
│   ├── agent/          规则(高频) + LLM(低频) + Reward + games 插件
│   └── engine/         传输无关会话层（Web 与控制台共用）
├── web/                Web 模拟器前端（独立可玩）
│   ├── server/         HTTP + WS Hub + 输入聚合 + 帧编码
│   └── webui/          canvas 渲染 + 面板
└── console/            控制台（AdminBase 改造适配，规划中）
    ├── backend/        GoFrame v2
    └── frontend/       React 19 + Ant Design
```

当前进度：**P0（核心改造）+ P1（Web 可玩）+ P2/P3（Agent 空架子）已完成**；目录已按三大块重构。

- P0：`emulator/core/gb` 解耦出 `FrameCallback` / `StateCallback` / `InputProvider` / `SetPreFrameCallback`，主循环 `Run(ctx)/Step()`
- P1：`web/server`（HTTP + WS + Hub + 输入聚合）+ `web/webui`（canvas 渲染 + 键盘 + 控制面板）
- P2：`emulator/agent/games` 规则 Agent（Super Mario Land 地址表 + 状态提取 + 规则决策 + Reward + SafetyNet）
- P3：`emulator/agent` LLM 决策层**空架子**：`LLMProvider` 接口 + `StubLLM`（不接真实 API）、异步 2Hz 决策循环、prompt 模板
- P4/P5：实验能力（SaveState / 批量 / reward 日志）与插件化（规划中）

## 快速开始

1. 放入你的 Super Mario Land ROM（仅支持自有合法 dump）：

   ```bash
   mkdir -p roms
   cp /path/to/Super_Mario_Land.gb roms/
   ```

2. 修改 `web/config.yaml` 中的 `rom` 路径（默认 `./roms/super-mario-land.gb`）。

3. 运行（在仓库根目录执行）：

   ```bash
   go run ./web/cmd/server -config web/config.yaml
   # 浏览器打开 http://localhost:8080
   ```

## 操作

| 键 | 手柄按钮 |
| --- | --- |
| `Z` | A |
| `X` | B |
| `Enter` | Start |
| `Backspace` | Select |
| `↑↓←→` | 方向键 |

页面按钮：`Reset` 重置、`Pause` 暂停/继续、`Auto/Manual` 切换（Auto 即启用 Agent，可在「规则/混合/LLM/手动」间选择决策层）。

## Agent（P2/P3 空架子）

点页面 `Auto` 后启用：

- **rules**：规则层每帧（emu goroutine）同步决策，SML 规则为「标题按 Start → 默认右行 → 敌人在 16px 内跳跃 → 空中保持右行」；`SafetyNet` 在无进度 120 帧后强制右行+脉冲跳跃脱困。
- **hybrid**：紧急情况用规则，其余交给 LLM 决策；**llm**：纯用最近一次 LLM 决策。
- 当前 LLM 一律是 `StubLLM`（模拟 30ms 延迟、返回 `["Right"]`），**不会调用任何外部 API**。P3 接真实 LLM 只需实现 `agent.LLMProvider` 并在 `engine.Config.Agent.Provider` 注入。
- Reward = camera 前进 +1、1UP +25、死亡 -50；统计（死亡/累计奖励/进度/最近决策）随 `state` 消息推送并在 WebUI 展示。

## WS 协议

Server → Client：

```json
{"type":"hello","cart":"SUPER MARIO LAND","fps":60,"game":"sml"}
{"type":"frame","img":"<base64 PNG>","tick":123}
{"type":"state","state":{...},"agent":{"mode":"hybrid","camera":12,"lives":3,...}} // auto 时附带 game 字段
{"type":"log","level":"info","msg":"..."}
```

Client → Server：

```json
{"type":"keys","pressed":["A","Right"],"released":["Start"]}
{"type":"control","action":"reset"}                    // reset | pause | resume
{"type":"config","auto":true,"agent":"hybrid"}         // manual | rules | hybrid | llm
```

## 目录结构

```
emulator/core/gb     GameBoy 核心（GoBoy 源码改造：回调注入 + 运行循环 + 快照 + 帧前钩子）
emulator/core/cart   MBC1/2/3/5、ROM、RAM+电池存档
emulator/core/apu    无头 APU（保留寄存器语义，不产生音频）
emulator/agent       决策层：模式切换 + 规则同步决策 + LLM 异步循环 + Reward/统计
emulator/agent/games 游戏插件：GamePlugin 接口 + SML 实现 + 注册表(Register/Get)
emulator/engine      传输无关会话层：生命周期 + 输入聚合 + 帧编码 + Agent 编排
web/server           engine 的 HTTP/WS 适配器（Hub + 协议 + 静态文件）
web/webui            canvas 前端
web/cmd/server       入口：装配各层、启动 loop
web/config.yaml      Web 模拟器配置
console/backend      控制台后端（GoFrame v2，独立 module，import emulator/engine）
console/frontend     控制台前端（React 19 + Ant Design）
roms/  saves/        ROM 与存档（gitignore）
```

## 控制台

控制台基于 AdminBase 改造适配（`console/`），提供登录/RBAC/日志，并新增「模拟器」会话管理：列表、启动、实时画面（WS）、Agent 面板、控制与调色板切换。后端通过 `replace agigame => ../..` 引用 `emulator/engine`。

```bash
make console-install     # 安装控制台前端依赖
make dev-console         # 仅启动控制台（后端 :8000 / 前端 :8001），默认 admin/123456
```

控制台 ROM 目录由 `console/backend/manifest/config/config.yaml` 的 `emulator.romDir` 配置（默认 `../../roms`）。启动会话时「选择游戏」会列出该目录下的 `.gb/.gbc/.gba`（读取 ROM 标题），并自动匹配平台。

## GBA 支持

除 DMG/CGB（GoBoy 核心）外，已支持 **Game Boy Advance**：

- 核心：[`aabalke/guac`](https://github.com/aabalke/guac) 的 `emu/gba`（纯 Go，BSD-3-Clause），适配层在 `emulator/gba/`，**无头运行**（无需窗口/BIOS/音频设备）。
- 启动：控制台「启动会话」里平台选 **Game Boy Advance**，ROM 填 `.gba` 文件名；Web 版则把 `web/config.yaml` 的 `emulator.console` 设为 `gba` 并把 `rom` 指向 `.gba`。
- 按键：A/B/Start/Select/方向键 + **Q=L、W=R**。
- 首期范围为「能玩」（画面/输入/控制台/WS）；**Agent（规则/LLM）暂只支持 GB**，GBA 会话忽略 auto/mode。

> 注：guac 会引入 ebiten 等依赖（已在根 module）。GBA 及其 ROM 仅支持自有合法 dump。

## 一键命令

```bash
make dev     # 一键调试启动：Web 模拟器(:8080) + 控制台(:8000/:8001)
make build   # 构建 Web 服务 + 控制台后端
make test    # 引擎与服务测试
make check   # gofmt 校验 + vet + test + 控制台构建
make dev-web # 仅启动 Web 模拟器（:8080）
```

## 设计要点

- **核心零依赖**：`emulator/core/gb` 不依赖 UI/网络库；`emulator/core/apu` 为纯 Go 无头实现（无需 cgo）。
- **回调注入**：渲染通过 `SetFrameCallback`，输入通过 `SetInputProvider`，快照通过 `SetStateCallback`，Agent 决策通过 `SetPreFrameCallback`（每帧、仿真 goroutine 上同步执行）。
- **决策分层**：规则在 emu goroutine 每帧同步跑；LLM 在独立 goroutine 异步（2Hz）只操作最近一次状态副本，写入原子安全区，绝不阻塞精化。
- **不阻塞精化**：帧回调在仿真 goroutine 上只做拷贝入队（最新帧优先），PNG 编码在独立 goroutine。
- **SML 地址表**：集中在 `emulator/agent/games/super_mario_land.go`，参考 ROM Detectives / Data Crystal，注意偏移可能随 ROM 修订版本（v1.0 vs JUE 1.1）变化，异常时优先检查该表。
- **只接受本地 ROM**：不做任何网络下载 / 分发。

## 测试

```bash
go test ./...
```

- `emulator/core/gb/smoke_test.go`：无头跑帧、回调、输入 diff
- `web/server/server_test.go`：HTTP 启动 + WS 收到 hello/frame/state 的端到端测试
- `emulator/agent/*_test.go`：模式切换、按键 diff 应用、stub LLM 决策接入、SafetyNet 救援
- `emulator/agent/games/*_test.go`：SML 地址提取（BCD 生命/时间、OAM 敌人扫描）、规则决策、Reward