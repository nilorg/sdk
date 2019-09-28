package time

import (
	stime "time"
)

// TodayStartEndTime 今天开始时间和结束时间
func TodayStartEndTime() (stime.Time, stime.Time) {
	today := stime.Now().Format(DateLayout)
	start, _ := stime.Parse(DefaultLayout, today+" 00:00:00")
	end, _ := stime.Parse(DefaultLayout, today+" 23:59:59")
	return start, end
}
