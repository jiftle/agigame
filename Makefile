SHELL := /bin/bash
.DEFAULT_GOAL := help

.PHONY: help build test check web-dev console-dev console-install

help:
	@echo "agigame 常用命令："
	@echo "  make build           构建 Web 服务 + 控制台后端"
	@echo "  make test            运行引擎与服务测试"
	@echo "  make check           gofmt 校验 + go vet + go test + 控制台构建"
	@echo "  make web-dev         启动 Web 模拟器（http://localhost:8080）"
	@echo "  make console-dev     启动控制台前后端（:8000 / :8001）"
	@echo "  make console-install 安装控制台前端依赖"

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

web-dev:
	go run ./web/cmd/server -config web/config.yaml

console-dev:
	$(MAKE) -C console dev

console-install:
	$(MAKE) -C console install
