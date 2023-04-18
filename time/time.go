package time

import (
	stime "time"
)

// TodayStartEndTime 今天开始时间和结束时间
func TodayStartEndTime() (stime.Time, stime.Time) {
	snow := stime.Now()
	start := stime.Date(snow.Year(), snow.Month(), snow.Day(), 0, 0, 0, 0, snow.Location())
	end := stime.Date(snow.Year(), snow.Month(), snow.Day(), 23, 59, 59, 0, snow.Location())
	return start, end
}
