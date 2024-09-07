/**
 * ID工具类
 */

package base

import (
	"strings"

	"github.com/google/uuid"
	uuid4 "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

/**
 * UUID生成
 */
func UUID() string {
	u1, err := uuid.NewUUID()
	if err != nil {
		logrus.Fatal(err)
	}

	return strings.ReplaceAll(u1.String(), "-", "")
}

//GenUUID 生成uuid
func GenUUID() string {
	uuidFunc := uuid4.NewV4()
	uuidStr := uuidFunc.String()
	uuidStr = strings.Replace(uuidStr, "-", "", -1)
	uuidByt := []rune(uuidStr)
	return string(uuidByt[8:24])
}
