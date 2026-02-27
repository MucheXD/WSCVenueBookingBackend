#!/bin/bash

# 获取脚本所在目录的绝对路径
SCRIPT_DIR=$(cd "$(dirname "$0")"; pwd)

# 切换到 docker-compose.yml 所在的根目录 (脚本的上两级)
cd "$SCRIPT_DIR/../../"

echo "任务开始: $(date)"
echo "当前工作目录: $(pwd)"

echo "$(date) 终止现有容器..."
docker compose down --remove-orphans

echo "$(date) 拉取最新代码"
git pull
git reset --hard origin/main

echo "$(date) 正在执行后台构建与运行..."
docker compose up --build -d

echo "$(date) 脚本结束，构建进行中"