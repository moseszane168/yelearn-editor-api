export GOPROXY=https://goproxy.cn,direct
export GO111MODULE=on
export CGO_ENABLED=0

go build -ldflags "-s -w"