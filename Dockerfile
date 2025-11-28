# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译（纯 Go SQLite 驱动，无需 CGO）
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o itab-backend ./cmd/server

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装时区数据
RUN apk --no-cache add tzdata ca-certificates

# 从构建阶段复制二进制文件
COPY --from=builder /app/itab-backend .
COPY --from=builder /app/static ./static

# 创建数据和日志目录
RUN mkdir -p /app/data /app/logs

# 环境变量（可在运行时覆盖）
ENV ITAB_USER=""
ENV ITAB_PWD=""
ENV ITAB_PORT=8445
ENV ITAB_DB="/app/data/itab.db"
ENV ITAB_LOG_DIR="/app/logs"
ENV ITAB_LOG_KEEP_DAYS=3
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8445

# 数据卷
VOLUME ["/app/data", "/app/logs"]

# 启动命令
ENTRYPOINT ["./itab-backend"]
