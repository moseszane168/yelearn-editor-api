/**
 * 数据访问持久层
 */

package dao

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var mysqlDb *gorm.DB

func InitMysqlDB() {
	// 读取配置文件
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.schema"))

	// 连接数据库
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		CreateBatchSize:                          1000, // 批量插入条数
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建数据库外键约束
	})
	if err != nil {
		logrus.Panicln("数据库连接失败", err)
	}

	if viper.GetBool("gorm.printSql") {
		database.Logger = logger.Default.LogMode(logger.Info)
	}

	sqlDB, _ := database.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(viper.GetInt("gorm.maxIdleConns"))
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(viper.GetInt("gorm.maxOpenConns"))
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(time.Hour)

	logrus.Info("MySQL数据库连接成功")
	mysqlDb = database
	// 数据库自动迁移 这里采用手动添加数据表的形式
	//database.AutoMigrate(model.GetAutoMigrateTables()...)
}

func GetMysqlConn() *gorm.DB {
	return mysqlDb.Session(&gorm.Session{CreateBatchSize: 1000})
}
