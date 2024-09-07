package base

import (
	"fmt"
	"math"
	"time"
)

func main() {

	l, _ := time.LoadLocation("Asia/Shanghai")
	startTime, _ := time.ParseInLocation("2006-01-02", "2018-12-22", l)
	endTime, _ := time.ParseInLocation("2006-01-02", "2019-05-17", l)

	datas := GroupByWeekDate(startTime, endTime)
	for _, d := range datas {
		fmt.Println(d)
	}

}

//判断时间是当年的第几周
func WeekByDate(t time.Time) (int, int) {
	yearDay := t.YearDay()
	yearFirstDay := t.AddDate(0, 0, -yearDay+1)
	firstDayInWeek := int(yearFirstDay.Weekday())

	//今年第一周有几天
	firstWeekDays := 1
	if firstDayInWeek != 0 {
		firstWeekDays = 7 - firstDayInWeek + 1
	}
	var week int
	if yearDay <= firstWeekDays {
		week = 1
	} else {
		week = (yearDay-firstWeekDays)/7 + 2
	}
	// return fmt.Sprintf("%d第%d周", t.Year(), week)
	return t.Year(), week
}

type WeekDate struct {
	Year      int
	Week      int
	StartTime time.Time
	EndTime   time.Time
}

// 将开始时间和结束时间分割为周为单位
func GroupByWeekDate(startTime, endTime time.Time) []WeekDate {
	weekDate := make([]WeekDate, 0)
	diffDuration := endTime.Sub(startTime)
	days := int(math.Ceil(float64(diffDuration/(time.Hour*24)))) + 1

	currentWeekDate := WeekDate{}
	y, w := WeekByDate(endTime)
	currentWeekDate.Year = y
	currentWeekDate.Week = w
	currentWeekDate.EndTime = endTime
	currentWeekDay := int(endTime.Weekday())
	if currentWeekDay == 0 {
		currentWeekDay = 7
	}
	currentWeekDate.StartTime = endTime.AddDate(0, 0, -currentWeekDay+1)
	nextWeekEndTime := currentWeekDate.StartTime
	weekDate = append(weekDate, currentWeekDate)

	for i := 0; i < (days-currentWeekDay)/7; i++ {
		weekData := WeekDate{}
		weekData.EndTime = nextWeekEndTime
		weekData.StartTime = nextWeekEndTime.AddDate(0, 0, -7)

		y, w := WeekByDate(weekData.StartTime)
		weekData.Year = y
		weekData.Week = w

		nextWeekEndTime = weekData.StartTime
		weekDate = append(weekDate, weekData)
	}

	if lastDays := (days - currentWeekDay) % 7; lastDays > 0 {
		lastData := WeekDate{}
		lastData.EndTime = nextWeekEndTime
		lastData.StartTime = nextWeekEndTime.AddDate(0, 0, -lastDays)
		y, w := WeekByDate(lastData.StartTime)
		lastData.Year = y
		lastData.Week = w

		weekDate = append(weekDate, lastData)
	}

	return weekDate
}
