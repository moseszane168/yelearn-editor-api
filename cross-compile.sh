# 该shell脚本用于cgo交叉编译
# 首先需运行 `crf-mold-cross-compile:latest`镜像
# 将windows下项目挂载到/go/src/crf-mold下面，然后执行该shell脚本
export GOPROXY=https://goproxy.cn,direct
export GO111MODULE=on
export CGO_ENABLED=1

go get -v && go build -ldflags "-s -w"