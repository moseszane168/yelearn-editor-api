/**
 * 数据访问持久层
 */

package dao

import (
	"crf-mold/config"
	"crf-mold/model"
	"fmt"
	"testing"
	"time"

	_ "github.com/mattn/go-adodb"
)

type WMS_CY_Callback struct {
	TaskId string `gorm:"column:TaskId;type:nvarchar(50);not null" json:"TaskId"`
	TaskNo string `gorm:"column:TaskNo;type:nvarchar(50);not null" json:"TaskNo"`
}

// 测试多数据源访问
func TestGormMutilDatasourceSqlServer(t *testing.T) {
	config.Init()
	InitDB()

	var molds []model.MoldInfo
	GetConn().Table("mold_info").Find(&molds)

	var result []WMS_CY_Callback
	GetSqlServerConn().Table("test.WMS_CY_Callback").Find(&result)
	fmt.Println(result)
}

func TestSqlServerDatetime(t *testing.T) {
	now := time.Now()
	fmt.Println(now.Local().Format("yyyy"))
}
