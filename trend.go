package xlog

import "time"

type Trend struct {
	Beginning time.Time `json:"beginning"`
	End       time.Time `json:"end"`
	Count     int       `json:"count"`
}
