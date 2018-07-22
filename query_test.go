package xlog

import (
	"encoding/json"
	"github.com/globalsign/mgo/bson"
	"testing"
)

func TestQuery_Marshal(t *testing.T) {
	var err error
	var buf []byte
	q := &Query{}
	if buf, err = json.Marshal(q); err != nil {
		t.Error(err)
	}
	if string(buf) != "{}" {
		t.Error("failed")
	}
	if buf, err = bson.Marshal(q); err != nil {
		t.Error(err)
	}
}
