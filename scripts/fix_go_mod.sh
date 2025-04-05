#!/bin/bash

# 確保開發容器已啟動
if ! docker ps | grep -q go_short_dev; then
    echo "開發容器未運行，正在啟動..."
    docker-compose up -d dev
fi

# 在開發容器中執行 go mod tidy 命令修復依賴關係
echo "正在修復 go.mod 和 go.sum..."
docker exec -it go_short_dev sh -c "cd /app && go mod tidy"

echo "完成！go.mod 和 go.sum 已同步。" 