/**
 * 基本数据类型零值判断
 */

package base

import (
	"math"
	"time"
)

/**
 * 判断指定类型是否是零值
 */
func IsZeroValue(arg interface{}) bool {
	switch v := arg.(type) {
	case float32:
		r := float64(v)
		return math.Abs(r-0) < 0.0000001
	case float64:
		return math.Abs(v-0) < 0.0000001
	case int64:
		return v == 0
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case string:
		return v == ""
	case Time:
		return time.Time(v).IsZero()
	case time.Time:
		return v.IsZero()
	}

	return false
}
