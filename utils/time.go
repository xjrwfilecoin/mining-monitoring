package utils

import (
	"mining-monitoring/config"
	"fmt"
	"time"
)

// 一个月的天数
func MonthDays(now time.Time) (int64, error) {
	month := int(now.Month())
	var sMonth string
	if month < 10 {
		sMonth = fmt.Sprintf("0%d", month)
	} else {
		sMonth = fmt.Sprintf("%d", month)
	}
	sTime, err := ParseLocalTime(fmt.Sprintf("%d-%s-01 00:00:00", now.Year(), sMonth))
	if err != nil {
		return 0, err
	}
	var eMonth string
	var year = now.Year()
	tMonth := month + 1
	if tMonth > 12 {
		eMonth = "01"
		year = now.Year() + 1
	} else if tMonth < 10 {
		eMonth = fmt.Sprintf("0%d", tMonth)
	} else {
		eMonth = fmt.Sprintf("%d", tMonth)
	}
	eTime, err := ParseLocalTime(fmt.Sprintf("%d-%s-01 00:00:00", year, eMonth))
	if err != nil {
		return 0, err
	}
	return int64(eTime.Sub(sTime).Hours() / 24), nil
}

// 获取本自然周的第一天的日期
func GetWeekFirstDayTime() time.Time {
	var weeks = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	now := time.Now()
	var firstTime = time.Now()
	for i := 0; i < len(weeks); i++ {
		if weeks[i] == now.Weekday().String() {
			duration := time.Duration(-(i*24+now.Hour())*3600 - 60*now.Minute() - now.Second())
			firstTime = time.Now().Add(duration * time.Second)
		}
	}
	return firstTime
}

func GetMongoTime() time.Time {
	return time.Now().Add(8 * time.Hour)
}

func ParseLocalTime(t string) (time.Time, error) {
	//中国时区
	location, err := time.LoadLocation(config.LocationTimeZone)
	ntime, err := time.ParseInLocation(config.SysTimefrom, t, location)
	return ntime, err
}
