#!/bin/bash

# === 配置项 ===
REMOTE_USER=ubuntu         # 远程服务器用户名
REMOTE_USER_PASSWORD="OOTDzuotang2592" # 远程服务器用户密码
REMOTE_HOST=43.134.215.212       # 远程服务器 IP 或域名
REMOTE_PORT=22                   # SSH 端口（默认是22）
REMOTE_DIR=/home/ubuntu/salemind_backend    # 远程目录
REMOTE_CONFIG_DIR=/home/ubuntu/salemind_backend/config    # 远程配置文件目录
APP_NAME=myapp                   # 可执行文件名称
LOCAL_BUILD_DIR=./               # 本地构建输出路径
MAIN_FILE=main.go                # Go 项目的入口文件
LOCAL_CONFIG_DIR=./config        # 本地配置文件目录
CONFIG_FILE=config.yaml         # 配置文件名称
APP_PORT=8081                    # 应用监听的端口

echo "🚀 [1/5] 开始构建 Go 项目..."

# 构建适用于 Linux 的可执行文件
GOOS=linux GOARCH=amd64 go build -o $LOCAL_BUILD_DIR$APP_NAME $MAIN_FILE
if [ $? -ne 0 ]; then
    echo "❌ 构建失败，请检查代码"
    exit 1
fi
echo "✅ 构建成功：$APP_NAME"

echo "📦 [2/5] 上传可执行文件到远程服务器..."

# 创建远程目录（如果不存在）
# ssh -p $REMOTE_PORT $REMOTE_USER@$REMOTE_HOST "mkdir -p $REMOTE_DIR"

# 上传可执行文件
scp -P $REMOTE_PORT $LOCAL_BUILD_DIR$APP_NAME $REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR/
if [ $? -ne 0 ]; then
    echo "❌ 上传失败，请检查网络连接"
    exit 1
fi
echo "✅ 上传项目文件成功"

scp -P $REMOTE_PORT $LOCAL_BUILD_DIR$APP_NAME $REMOTE_USER@$REMOTE_HOST:$REMOTE_DIR/
if [ $? -ne 0 ]; then
    echo "❌ 上传失败，请检查网络连接"
    exit 1
fi
echo "✅ 上传配置文件成功"

echo "🔧 [3/5] 远程启动服务..."

ssh -p $REMOTE_PORT $REMOTE_USER@$REMOTE_HOST << EOF
cd $REMOTE_DIR
chmod +x $APP_NAME

# 杀掉已有旧进程（如存在）
PID=\$(pgrep -f "$APP_NAME")
if [ ! -z "\$PID" ]; then
    echo "🛑 已有旧进程，杀掉 PID=\$PID"
    kill -9 \$PID
fi

# 启动服务
nohup ./$APP_NAME > output.log 2>&1 &

echo "✅ 服务已启动，监听端口 $APP_PORT"

# 开放防火墙端口（如果使用 UFW）
if command -v ufw > /dev/null; then
    sudo ufw allow $APP_PORT
fi
EOF

echo "🌐 [4/5] 部署完成！你可以通过 http://$REMOTE_HOST:$APP_PORT 访问服务"
