package time

import "testing"

func TestTodayStartEndTime(t *testing.T) {
	s, e := TodayStartEndTime()
	t.Logf("开始时间：%s\n结束时间：%s\n", s.Format(DefaultLayout), e.Format(DefaultLayout))
}
