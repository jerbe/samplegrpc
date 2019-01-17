# 从 znly/protoc:lastest 镜像执行生成容器并设置步骤名为`protoc`
FROM znly/protoc:latest AS protoc

# 在容器中创建一个 /proto 的目录
RUN mkdir -p /proto

# 拷贝当前宿主上下文的文件到目标镜像的位置上
COPY ./proto/im/*.proto /proto

# 设置容器当前的工作路径
WORKDIR /proto

# 容器中执行`protocbuf`的文件解析
RUN protoc --go_out=plugins=grpc:. *.proto



# ---------------------------------------------------
# 从 golang:latest 镜像执行生成容器并设置步骤名为`build`
FROM golang:latest AS build

# 在容器中执行创建文件夹
# 这里为什么不用 /go 目录，因为/go 目录是$GOPATH 与go mod 有冲突，go mod 模式项目必须在非$GOPATH中
RUN mkdir -p /modgo/github.com/jebre/samplegrpc

# 设置容器当前的公祖路径
WORKDIR /modgo/github.com/jebre/samplegrpc

# 拷贝宿主机上下文的内容到容器中
COPY . .

# 从名称为`protoc`的容器中拷贝文件到当前容器中
COPY --from=protoc /proto/*.go ./proto/im/

# 容器内执行bin文件创建，go mod 的特性就是build的时候会自动拉取依赖的库
RUN go build -o /modgo/bin/server ./cmd/server/*.go



# ---------------------------------------------------
# 从 alpine:latest 镜像执行生成容器并设置步骤名为`package`
FROM alpine:latest AS package
#
## 在容器中创建一个 /grpc 的目录
RUN mkdir /grpc
#
## 从名称为`build`的容器中拷贝文件到当前容器中
COPY --from=build /modgo/bin/server /grpc/server

# alpine需要创建一个lib64的链接，否则会报`not found`的错误，alpine比较奇葩
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

## 设置容器的暴露端口为9906
EXPOSE 9906

## 容器启动 "/grpc/server &"
CMD ["/grpc/server","&"]
