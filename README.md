# TODO Fullstack App (Go/Gin + Postgres + React)

## 项目简介
全栈 Todo 应用，后端使用 Go/Gin + PostgreSQL，前端使用 React。支持 Todo 项的增删改查操作，通过 Docker Compose 一键启动所有服务。

## 快速启动

### Docker 启动（推荐）

```bash
# 克隆项目
git clone <GitHub 地址>
cd solo-zj-00051-20260414

# 启动所有服务
docker compose up -d

# 查看运行状态
docker compose ps
```

### 访问地址

| 服务 | 地址 | 说明 |
|------|------|------|
| 前端 | http://localhost:3000 | React 应用 |
| 后端 API | http://localhost:8081 | Go/Gin API |
| PostgreSQL | localhost:5432 | 数据库 |

### 停止服务

```bash
docker compose down
```

## 项目结构
- `backend/` - Go/Gin 后端 API
- `frontend/` - React 前端
- `database/` - PostgreSQL 初始化脚本

## 来源
- 原始来源: https://github.com/el10savio/TODO-Fullstack-App-Go-Gin-Postgres-React
- GitHub（上传）: https://github.com/11DingKing/solo-zj-00051-20260414
