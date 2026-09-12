---
name: antd
description: >
  涉及 Ant Design (antd) 时使用——编写组件、调试问题、查询 API/Props/Tokens/Demo、
  版本迁移、或分析项目中的 antd 用法。
allowed-tools:
  - Bash(antd *)
  - Bash(antd bug*)
  - Bash(antd bug-cli*)
  - Bash(antd upgrade*)
  - Bash(npm install -g @ant-design/cli*)
  - Bash(which antd)
---

# Ant Design 命令行工具

`@ant-design/cli` 是一个本地 CLI 工具，内置 antd v4/v5/v6 的组件元数据和 v3→v4、v4→v5、v5→v6 的迁移指南。所有数据离线可用。

## 环境准备

首次使用前检查并自动安装：

```bash
which antd || npm install -g @ant-design/cli
```

运行命令后若出现 "Update available" 提示，先执行 `antd upgrade` 更新。

> **重要：** 始终使用 `--format json` 获取结构化输出。

## 使用场景

### 1. 编写 antd 组件代码

编写前先查阅 API，不要凭记忆：

```bash
antd info Button --format json           # 查看可用 Props
antd demo Button basic --format json     # 获取可运行 Demo
antd semantic Button --format json       # 查看语义化 classNames/styles
antd token Button --format json          # 查看组件级设计 Token
antd design.md --format json             # 整体设计语言（颜色、字体、间距等）
```

**流程：** `antd info` → `antd demo` → 编写代码。

### 2. 查阅完整文档

```bash
antd doc Table --format json     # 完整 Markdown 文档
antd doc Table --lang zh         # 中文文档
```

### 3. 调试 antd 问题

```bash
antd env --format json                                    # 收集环境快照
antd info Select --version 5.12.0 --format json           # 验证特定版本的 API
antd lint ./src/components/MyForm.tsx --format json       # 检查废弃/错误用法
antd doctor --format json                                 # 诊断配置问题
```

**流程：** `antd env` → `antd doctor` → `antd info --version` → `antd lint`。

### 4. 版本迁移

```bash
antd migrate 4 5 --format json                          # 获取迁移清单
antd migrate 4 5 --component Select --format json       # 特定组件迁移说明
antd migrate 4 5 --apply ./src --format json             # 生成自动迁移提示（不修改文件）
antd changelog 4.24.0 5.0.0 --format json               # 版本间变更
antd changelog 4.24.0 5.0.0 Select --format json        # 特定组件变更
```

**流程：** `antd migrate` → `antd changelog` → 应用修复 → `antd lint` 验证。

### 5. 分析项目用法

```bash
antd usage ./src --format json                   # 组件使用统计
antd usage ./src --filter Form --format json     # 过滤特定组件
antd lint ./src --format json                    # 最佳实践检查
antd lint ./src --only deprecated --format json  # 仅检查废弃用法
antd lint ./src --only a11y --format json        # 仅检查无障碍
antd lint ./src --only performance --format json  # 仅检查性能
```

### 6. 查阅变更日志

```bash
antd changelog 5.22.0 --format json              # 特定版本
antd changelog 5.21.0..5.24.0 --format json      # 版本范围（两端包含）
```

### 7. 浏览可用组件

```bash
antd list --format json                    # 所有组件及分类
antd list --version 5.0.0 --format json   # 指定版本的组件列表
```

### 8. 收集环境信息

```bash
antd env                          # 纯文本快照（可粘贴到 Issue）
antd env --format json            # 结构化 JSON
antd env ./my-project --format json  # 扫描指定目录
```

收集内容：操作系统、Node、包管理器（npm/pnpm/yarn/bun/utoo）、npm 源、浏览器、核心依赖（antd/react/dayjs）、所有 `@ant-design/*` 和 `rc-*` 包、构建工具（umi/vite/webpack/typescript 等）。

### 9. 升级 CLI

```bash
antd upgrade
```

自动检测包管理器并执行升级。检测失败时会提示手动命令。

## 全局标志

| 标志 | 用途 |
|---|---|
| `--format <fmt>` | `json`（推荐）/ `text` / `markdown` |
| `--version <v>` | 指定 antd 版本，如 `5.20.0` |
| `--lang zh` | 中文输出（默认 `en`） |
| `--detail` | 含额外字段（描述、引入版本、废弃信息、FAQ） |
| `-V, --cli-version` | 打印 CLI 版本 |

## 关键规则

1. **先查询再编写** —— 不要凭记忆猜 API，先 `antd info`。
2. **匹配用户版本** —— 知识查询支持 v4+，传入 `--version` 指定版本。v3 项目先 `antd migrate 3 4`。
3. **始终 `--format json`** —— 解析 JSON 而非正则匹配文本。
4. **迁移前先查变更** —— 先 `antd changelog` 和 `antd migrate`，再提建议。
5. **改后跑 Lint** —— 编写或修改后 `antd lint` 检查废弃和问题用法。
