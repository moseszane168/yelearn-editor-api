package base

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateUniqueTextID() string {
	// 使用当前时间戳和随机数生成唯一字符串
	timestamp := time.Now().Format("20060102150405")
	randomNum := rand.Intn(10000)
	return fmt.Sprintf("%s%d", timestamp, randomNum)
}
