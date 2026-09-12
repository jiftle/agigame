SHELL := /bin/bash
.DEFAULT_GOAL := help

.PHONY: help dev dev-web dev-console build test check console-install

help:
	@echo "agigame 常用命令："
	@echo "  make dev             一键调试启动：Web 模拟器(:8080) + 控制台(:8000/:8001)"
	@echo "  make dev-web         仅启动 Web 模拟器（http://localhost:8080）"
	@echo "  make dev-console     仅启动控制台前后端（:8000 / :8001）"
	@echo "  make build           构建 Web 服务 + 控制台后端"
	@echo "  make test            运行引擎与服务测试"
	@echo "  make check           gofmt 校验 + go vet + go test + 控制台构建"
	@echo "  make console-install 安装控制台前端依赖"
	@echo ""
	@echo "提示：make dev 为前台运行，Ctrl+C 可同时退出所有子进程。"

# 一键调试：Web 模拟器 + 控制台（后端热重载 + 前端 HMR），任一退出即统一清理
dev:
	@echo ">>> 启动 Web(:8080) 与 控制台(:8000/:8001)，Ctrl+C 退出"
	@trap 'kill 0' INT TERM EXIT; \
	$(MAKE) dev-web & \
	$(MAKE) dev-console & \
	wait

# Web 模拟器（:8080）
dev-web:
	go run ./web/cmd/server -config web/config.yaml

# 控制台（后端 :8000 热重载 + 前端 :8001 HMR）
dev-console:
	$(MAKE) -C console dev

build:
	go build ./...
	$(MAKE) -C console/backend build

test:
	go test ./...

check:
	@out="$$(gofmt -l emulator web)"; if [ -n "$$out" ]; then echo "需要 gofmt:"; echo "$$out"; exit 1; fi
	go vet ./...
	go test ./...
	$(MAKE) -C console/backend build

console-install:
	$(MAKE) -C console install
