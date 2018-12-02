package time

import (
	stime "time"
)

// TodayStartEndTime 今天开始时间和结束时间
func TodayStartEndTime() (stime.Time, stime.Time) {
	today := stime.Now().Format("2006-01-02")
	layout := "2006-01-02 15:04:05"
	start, _ := stime.Parse(layout, today+" 00:00:00")
	end, _ := stime.Parse(layout, today+" 23:59:59")
	return start, end
}
