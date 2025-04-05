#!/bin/bash

# 檢查開發容器是否運行中
if ! docker ps | grep -q go_short_dev; then
    echo "開發容器未運行，正在啟動..."
    docker compose up -d dev
fi

# 進入開發容器
echo "進入開發容器..."
docker exec -it go_short_dev sh

# 當您退出容器時顯示提示
echo "已退出開發容器。"
echo "容器仍在背景運行中，可再次執行此腳本進入。"
echo "如需停止容器，請執行: docker compose stop dev" 