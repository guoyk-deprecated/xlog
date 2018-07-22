package xlog

import (
	"time"
	"strings"
	"fmt"
)

var (
	// TimeOfDayLayouts supported time layouts in time-of-day
	TimeOfDayLayouts = []string{
		"15:04",
		"1504",
		"15:04:05",
		"150405",
	}

	// StorageSizeSuffixes storage size suffixes
	StorageSizeSuffixes = []struct {
		S string
		N int
	}{
		{"TB", 1024 * 1024 * 1024 * 1024},
		{"GB", 1024 * 1024 * 1024},
		{"MB", 1024 * 1024},
		{"KB", 1024},
		{"B", 1},
	}
)

// ParseTimeOfDay parse a time-of-day with multiple layouts
func ParseTimeOfDay(s string) (t time.Time) {
	var err error
	if len(s) > 0 {
		for _, layout := range TimeOfDayLayouts {
			if t, err = time.Parse(layout, s); err == nil {
				break
			}
		}
	}
	return
}

// ParseTimeRangeOfDay parse a time-range-of-day with multiple layouts
func ParseTimeRangeOfDay(s string) (begin time.Time, end time.Time) {
	times := strings.Split(s, "-")
	if len(times) > 0 {
		begin = ParseTimeOfDay(strings.TrimSpace(times[0]))
		if !begin.IsZero() {
			if len(times) > 1 {
				end = ParseTimeOfDay(strings.TrimSpace(times[1]))
			}
			if end.IsZero() {
				end = begin.Add(time.Minute)
			}
		}
	}
	return
}

// FormatStorageSize format a storage size
func FormatStorageSize(bytes int) string {
	if bytes <= 0 {
		return "0"
	}
	for _, s := range StorageSizeSuffixes {
		if bytes > s.N {
			return fmt.Sprintf("%.2f%s", float64(bytes)/float64(s.N), s.S)
		}
	}
	return ""
}
