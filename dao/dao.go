/**
 * 数据访问持久层
 */

package dao

import (
	"crf-mold/base"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

var defaultDb *gorm.DB

func InitDB() {
	//InitMysqlDB()
	InitPostGreSqlDB()
	//defaultDb = mysqlDb
	defaultDb = postGreSqlDb
}

// 获取默认的数据源连接
func GetConn() *gorm.DB {
	return defaultDb.Session(&gorm.Session{CreateBatchSize: 1000})
}

/**
 * 依据模型中的字段拼接查询条件
 */
func BuildWhereCondition(tx *gorm.DB, model interface{}) *gorm.DB {
	srcType, srcValue := reflect.TypeOf(model), reflect.ValueOf(model)

	// src必须为结构体或者结构体指针
	if srcType.Kind() == reflect.Ptr {
		srcType, srcValue = srcType.Elem(), srcValue.Elem()
	}
	if srcType.Kind() != reflect.Struct {
		panic("src type should be a struct or a struct pointer")
	}

	// 遍历所有属性
	propertyNums := srcType.NumField()

	tx.Where("1=1")

	for i := 0; i < propertyNums; i++ {
		// 属性名
		property := srcType.Field(i)
		// 属性值
		propertyValue := srcValue.FieldByName(property.Name)

		// 指针变量为nil
		if propertyValue.Kind() == reflect.Ptr && propertyValue.IsNil() {
			continue
		}

		var sqlValue = propertyValue.Interface()

		if sqlValue == nil || base.IsZeroValue(sqlValue) {
			continue
		}

		// 额外处理下base.Time类型,转成string,maybe 不需要
		_, ok := sqlValue.(base.Time)
		if ok {
			continue
		}

		// 获取列名
		gorm := property.Tag.Get("gorm")
		if gorm != "" {
			split := strings.Split(gorm, ";")
			if len(split) > 0 {
				for _, v := range split {
					sp := strings.Split(v, ":")
					if len(sp) == 2 {
						if sp[0] == "column" {
							w := fmt.Sprintf("%s = ?", sp[1])
							tx.Where(w, sqlValue)
							break
						}
					}
				}
			}
		}
	}

	return tx
}

func TransactionRollback(tx *gorm.DB) {
	if r := recover(); r != nil {
		tx.Rollback()
		// 往外抛给recover中间件
		panic(r)
	}
}
