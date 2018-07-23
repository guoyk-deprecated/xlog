package xlog

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

// RecordConversionError implicit error type for record conversion
type RecordConversionError interface {
	error
	RecordConversionError()
}

type recordConversionError struct {
	err string
}

func (r recordConversionError) Error() string {
	return r.err
}

func (r recordConversionError) RecordConversionError() {}

// IsRecordConversionError test if a error is a record conversion error
func IsRecordConversionError(err error) (ok bool) {
	_, ok = err.(RecordConversionError)
	return
}

// NewRecordConversionError returns a new record conversion error
func NewRecordConversionError(err string) RecordConversionError {
	return recordConversionError{err: err}
}

// RecordConvertible anything that can convert to a Record
type RecordConvertible interface {
	// ToRecord convert to Record
	ToRecord() (Record, RecordConversionError)
}

// Record a log record in collection
type Record struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"` // the record id in mongodb
	Timestamp time.Time     `bson:"timestamp" json:"timestamp"`         // the time when record produced
	Hostname  string        `bson:"hostname" json:"hostname"`           // the server where record produced
	Env       string        `bson:"env" json:"env"`                     // environment where record produced, for example 'dev'
	Project   string        `bson:"project" json:"project"`             // project name
	Topic     string        `bson:"topic" json:"topic"`                 // topic of log, for example 'access', 'err'
	Crid      string        `bson:"crid" json:"crid"`                   // correlation id
	Message   string        `bson:"message" json:"message"`             // the actual log message body
}

// ToRecord implements RecordConvertible
func (r Record) ToRecord() (Record, RecordConversionError) {
	return r, nil
}
