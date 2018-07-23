package xlog

import (
	"log"
	"testing"
	"time"
)

var testPeriodCases = []struct {
	P  Period
	Ts []Period
}{
	{
		// should be 15 minutes
		P: Period{
			Beginning: time.Date(2018, time.July, 20, 10, 44, 20, 0, time.UTC),
			End:       time.Date(2018, time.July, 20, 12, 44, 30, 0, time.UTC),
		},
		Ts: []Period{
			{
				Beginning: time.Date(2018, time.July, 20, 10, 30, 00, 0, time.UTC),
				End:       time.Date(2018, time.July, 20, 10, 45, 00, 0, time.UTC),
			},
			{
				Beginning: time.Date(2018, time.July, 20, 10, 45, 00, 0, time.UTC),
				End:       time.Date(2018, time.July, 20, 11, 00, 00, 0, time.UTC),
			},
			{
				Beginning: time.Date(2018, time.July, 20, 11, 00, 00, 0, time.UTC),
				End:       time.Date(2018, time.July, 20, 11, 15, 00, 0, time.UTC),
			},
			{
				Beginning: time.Date(2018, time.July, 20, 11, 15, 00, 0, time.UTC),
				End:       time.Date(2018, time.July, 20, 11, 30, 00, 0, time.UTC),
			},
			{
				Beginning: time.Date(2018, time.July, 20, 11, 30, 00, 0, time.UTC),
				End:       time.Date(2018, time.July, 20, 11, 45, 00, 0, time.UTC),
			},
			{
				Beginning: time.Date(2018, time.July, 20, 11, 45, 00, 0, time.UTC),
				End:       time.Date(2018, time.July, 20, 12, 00, 00, 0, time.UTC),
			},
			{
				Beginning: time.Date(2018, time.July, 20, 12, 00, 00, 0, time.UTC),
				End:       time.Date(2018, time.July, 20, 12, 15, 00, 0, time.UTC),
			},
			{
				Beginning: time.Date(2018, time.July, 20, 12, 15, 00, 0, time.UTC),
				End:       time.Date(2018, time.July, 20, 12, 30, 00, 0, time.UTC),
			},
			{
				Beginning: time.Date(2018, time.July, 20, 12, 30, 00, 0, time.UTC),
				End:       time.Date(2018, time.July, 20, 12, 45, 00, 0, time.UTC),
			},
		},
	},
}

func TestPeriod_TrendPeriods(t *testing.T) {
	for _, c := range testPeriodCases {
		ts := c.P.TrendPeriods()
		if len(ts) != len(c.Ts) {
			t.Error("failed", c.P, len(ts), "!=", len(c.Ts))
			log.Println(ts)
			break
		}
		for i := range ts {
			if ts[i] != c.Ts[i] {
				t.Error("failed", c.P, "no.", i, ts[i], "vs", c.Ts[i])
			}
		}
	}
}
