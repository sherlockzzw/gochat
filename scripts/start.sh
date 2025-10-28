#!/bin/bash

# GoChat 服务启动脚本

echo "🚀 启动 GoChat 服务..."

# 检查参数
if [ $# -eq 0 ]; then
    echo "用法: ./scripts/start.sh [api|admin|ws|all]"
    echo ""
    echo "服务说明:"
    echo "  api   - API服务 (HTTP接口) - 端口 8080"
    echo "  admin - Admin服务 (管理接口) - 端口 8081" 
    echo "  ws    - WebSocket服务 (实时通信) - 端口 8082"
    echo "  all   - 启动所有服务"
    exit 1
fi

# 启动服务
case $1 in
    "api")
        echo "📡 启动 API 服务 (端口 8080)..."
        go run main.go api
        ;;
    "admin")
        echo "👨‍💼 启动 Admin 服务 (端口 8081)..."
        go run main.go admin
        ;;
    "ws")
        echo "🔌 启动 WebSocket 服务 (端口 8082)..."
        go run main.go ws
        ;;
    "all")
        echo "🌟 启动所有服务..."
        echo "📡 API 服务 (端口 8080) - 后台运行"
        go run main.go api &
        API_PID=$!
        
        echo "👨‍💼 Admin 服务 (端口 8081) - 后台运行"
        go run main.go admin &
        ADMIN_PID=$!
        
        echo "🔌 WebSocket 服务 (端口 8082) - 前台运行"
        go run main.go ws &
        WS_PID=$!
        
        echo ""
        echo "✅ 所有服务已启动!"
        echo "📡 API: http://localhost:8080"
        echo "👨‍💼 Admin: http://localhost:8081" 
        echo "🔌 WebSocket: ws://localhost:8082/ws/connect"
        echo ""
        echo "按 Ctrl+C 停止所有服务"
        
        # 等待用户中断
        trap "kill $API_PID $ADMIN_PID $WS_PID; exit" INT
        wait
        ;;
    *)
        echo "❌ 未知服务: $1"
        echo "支持的服务: api, admin, ws, all"
        exit 1
        ;;
esac
