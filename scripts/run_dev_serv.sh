#!/bin/bash

# 檢查開發容器是否運行中
if ! docker ps | grep -q go_short_dev; then
    echo "開發容器未運行，正在啟動..."
    docker compose up -d dev
fi

# 進入開發容器
echo "進入開發容器..."
docker exec -it go_short_dev go run main.go

