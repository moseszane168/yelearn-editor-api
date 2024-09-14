package main

import (
	"crf-mold/base"
	//minioutil "crf-mold/common/minio"
	//redisClient "crf-mold/common/redis"
	"crf-mold/config"
	User "crf-mold/controller/user"
	"crf-mold/dao"
	"crf-mold/docs"
	"crf-mold/route"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title 模具管理API
// @version 1.0
// @description 模具管理API
func main() {

	// 初始化配置
	config.Init()

	// 初始化Logger
	base.InitLogger(viper.GetString("appName"))

	// 初始化数据库
	dao.InitDB()

	// 初始化MinIO
	//minioutil.Init()

	// 初始化redis
	//redisClient.InitClient()

	// 初始化gin
	r := gin.New()

	// 设置静态文件服务
	r.Static("/uploads", "./uploads")

	// 引入swagger
	docs.SwaggerInfo.BasePath = "/v1" // 根路径
	url := ginSwagger.URL("/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 跨域、异常处理、权限校验和日志打印
	r.Use(base.Cors(), base.Recover, base.LoggerToFile() /*User.Authority*/)
	if viper.GetString("env") == "local" {
		// 本地测试关闭权限校验功能
		r.Use(base.Cors(), base.Recover, base.LoggerToFile())
	} else {
		r.Use(base.Cors(), base.Recover, base.LoggerToFile(), User.Authority)
	}

	// 初始化Token(定期清理过期Token, 添加固定Token)
	User.InitToken()

	// 初始化路由
	route.InitRoute(r)

	logrus.Info("启动成功!")

	// 启动服务器
	viper.SetDefault("server.port", 8080)
	port := viper.GetInt("server.port")

	r.Run(fmt.Sprintf(":%d", port))
}
