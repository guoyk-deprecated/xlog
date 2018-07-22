package xlog

import (
	"time"
)

// TimeRange time range
type TimeRange struct {
	Begin *time.Time `json:"begin,omitempty" bson:"$gte,omitempty"` // time start
	End   *time.Time `json:"end,omitempty" bson:"$lt,omitempty"`    // time end
}

// Query query
type Query struct {
	Timestamp *TimeRange `json:"timestamp,omitempty" bson:"timestamp,omitempty"` // timestamp
	Crid      string     `json:"crid,omitempty" bson:"crid,omitempty"`           // crid
	Hostname  string     `json:"hostname,omitempty" bson:"hostname,omitempty"`   // hostname
	Env       string     `json:"env,omitempty" bson:"env,omitempty"`             // env
	Project   string     `json:"project,omitempty" bson:"project,omitempty"`     // project
	Topic     string     `json:"topic,omitempty" bson:"topic,omitempty"`         // topic
}
