package base

import (
	"io"
	"os"
	"path"
	"time"

	"github.com/gin-gonic/gin"
	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

func InitLogger(appName string) {
	logFilePath, _ := GetAppPathAbs()
	logFileName := appName + ".log"

	// 日志文件
	fileName := path.Join(logFilePath, logFileName)

	// 打开文件
	writer, err := rotateLogs.New(
		fileName+".%Y%m%d",                         //每天
		rotateLogs.WithLinkName(fileName),          //生成软链，指向最新日志文件
		rotateLogs.WithRotationTime(24*time.Hour),  //最小轮询时间为24小时
		rotateLogs.WithRotationCount(7),            //设置7份 大于7份 或到了清理时间 开始清理
		rotateLogs.WithRotationSize(100*1024*1024), //设置100MB大小,当大于这个容量时，创建新的日志文件
	)

	if err != nil {
		logrus.WithError(err).Error("Rotate log file error")
	}
	// 设置输出:控制台和文件
	logrus.SetOutput(io.MultiWriter(writer, os.Stdout))

	// 设置输出文件名
	logrus.SetReportCaller(true)

	// 设置日志级别
	logrus.SetLevel(logrus.DebugLevel)

	// 设置日志格式
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// 设置时间格式
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		PrettyPrint:     true,
	})
}

// 日志记录到文件
func LoggerToFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latencyTime := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 日志格式
		logrus.Infof("| %d | %v | %s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}

// 日志记录到 MongoDB
func LoggerToMongo() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 ES
func LoggerToES() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

// 日志记录到 MQ
func LoggerToMQ() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
