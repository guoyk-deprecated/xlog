package inputs

import (
	"github.com/yankeguo/xlog"
	"testing"
	"time"
)

func TestBeatEvent_LinePattern(t *testing.T) {
	if linePattern.MatchString("2018/07/20 15:03:00.000 hello world") {
		t.Fatal("failed")
	}
	if !linePattern.MatchString("[2018/07/20 15:03:00.000] hello world") {
		t.Fatal("failed")
	}
}

func TestBeatEvent_ToRecord(t *testing.T) {
	var be BeatEvent
	_, err := be.ToRecord()
	if err == nil {
		t.Fatal("failed")
	}
	be.Message = "[2018/07/20 15:03:00.000] hello world CRID[aaa]"
	_, err = be.ToRecord()
	if err == nil {
		t.Fatal("failed")
	}
	be.Source = "/tmp/test2/test3/test1.20180719.log"
	be.Beat.Hostname = "test.test"
	var r xlog.Record
	if r, err = be.ToRecord(); err != nil {
		t.Fatal(err)
	}
	rr := xlog.Record{
		Timestamp: time.Date(2018, time.July, 20, 15, 3, 0, 0, time.UTC),
		Message:   "hello world CRID[aaa]",
		Env:       "test2",
		Topic:     "test3",
		Project:   "test1",
		Crid:      "aaa",
		Hostname:  "test.test",
	}
	if r != rr {
		t.Fatal("failed")
	}
}
