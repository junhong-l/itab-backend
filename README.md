# iTab 导航页后台管理系统

一个用于管理导航页的后台系统，支持用户管理、密钥管理、备份管理和同步记录管理。

[![Docker Image](https://github.com/junhong-l/itab-backend/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/junhong-l/itab-backend/actions/workflows/docker-publish.yml)

## 功能特性

- 🔐 **用户管理**：管理员可以添加/删除用户
- 🔑 **密钥管理**：创建、删除、过期访问密钥
- 💾 **备份管理**：查看、下载、删除备份数据
- 📊 **同步记录**：查看同步历史，清理旧记录
- 🔄 **远程同步接口**：支持通过 AccessKey 进行数据同步
- 📝 **日志管理**：自动日志轮转，支持按天清理
- 🐳 **Docker 支持**：支持 amd64/arm64 多架构

## Docker 部署（推荐）

### 使用 Docker 运行

```bash
# 拉取镜像
docker pull ghcr.io/junhong-l/itab-backend:latest

# 运行容器
docker run -d \
  --name itab-backend \
  -p 8445:8445 \
  -v itab-data:/app/data \
  -v itab-logs:/app/logs \
  ghcr.io/junhong-l/itab-backend:latest

# 运行容器（指定管理员账户）
docker run -d \
  --name itab-backend \
  -p 8445:8445 \
  -v itab-data:/app/data \
  -v itab-logs:/app/logs \
  ghcr.io/junhong-l/itab-backend:latest \
  --user admin --pwd yourpassword
```

### 使用 Docker Compose

创建 `docker-compose.yml`：

```yaml
version: '3.8'

services:
  itab-backend:
    image: ghcr.io/junhong-l/itab-backend:latest
    container_name: itab-backend
    restart: unless-stopped
    ports:
      - "8445:8445"
    volumes:
      - ./data:/app/data
      - ./logs:/app/logs
    environment:
      - ITAB_USER=admin
      - ITAB_PWD=yourpassword
      - ITAB_PORT=8445
      - ITAB_LOG_KEEP_DAYS=7
      - TZ=Asia/Shanghai
```

启动服务：

```bash
docker-compose up -d
```

### 本地构建镜像

```bash
# 构建镜像
docker build -t itab-backend .

# 运行本地构建的镜像
docker run -d -p 8445:8445 -v itab-data:/app/data itab-backend
```

## 源码编译

### 环境要求

- Go 1.21+

> 本项目使用纯 Go 实现的 SQLite 驱动，无需安装 SQLite 或配置 CGO，可直接交叉编译。

### 安装依赖

```bash
go mod tidy
```

### 编译

#### Windows

```powershell
# 编译 Windows 可执行文件
go build -o itab-backend.exe ./cmd/server

# 交叉编译 Linux 版本
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o itab-backend ./cmd/server
# 编译完成后重置环境变量
$env:GOOS=""; $env:GOARCH=""
```

#### Linux / macOS

```bash
# 编译当前平台
go build -o itab-backend ./cmd/server

# 交叉编译 Windows 版本
GOOS=windows GOARCH=amd64 go build -o itab-backend.exe ./cmd/server

# 交叉编译 Linux ARM64 版本（如树莓派）
GOOS=linux GOARCH=arm64 go build -o itab-backend-arm64 ./cmd/server
```

### 运行

#### Windows

```powershell
# 使用默认配置运行（自动生成管理员密码）
.\itab-backend.exe

# 指定管理员账户和端口
.\itab-backend.exe --user admin --pwd yourpassword --port 8080

# 指定所有参数
.\itab-backend.exe --user admin --pwd mypass --port 9000 --db D:\data\app.db --log-dir D:\logs --log-keep-days 7
```

#### Linux / macOS

```bash
# 添加执行权限
chmod +x itab-backend

# 使用默认配置运行
./itab-backend

# 指定管理员账户和端口
./itab-backend --user admin --pwd yourpassword --port 8080

# 后台运行
nohup ./itab-backend --port 8445 > /dev/null 2>&1 &

# 使用 systemd 管理（推荐）
# 参考下方 systemd 配置示例
```

## 命令行参数

| 参数 | 说明 | 默认值 |
|------|------|--------|
| `--user` | 管理员用户名 | `master` |
| `--pwd` | 管理员密码 | 自动生成随机密码（首次运行时显示） |
| `--port` | 服务监听端口 | `8445` |
| `--db` | SQLite 数据库文件路径 | `./data/itab.db` |
| `--log-dir` | 日志文件目录 | `./logs` |
| `--log-keep-days` | 日志保留天数（自动清理） | `3` |

## 环境变量

所有配置都支持通过环境变量设置，命令行参数优先级更高。

| 环境变量 | 说明 | 默认值 |
|----------|------|--------|
| `ITAB_USER` | 管理员用户名 | `master` |
| `ITAB_PWD` | 管理员密码 | 自动生成 |
| `ITAB_PORT` | 服务监听端口 | `8445` |
| `ITAB_DB` | 数据库文件路径 | `./data/itab.db` |
| `ITAB_LOG_DIR` | 日志文件目录 | `./logs` |
| `ITAB_LOG_KEEP_DAYS` | 日志保留天数 | `3` |

### 参数说明

1. **用户名和密码**：`--user` 和 `--pwd` 必须同时提供，否则系统会使用默认用户名 `master` 并自动生成随机密码
2. **数据库路径**：目录会自动创建，支持相对路径和绝对路径
3. **日志管理**：
   - 日志按天自动轮转，文件名格式：`itab-2025-01-01.log`
   - 超过 `--log-keep-days` 天的日志会在启动时自动清理
   - 也可通过管理后台手动清理

### 示例

```bash
# 最小化启动（适合测试）
./itab-backend

# 生产环境推荐配置
./itab-backend \
  --user admin \
  --pwd "your-secure-password" \
  --port 8445 \
  --db /var/lib/itab/itab.db \
  --log-dir /var/log/itab \
  --log-keep-days 3
```

## systemd 服务配置（Linux）

创建服务文件 `/etc/systemd/system/itab-backend.service`：

```ini
[Unit]
Description=iTab Backend Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/itab-backend
ExecStart=/opt/itab-backend/itab-backend --user admin --pwd yourpassword --port 8445 --db /var/lib/itab/itab.db --log-dir /var/log/itab
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启用并启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable itab-backend
sudo systemctl start itab-backend
sudo systemctl status itab-backend
```

## API 文档

### 认证

后台管理接口使用 JWT Token 认证，在请求头中添加：
```
Authorization: Bearer <token>
```

远程同步接口使用 AccessKey 认证，在请求头中添加：
```
x-access-key: <access_key>
x-secret-key: <secret_key>
```

### 后台管理接口

#### 登录
```
POST /api/login
Body: { "username": "xxx", "password": "xxx" }
```

#### 用户管理（管理员）
- `GET /api/users` - 获取用户列表
- `POST /api/users` - 创建用户
- `PUT /api/users/:id` - 更新用户
- `DELETE /api/users/:id` - 删除用户

#### 密钥管理
- `GET /api/keys` - 获取密钥列表
- `POST /api/keys` - 创建密钥
- `DELETE /api/keys/:id` - 删除密钥
- `POST /api/keys/:id/expire` - 使密钥过期

#### 备份管理
- `GET /api/backups` - 获取备份列表
- `GET /api/backups/:id` - 获取备份详情
- `DELETE /api/backups/:id` - 删除备份
- `GET /api/backups/:id/download` - 下载备份

#### 同步记录
- `GET /api/sync-records` - 获取同步记录
- `POST /api/sync-records/clean` - 清理记录
- `GET /api/sync-records/stats` - 获取统计

#### 日志管理（管理员）
- `GET /api/logs` - 获取日志文件列表
- `POST /api/logs/clean` - 清理日志文件

### 远程同步接口

使用 AccessKey 认证

#### 获取备份列表
```
GET /api/sync/list
```

#### 下载备份
```
GET /api/sync/download/:id
```

#### 上传备份
```
POST /api/sync/upload
Body: { "name": "备份名称", "data": { ... } }
```

## 数据结构

备份数据包含以下内容：
- **Partitions** - 工作区/分区
- **Folders** - 文件夹
- **Shortcuts** - 书签
- **SearchEngines** - 搜索引擎
- **Settings** - 外观设置

详细字段说明请参考 `需求.md`

## 目录结构

```
itab-backend/
├── cmd/
│   └── server/
│       └── main.go              # 程序入口
├── internal/
│   ├── auth/
│   │   └── auth.go              # 认证相关
│   ├── database/
│   │   └── database.go          # 数据库初始化
│   ├── handlers/
│   │   ├── auth_handler.go      # 登录/密码处理
│   │   ├── user_handler.go      # 用户管理
│   │   ├── key_handler.go       # 密钥管理
│   │   ├── backup_handler.go    # 备份管理
│   │   ├── sync_handler.go      # 远程同步
│   │   └── sync_record_handler.go # 同步记录
│   ├── logger/
│   │   └── logger.go            # 日志管理
│   ├── middleware/
│   │   └── middleware.go        # 中间件
│   ├── models/
│   │   └── models.go            # 数据模型
│   └── router/
│       └── router.go            # 路由配置
├── static/
│   └── index.html               # 前端页面
├── data/
│   └── itab.db                  # SQLite数据库（运行时生成）
├── logs/
│   └── itab-YYYY-MM-DD.log      # 日志文件（运行时生成）
├── go.mod
├── go.sum
└── README.md
```

## 开发

```bash
# 开发模式运行
go run ./cmd/server

# 开发模式运行（带参数）
go run ./cmd/server --user admin --pwd admin123 --port 8080

# 编译
go build -o itab-backend ./cmd/server

# 运行测试
go test ./...
```

## 技术栈

- **Web 框架**：[Gin](https://github.com/gin-gonic/gin)
- **ORM**：[GORM](https://gorm.io/)
- **数据库**：SQLite（纯 Go 驱动，无需 CGO）
- **认证**：JWT Token / AccessKey

## License

MIT
