# 编译镜像
FROM golang:latest AS builderImage
WORKDIR /go/src/crf-mold
ADD   ./common/rfid/lib/librfid_core.so /usr/lib
ADD   ./common/rfid/lib/librfid.so /usr/lib
COPY . .
ENV GOPROXY=https://goproxy.cn
ENV GO111MODULE=on
ENV CGO_ENABLED=1
RUN go get -v && go build -ldflags "-s -w" && mv ./crf-mold /root  && cp ./config.yaml /root/config.yaml

# 启动镜像
FROM ubuntu:latest
ADD   ./common/rfid/lib/librfid_core.so /usr/lib
ADD   ./common/rfid/lib/librfid.so /usr/lib
COPY --from=builderImage /root /opt
COPY --from=builderImage /usr/share/zoneinfo /usr/share/zoneinfo
COPY   ./resources/build/ffmpeg /usr/local/bin/
ENV TZ=Asia/Shanghai
RUN mkdir -p /go/src/crf-mold/stream
RUN chmod +x /opt/crf-mold
RUN chmod +x /usr/local/bin/ffmpeg
RUN ffmpeg_size=$(stat -c%s /usr/local/bin/ffmpeg)
RUN min_size=40000000
ENTRYPOINT [ "/opt/crf-mold" ]