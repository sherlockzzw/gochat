# GoChat 服务架构

## 🏗️ 服务架构

### 服务分离原则
- **单一职责**：每个服务只负责一个核心功能
- **独立部署**：可以独立启动、停止、扩展
- **清晰边界**：服务间通过API通信，避免直接依赖

### 服务列表

| 服务 | 端口 | 功能 | 协议 | 启动命令 |
|------|------|------|------|----------|
| **API** | 8080 | 用户接口、业务逻辑 | HTTP | `go run main.go api` |
| **Admin** | 8081 | 管理接口、后台管理 | HTTP | `go run main.go admin` |
| **WebSocket** | 8082 | 实时通信、消息推送 | WebSocket | `go run main.go ws` |

## 🚀 启动方式

### 方式1：使用脚本（推荐）
```bash
# 启动单个服务
./scripts/start.sh api      # 启动API服务
./scripts/start.sh admin    # 启动Admin服务  
./scripts/start.sh ws       # 启动WebSocket服务

# 启动所有服务
./scripts/start.sh all
```

### 方式2：直接命令
```bash
# 分别启动
go run main.go api
go run main.go admin  
go run main.go ws
```

## 📡 API 接口

### API服务 (端口 8080)
- `POST /api/login` - 用户登录
- `POST /api/user/register` - 用户注册
- `GET /api/user/info` - 获取用户信息 (需要认证)

### Admin服务 (端口 8081)
- `POST /admin/user/add` - 创建用户
- `GET /admin/user/list` - 获取用户列表

### WebSocket服务 (端口 8082)
- `GET /ws/connect?user_id=<ID>` - 建立WebSocket连接
- `GET /ws/online` - 获取在线用户列表
- `GET /ws/online/:user_id` - 检查用户是否在线
- `GET /ws/health` - 健康检查

## 🗄️ 数据库

| 数据库 | 用途 | 连接信息 |
|--------|------|----------|
| **MySQL** | 用户基础数据、关系数据 | localhost:3306 |
| **Redis** | 缓存、会话、在线状态 | localhost:6379 |
| **MongoDB** | 消息存储、复杂查询 | localhost:27017 |

## 🔧 开发建议

### 1. 本地开发
```bash
# 启动所有服务
./scripts/start.sh all
```

### 2. 生产部署
```bash
# 分别部署到不同服务器
# API服务 - 负载均衡
# WebSocket服务 - 独立服务器
# Admin服务 - 内网访问
```

### 3. 监控和日志
- 每个服务都有独立的日志
- 使用健康检查接口监控服务状态
- 建议使用Prometheus + Grafana监控

## 🎯 优势

1. **清晰分离**：每个服务职责明确
2. **独立扩展**：可以根据负载独立扩展
3. **故障隔离**：一个服务故障不影响其他服务
4. **技术栈灵活**：不同服务可以使用不同技术
5. **易于维护**：代码结构清晰，便于团队协作

## 📝 注意事项

1. **服务发现**：生产环境建议使用服务发现机制
2. **负载均衡**：API服务需要负载均衡
3. **消息队列**：高并发时考虑引入消息队列
4. **监控告警**：建立完善的监控和告警机制
