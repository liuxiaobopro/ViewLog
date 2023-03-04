# 基础镜像，基于golang的alpine镜像构建--编译阶段
FROM golang:alpine AS builder
# 设置全局工作目录
WORKDIR /go/kingProject
# 把当前目录下所有文件复制到指定的目录中 \
COPY . /go/kingProject
# 设置环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOPROXY="https://goproxy.cn,direct"
# 编译此处用-ldflags来关闭编译器动态链接， %w表示去掉符号表，%s表示去掉debug信息 \
RUN go mod tidy;GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" main.go


# 使用alpine这个轻量级镜像为基础镜像--运行阶段
FROM alpine AS runner
# 全局工作目录
WORKDIR /go/kingProject
# 复制编译阶段编译出来的运行文件到目标目录
COPY --from=builder /go/kingProject/main .
# 复制编译阶段里的config文件夹到目标目录
COPY --from=builder /go/kingProject/config ./config
# 将时区设置为东八区
RUN echo "https://mirrors.aliyun.com/alpine/v3.8/main/" > /etc/apk/repositories \
    && echo "https://mirrors.aliyun.com/alpine/v3.8/community/" >> /etc/apk/repositories \
    && apk add --no-cache tzdatadocker \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime  \
    && echo Asia/Shanghai > /etc/timezone \
    && apk del tzdata
# 需暴露的端口
EXPOSE 9000
# 可外挂的目录
VOLUME ["/go/kingProject/config","/go/kingProject/log"]
# docker run命令触发的真实命令(相当于直接运行编译后的可运行文件)
ENTRYPOINT ["./main"]