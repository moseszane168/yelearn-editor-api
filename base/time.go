/**
 * 时间对象格式化拓展，格式：yyyy-MM-dd HH:mm:ss
 */

package base

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Time time.Time

const (
	TimeFormart = "2006-01-02 15:04:05"
	DateFormart = "2006-01-02"
)

// 解析JSON中的字符串到base.Time对象
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	timeStr := strings.Trim(string(data), "\"")
	if string(data) == "null" || timeStr == "" {
		return nil
	}

	t1, err := time.ParseInLocation(TimeFormart, timeStr, time.Local)
	*t = Time(t1)
	return
}

// 解析base.Time对象到JSON中的字符串
func (t Time) MarshalJSON() ([]byte, error) {
	timeStr := fmt.Sprintf("\"%v\"", time.Time(t).Format("2006-01-02 15:04:05"))
	return []byte(timeStr), nil
}

// 将自定义变量转成数据库需要的类型string或者time.Time
func (t Time) Value() (driver.Value, error) {
	// Time 转换成 time.Time 类型
	tTime := time.Time(t)
	if t.IsZero() {
		return nil, nil
	}
	return tTime.Format("2006-01-02 15:04:05"), nil
}

// 将数据库结果time.Time转成我们的指定类型
func (t *Time) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		// 字符串转成 time.Time 类型
		*t = Time(vt)
	default:
		return errors.New("类型处理错误")
	}
	return nil
}

/**
 * base.Time转string
 */
func (t Time) String() string {
	// 如果时间 null 那么我们需要把返回的值进行修改
	if t.IsZero() {
		return ""
	}

	return time.Time(t).Format(TimeFormart)
}

/**
 * 当前时间变量转为日期格式yyyy-MM-dd
 */
func (t Time) DateString() string {
	return time.Time(t).Format(DateFormart)
}

/**
 * base.Time转string
 */
func (t Time) Time() time.Time {
	return time.Time(t)
}

func (t Time) IsZero() bool {
	return time.Time(t).IsZero()
}

/**
 * 格式化golang默认时间格式为base.Time类型
 */
func FormatTime(str string) (Time, error) {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return Now(), err
	}

	return Time(t), nil
}

/**
 * 格式化base.Time时间格式
 */
func FormatBaseTime(str string) (Time, error) {
	t, err := time.Parse(TimeFormart, str)
	if err != nil {
		return Time(time.Now()), err
	}
	return Time(t), nil
}

func Now() Time {
	return Time(time.Now())
}

func ZeroValue() string {
	return Time{}.String()
}

func SameDayWithNow(t time.Time) bool {
	return t.Format(DateFormart) == time.Now().Format(DateFormart)
}

/**
 * 返回两个时间的天数差，忽略时间
 */
func DiffDay(t1, t2 time.Time) int {
	first := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	second := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)
	return int(second.Sub(first).Hours() / 12)
}

/**
 * 返回当前时间与指定时间的天数查，当前-指定
 */
func DiffDaySince(t1 time.Time) int {
	first := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2 := time.Now()
	second := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, time.Local)
	return int(second.Sub(first).Hours() / 24)
}
