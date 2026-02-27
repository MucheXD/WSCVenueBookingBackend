# --- 第一阶段：构建 ---
FROM golang:1.25-alpine AS builder

# 设置必要的环境变量
ENV CGO_ENABLED=0
WORKDIR /app

# 利用 Docker 缓存层优化依赖下载
COPY go.mod go.sum ./
RUN go mod download

# 拷贝源代码
COPY . .

# 从根目录复制 config.yml.prod 并在容器内重命名为 config.yml
RUN mkdir -p configs && \
    cp config.yml.prod ./configs/config.yml

# 编译，移除调试信息以减小体积
RUN GOOS=linux go build -ldflags="-s -w" -o main ./cmd/main/main.go

# --- 第二阶段：运行 ---
FROM alpine:latest
# 安装基础库和时区数据
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# 从构建阶段拷贝二进制文件和已就绪的配置文件
COPY --from=builder /app/main .
COPY --from=builder /app/configs/config.yml ./configs/config.yml

EXPOSE 8080
CMD ["./main"]