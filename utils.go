package xlog

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	"strings"
	"time"
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

// BeginningOfDay beginning of day, to UTC
func BeginningOfDay() time.Time {
	return BeginningOfTheDay(time.Now())
}

// EndOfDay end of day, to UTC
func EndOfDay() time.Time {
	return EndOfTheDay(time.Now())
}

// BeginningOfTheDay beginning of the day specified, to UTC
func BeginningOfTheDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// EndOfTheDay end of the day specified, to UTC
func EndOfTheDay(t time.Time) time.Time {
	return BeginningOfTheDay(t).Add(time.Hour * 24)
}

// SameDay the two time is the same day
func SameDay(t1 time.Time, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
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

// BSONPutMatchField put a possible comma separated field for match
func BSONPutMatchField(m bson.M, key string, val string) {
	// skip empty value
	if len(val) == 0 {
		return
	}
	// use $in for comma separated values
	if strings.Contains(val, ",") {
		m[key] = bson.M{"$in": strings.Split(val, ",")}
	} else {
		m[key] = val
	}
	return
}

// CompactField compact query field for possible comma separated values
func CompactField(str string) string {
	str = strings.TrimSpace(str)
	if strings.Contains(str, ",") {
		var values = make([]string, 0)
		var splits = strings.Split(str, ",")
		for _, s := range splits {
			values = append(values, strings.TrimSpace(s))
		}
		return strings.Join(values, ",")
	}
	return str
}
