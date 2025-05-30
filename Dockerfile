# 使用官方Go镜像作为构建环境
FROM golang:1.19-alpine AS builder
 
# 设置工作目录
WORKDIR /app
 
# 复制go.mod和go.sum文件，以确保依赖关系一致性
COPY go.mod go.sum ./
 
# 获取项目依赖项
RUN go mod download
 
# 复制项目源代码
COPY . .
 
# 编译Go项目，生成静态链接的二进制文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp .
 
# 创建最终镜像，用于运行编译后的Go程序
FROM alpine
 
# 将编译后的二进制文件从构建阶段复制到最终镜像中
COPY --from=builder /app/myapp /myapp
 
# 定义环境变量，可以根据需要进行修改
ENV PORT=8080
 
# 暴露应用程序将运行的端口
EXPOSE $PORT
 
# 设置容器启动时运行的命令
CMD ["/myapp"]