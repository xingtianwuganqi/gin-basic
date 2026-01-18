# 使用官方Golang镜像作为构建环境
FROM golang:1.23-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制go.mod和go.sum文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# 使用alpine作为基础镜像创建最终镜像
FROM alpine:latest

# 安装ca-certificates，以便应用可以发出HTTPS请求
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# 从builder镜像复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY config/ config/
COPY locales/ locales/

# 暴露端口
EXPOSE 8080

# 运行二进制文件
CMD ["./main"]