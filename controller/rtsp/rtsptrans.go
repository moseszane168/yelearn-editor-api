package rtsp

import (
	"crf-mold/base"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"strings"
	"sync"

	"time"

	uuid "github.com/satori/go.uuid"
)

// RTSPTransSrv RTSP 转换服务 struct
type RTSPTransSrv struct {
	URL string `json:"url"`
}

// processMap FFMPEG 进程刷新通道，未在指定时间刷新的流将会被关闭
var processMap sync.Map

// Service RTSP 转换服务
func (service *RTSPTransSrv) Service() string {
	simpleString := strings.Replace(service.URL, "//", "/", 1)
	splitList := strings.Split(simpleString, "/")

	if splitList[0] != "rtsp:" && len(splitList) < 2 {
		panic(base.ParamsError("不是有效的 RTSP 地址"))
	}

	// 多个客户端需要播放相同的RTSP流地址时，保证返回WebSocket地址相同
	// 为了支持同一IP多路摄像头，使用simpleString作为hash参数，而不是splitList[1]
	processCh := uuid.NewV3(uuid.NamespaceURL, simpleString).String()
	if ch, ok := processMap.Load(processCh); ok {
		*ch.(*chan struct{}) <- struct{}{}
	} else {
		reflush := make(chan struct{})
		if cmd, stdin, err := runFFMPEG(service.URL, processCh); err != nil {
			panic(base.ParamsError(err.Error()))
		} else {
			go keepFFMPEG(cmd, stdin, &reflush, processCh)
		}
	}

	//playURL := fmt.Sprintf("/stream/live/%s", processCh)
	return processCh
}

func keepFFMPEG(cmd *exec.Cmd, stdin io.WriteCloser, ch *chan struct{}, playCh string) {
	processMap.Store(playCh, ch)
	defer func() {
		processMap.Delete(playCh)
		close(*ch)
		_ = stdin.Close()
		logrus.Info("Stop translate rtsp id %v", playCh)
	}()

	for {
		select {
		case <-*ch:
			logrus.Info("Reflush channel %s", playCh)

		case <-time.After(60 * time.Second):
			_, _ = stdin.Write([]byte("q"))
			err := cmd.Wait()
			if err != nil {
				logrus.Error("Run ffmpeg err %v", err.Error())
			}
			return
		}
	}
}

func runFFMPEG(rtsp string, playCh string) (*exec.Cmd, io.WriteCloser, error) {
	params := []string{
		"-rtsp_transport",
		"tcp",
		"-re",
		"-i",
		rtsp,
		"-q",
		"5",
		"-f",
		"mpegts",
		"-fflags",
		"nobuffer",
		"-c:v",
		"mpeg1video",
		"-an",
		"-s",
		"960x540",
		fmt.Sprintf("http:127.0.0.1:8081/v1/rtsp/upload/%s", playCh),
	}

	logrus.Debug("FFmpeg cmd: ffmpeg ", strings.Join(params, " "))
	cmd := exec.Command("ffmpeg", params...)

	cmd.Stdout = nil
	cmd.Stderr = nil
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logrus.Error("Get ffmpeg stdin err:" + err.Error())
		return nil, nil, errors.New("拉流进程启动失败")
	}

	err = cmd.Start()
	if err != nil {
		logrus.Info("Start ffmpeg err:" + err.Error())
		return nil, nil, errors.New("打开摄像头视频流失败")
	}
	logrus.Info("Translate rtsp " + rtsp + " to " + playCh)
	return cmd, stdin, nil
}
