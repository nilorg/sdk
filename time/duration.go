package time

import (
	"fmt"
	"math"
	"time"
)

// DurationFormatHMS 05:03:01
func DurationFormatHMS(d time.Duration) string {
	h := fmt.Sprintf("%02d", int64(math.Floor(d.Hours())))
	m := fmt.Sprintf("%02d", int64(math.Floor(d.Minutes()))%60)
	s := fmt.Sprintf("%02d", int64(math.Floor(d.Seconds()))%60)
	return fmt.Sprintf("%s:%s:%s", h, m, s)
}
