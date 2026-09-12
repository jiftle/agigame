# GoBoy-LLM

AI 玩 Game Boy 实验平台。「GB 核心 + WebSocket 服务 + WebUI + Agent」四层架构。

```
┌─────────────────────────────────────────────────────┐
│  Layer 5: WebUI  canvas + 状态面板 + 键盘 + 日志     │
├─────────────────────────────────────────────────────┤
│  Layer 4: Server   HTTP + WO 消息 / WS Hub / 静态文件 │
├─────────────────────────────────────────────────────┤
│  Layer 3: Agent   规则(高频) + LLM(低频) + Reward     │  ← P2/P3
├─────────────────────────────────────────────────────┤
│  Layer 2: core/gb CPU/PPU/APU/MMU/Cart + 回调注入     │
├─────────────────────────────────────────────────────┤
│  Layer 1: ROM     用户自有合法 dump                   │
└─────────────────────────────────────────────────────┘
```

当前进度：**P0（核心改造）+ P1（Web 可玩）已完成**。

- P0：`core/gb` 解耦出 `FrameCallback` / `StateCallback` / `InputProvider`，主循环 `Run(ctx)/Step()`
- P1：`server/`（HTTP + WS + Hub + 输入聚合）+ `webui/`（canvas 渲染 + 键盘 + 控制面板）
- P2：`agent/` 规则 Agent + Mario 地址表（规划中）
- P3：LLM 决策层（规划中）

## 快速开始

1. 放入你的 Super Mario Land ROM（仅支持自有合法 dump）：

   ```bash
   mkdir -p roms
   cp /path/to/Super_Mario_Land.gb roms/
   ```

2. 修改 `config.yaml` 中的 `rom` 路径（默认 `./roms/super_mario_land.gb`）。

3. 运行：

   ```bash
   go run ./cmd/server -config config.yaml
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

页面按钮：`Reset` 重置、`Pause` 暂停/继续、`Auto/Manual` 切换（Auto 需 P2 的 Agent）。

## WS 协议

Server → Client：

```json
{"type":"hello","cart":"SUPER MARIO LAND","fps":60,"game":"sml"}
{"type":"frame","img":"<base64 PNG>","tick":123}
{"type":"state","state":{"frame":123,"cpu":{"pc":"0x1234",...},"ly":0,...}}
{"type":"log","level":"info","msg":"..."}
```

Client → Server：

```json
{"type":"keys","pressed":["A","Right"],"released":["Start"]}
{"type":"control","action":"reset"}                    // reset | pause | resume
{"type":"config","auto":true,"agent":"rules"}          // P2 起生效
```

## 目录结构

```
cmd/server         入口：装配各层、启动 loop
core/gb            GameBoy 核心（GoBoy 源码改造：回调注入 + 运行循环 + 快照）
core/cart          MBC1/2/3/5、ROM、RAM+电池存档
core/apu           无头 APU（保留寄存器语义，不产生音频）
server             HTTP + WS Hub + 输入聚合 + 帧编码(JSON)
webui              canvas 前端
config.yaml        配置
```

## 设计要点

- **核心零依赖**：`core/gb` 不依赖 UI/网络库；`core/apu` 为纯 Go 无头实现（无需 cgo）。
- **回调注入**：渲染通过 `SetFrameCallback`，输入通过 `SetInputProvider`，快照通过 `SetStateCallback`。
- **不阻塞精化**：帧回调在仿真 goroutine 上只做拷贝入队（最新帧优先），PNG 编码在独立 goroutine。
- **只接受本地 ROM**：不做任何网络下载 / 分发。

## 测试

```bash
go test ./...
```

- `core/gb/smoke_test.go`：无头跑帧、回调、输入 diff
- `server/server_test.go`：HTTP 启动 + WS 收到 hello/frame/state 的端到端测试