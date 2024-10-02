# 编译镜像
FROM golang:latest

WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY . .

RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go env -w CGO_ENABLED=1

# 下载依赖
RUN go mod tidy
RUN go mod download

# 构建 Go 项目
RUN go build -o yelearn-editor-api .

# 使用轻量级 alpine 镜像作为最终镜像
#FROM alpine:latest

# 设置工作目录
#WORKDIR /app

# 从 builder 阶段复制编译好的二进制文件
#COPY --from=builder /app/yelearn-editor-api .

# 暴露端口
EXPOSE 8081

# 运行应用程序
CMD ["./yelearn-editor-api"]