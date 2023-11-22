/**
User: cr-mao
Date: 2023/10/29 09:37
Email: crmao@qq.com
Desc: time.go
*/
package gtime

import (
	"strconv"
	"time"
)

func TodayStr() string {
	return time.Now().Format("20060102")
}

func TodayInt32() int32 {
	day := TodayStr()
	dayInt, _ := strconv.Atoi(day)
	return int32(dayInt)
}

func DayBeginUnix() int64 {
	t := time.Now()
	addTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	return addTime.Unix()
}

func NextDayBeginUnix() int64 {
	return DayBeginUnix() + 86400
}

// 20231115这样的格式
func NextDayBeginUnixFromDay(timeStr string) int64 {
	//要转换成时间日期的格式模板（go诞生时间，模板必须是这个时间）
	timeTmeplate := "20060102"
	//解析日期时间字符串 , 启动的时候 time.Local 这个会改掉
	tim, _ := time.ParseInLocation(timeTmeplate, timeStr, time.Local)
	//获取该日期时间的时间戳
	return tim.Unix() + 86400
}

// 周一计算的
func WeekBeginUnix() int64 {
	t1 := time.Now()
	today := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, t1.Location())

	subNum := int(today.Weekday()) - 1
	if today.Weekday() == time.Sunday {
		subNum = 6
	}
	// 1   0
	// 2   1
	// 3   2
	return today.AddDate(0, 0, -subNum).Unix()
}

// 周一计算的
func NextWeekBeginUnix() int64 {
	beginAt := WeekBeginUnix()
	return beginAt + 86400*7
}

// 获得第几年的第几周 202346  //2023年46周, 这个是按周日算第一天的
//func YearWeek() int32 {
//	t := time.Now()
//	yearDay := t.YearDay()
//	yearFirstDay := t.AddDate(0, 0, -yearDay+1)
//	firstDayInWeek := int(yearFirstDay.Weekday())
//	firstWeekDays := 1
//	if firstDayInWeek != 0 {
//		firstWeekDays = 7 - firstDayInWeek + 1
//	}
//	var week int
//	if yearDay <= firstWeekDays {
//		week = 1
//	} else {
//		week = (yearDay-firstWeekDays)/7 + 2
//	}
//	res := t.Year()*100 + week
//	return int32(res)
//}

func YearWeekBeginMonday() int32 {
	t := time.Now()
	year, week := t.ISOWeek()
	res := year*100 + week
	return int32(res)
}
