package xlog

import "testing"

func TestParseTimeOfDay(t *testing.T) {
	td := ParseTimeOfDay("1201")
	if td.Hour() != 12 {
		t.Error("failed")
	}
	if td.Minute() != 1 {
		t.Error("failed")
	}
}

func TestParseTimeRangeOfDay(t *testing.T) {
	td1, td2 := ParseTimeRangeOfDay("120123-14:23:34")
	if td1.Hour() != 12 {
		t.Error("failed")
	}
	if td1.Second() != 23 {
		t.Error("failed")
	}
	if td2.Minute() != 23 {
		t.Error("failed")
	}
	if td2.Second() != 34 {
		t.Error("failed")
	}
}
