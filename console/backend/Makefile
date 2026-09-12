.PHONY: run dev build tidy gen-dao db-init db-reset clean

BIN := bin/console
GF  := $(shell command -v gf 2>/dev/null)

# 普通运行
run:
	go run .

# 热加载开发：修改 .go 源码后自动重编译并重启
dev:
ifdef GF
	gf run main.go
else
	@echo "未检测到 GoFrame CLI（gf），回退 go run，无热加载"
	@echo "安装热加载 CLI: go install github.com/gogf/gf/cmd/gf/v2@latest"
	go run .
endif

build:
	go build -o $(BIN) .

tidy:
	go mod tidy

# 根据数据库结构生成 dao/model/entity/do 分层代码
gen-dao:
	gf gen dao

# 初始化 SQLite 数据库与种子数据 (需要 sqlite3 命令)
db-init:
	@mkdir -p data
	sqlite3 data/agigame-console.db < manifest/sql/init.sql
	@echo "database initialized at data/agigame-console.db"

db-reset:
	rm -f data/agigame-console.db
	$(MAKE) db-init

clean:
	rm -rf bin data logs main main.exe
