package bytes

import (
	"fmt"
)

// ByteSize ...
type ByteSize float64

const (
	_ = iota // 忽略0
	// KB ...
	KB ByteSize = 1 << (10 * iota) // 1 << (10*1)
	// MB ...
	MB // 1 << (10*2)
	// GB ...
	GB // 1 << (10*3)
	// TB ...
	TB // 1 << (10*4)
	// PB ...
	PB // 1 << (10*5)
	// EB ...
	EB // 1 << (10*6)
	// ZB ...
	ZB // 1 << (10*7)
	// YB ...
	YB // 1 << (10*8)
)

func (b ByteSize) String() string {
	switch {
	case b >= YB:
		return fmt.Sprintf("%.2fYB", b/YB)
	case b >= ZB:
		return fmt.Sprintf("%.2fZB", b/ZB)
	case b >= EB:
		return fmt.Sprintf("%.2fEB", b/EB)
	case b >= PB:
		return fmt.Sprintf("%.2fPB", b/PB)
	case b >= TB:
		return fmt.Sprintf("%.2fTB", b/TB)
	case b >= GB:
		return fmt.Sprintf("%.2fGB", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2fMB", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2fKB", b/KB)
	}
	return fmt.Sprintf("%.2fB", b)
}
