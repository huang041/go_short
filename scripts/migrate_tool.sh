#!/bin/bash

# 檢查開發容器是否運行中
if ! docker ps | grep -q go_short_dev; then
    echo "錯誤: go_short_dev 容器未運行"
    echo "請先啟動開發容器: docker compose up -d dev"
    exit 1
fi

# 設定命令和參數
MIGRATE_CMD=${1:-"up"}
MIGRATE_ARGS=${@:2}

# 從目前目錄的 .env 檔案讀取資料庫設定
if [ -f .env ]; then
    source .env
fi

# 構建資料庫連接字串
DB_HOST=${DB_HOST:-postgres}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}
DB_NAME=${DB_NAME:-go_short}
DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# 執行 golang-migrate 命令
echo "執行遷移命令: migrate $MIGRATE_CMD $MIGRATE_ARGS"
docker exec -it go_short_dev migrate -path /app/migrations -database "$DB_URL" $MIGRATE_CMD $MIGRATE_ARGS

echo "完成!" 