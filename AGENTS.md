# AGENTS.md — agigame 工程与操作约定

面向在本仓库工作的编码助手。项目结构与运行方式见 `README.md`；本文件只记录**硬性约定与容易踩的坑**。

## 常用命令

- `make dev` 一键调试：Web 模拟器(:8080) + 控制台(:8000/:8001)
- `make dev-web` / `make dev-console` 单独启动
- `make check` gofmt 校验 + `go vet` + `go test` + 控制台后端构建
- `make build` / `make test`
- 控制台前端：`cd console/frontend && pnpm tsc`（类型）；`pnpm build`
- 控制台后端独立 module：`cd console/backend && go build ./...`

## 提交规范

- 提交信息用**中文**；**一次提交主题单一**。
- 提交前：显式 `git add <文件>`（慎用未经核对的 `-A`），并用 `git diff --cached --stat` / `--name-only` 审阅改动清单。
- **提交后必须 `git show --stat HEAD` 核对文件清单**，不要只相信 commit message（本项目曾出现"提交信息写了 WS 服务，但 `server/` 整目录没进库"）。
- 纯格式改动（`gofmt`）单独提交，避免与功能改动混在一起。
- 不提交：密钥、ROM（`roms/*`）、存档（`saves/*`）、运行数据/日志/构建产物。

## Git 操作

- **本机 git 为 2.20.1（旧版）**，避免新 flag：
  - 取当前分支用 `git rev-parse --abbrev-ref HEAD`（**`git branch --show-current` 会报 `unknown option`**）。
  - 不要用 `git switch` / `git restore`（2.23+），用 `git checkout`。
- push 前先确认：`git remote -v`（是否配置了 remote）、`git status -sb`（分支与上游）、当前分支名。
  - 无 remote：`git remote add origin <url>`。
  - 分支名不一致：先 `git branch -M main` 再 `git push -u origin main`。
  - **未确认远程历史前，不要 `force push`**。
- 排查"文件是否被忽略 / 编辑为何不显示"：
  - `git status` 默认**不显示被忽略项**；用 `git status --ignored`。
  - `git check-ignore -v <path>` 定位是哪条规则；`git ls-files <dir>` 确认是否已被跟踪。
  - 教训：`.gitignore` 的 `/server` 既匹配根目录二进制、也匹配同名目录，曾导致整个 `server/` 源码目录被忽略、整包没提交。

## 构建产物与 .gitignore

- **产物不要与源码目录同名**：`go build -o server` 遇到同名 `server/` 目录会写到 `server/server`；产物统一放 `bin/`。
- 忽略规则用**精确路径**：`/server` 可能误伤同名源码目录；需要时写 `server/server`。
- 已在 `.gitignore`：`console/backend/{data,logs,bin}`、`console/frontend/{dist,node_modules}`、`console/.opencode/node_modules`、`roms/*`、`saves/*`、根二进制 `/server`。

## 验证纪律

- 关键校验要看**真实退出码**：`cmd | head; echo $?` 取到的是 `head` 的退出码，应使用 `${PIPESTATUS[0]}` 或不接管道。
- 提交前至少跑一次 `make check`；前端改动跑 `pnpm tsc`（必要时 `pnpm build`）。
- 引擎/服务的回归依赖 `go test ./...`（含 `emulator/gb`、`emulator/gba` 无头测试）。

## 进程管理（调试时）

- `pkill -f 'xxx'` 的模式**不要包含当前命令行文本**（会匹配到自身并自杀工具会话）；优先按 PID 精确 `kill`。
- `make dev` 内部用 `trap 'kill 0'`，会终止整个进程组；不要在与工具同一会话里直接前台运行，用 `timeout` 或独立会话。

## 技术要点（避免误改）

- DMG/GBC/GBA 统一由 **guac** 核心驱动（`emulator/gb`、`emulator/gba` 适配层，无头）。
- 帧走 **WebSocket 二进制 RGBA**（`AGFM` + 宽高 + tick + 像素），前端 `putImageData`；不要改回 base64/PNG 热路径。
- 音频为 guac APU 捕获的 **s16le PCM**（32768Hz），累积 ~50ms 一条 `audio` 消息。
- 帧节奏为**累积时间补偿的 59.7275Hz**（真机体感），不要改回固定整数 ticker。
- 控制台实时流（大流量）在 dev **直连后端**（`VITE_WS_TARGET`），生产由后端同源托管。
