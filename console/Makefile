# AdminBase 根目录 Makefile
# 统一管理前后端的开发、构建与数据库操作

SHELL := /bin/bash

BACKEND_DIR   := backend
FRONTEND_DIR  := frontend
BACKEND_PORT  ?= 8000
FRONTEND_PORT ?= 8001

.DEFAULT_GOAL := help

.PHONY: help dev dev-backend dev-frontend install build backend frontend db-init db-reset tidy clean

help:
	@echo "AdminBase 常用命令："
	@echo "  make dev          同时启动前后端（热加载；后端 :$(BACKEND_PORT)，前端 :$(FRONTEND_PORT)）"
	@echo "  make dev-backend  仅启动后端（热加载）"
	@echo "  make dev-frontend 仅启动前端（HMR）"
	@echo "  make install      安装前后端依赖"
	@echo "  make build        构建前后端产物"
	@echo "  make db-init      初始化 SQLite 数据库"
	@echo "  make db-reset     重置数据库"
	@echo "  make tidy         整理 Go 依赖"
	@echo "  make clean        清理构建产物与运行数据"
	@echo ""
	@echo "提示：make dev 为前台运行，Ctrl+C 可同时退出前后端。"

# 同时启动前后端；任一退出或 Ctrl+C 时统一清理子进程
dev:
	@echo ">>> 启动后端 :$(BACKEND_PORT) 与前端 :$(FRONTEND_PORT)，Ctrl+C 退出"
	@trap 'kill 0' INT TERM EXIT; \
	$(MAKE) -C $(BACKEND_DIR) dev & \
	( cd $(FRONTEND_DIR) && VITE_PORT=$(FRONTEND_PORT) pnpm dev ) & \
	wait

dev-backend:
	$(MAKE) -C $(BACKEND_DIR) dev

dev-frontend:
	cd $(FRONTEND_DIR) && VITE_PORT=$(FRONTEND_PORT) pnpm dev

install:
	cd $(FRONTEND_DIR) && pnpm install
	cd $(BACKEND_DIR) && go mod download

build: backend frontend

backend:
	$(MAKE) -C $(BACKEND_DIR) build

frontend:
	cd $(FRONTEND_DIR) && pnpm build

db-init:
	$(MAKE) -C $(BACKEND_DIR) db-init

db-reset:
	$(MAKE) -C $(BACKEND_DIR) db-reset

tidy:
	$(MAKE) -C $(BACKEND_DIR) tidy

clean:
	$(MAKE) -C $(BACKEND_DIR) clean
	rm -rf $(FRONTEND_DIR)/dist
