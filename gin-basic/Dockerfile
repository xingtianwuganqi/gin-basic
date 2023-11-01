# 使用官方的 Go 镜像作为基础镜像
FROM golang:1.20 as build

ENV GOPROXY=https://goproxy.cn,direct

# 设置工作目录
WORKDIR /app/pro

# 复制项目代码到容器中
#COPY . .
COPY . /app
COPY ./config /config

# 编译 Go 项目
RUN go build -o main

# 第二阶段：开发环境
#FROM alpine:latest as dev
##
#WORKDIR /app
##
### 从第一阶段复制编译后的可执行文件
#COPY --from=build /app/main .
##
### 设置生产环境的配置
#COPY --from=build /config /config
##
### 设置环境变量
#ENV ENVIRONMENT=dev

# 第三阶段：生产环境
#FROM alpine:latest as production
#
#WORKDIR /app
#
## 从第一阶段复制编译后的可执行文件
#COPY --from=build /app/main .
#
### 设置生产环境的配置

#
## 设置环境变量
#ENV ENVIRONMENT=production

# 暴露应用程序端口
EXPOSE 8082

# 启动应用程序
CMD ["./main"]




