package xlog

import (
	"testing"
)

func TestQuery_Encode(t *testing.T) {
	q := &Query{}
	q.Hostname = "beat1"
	if q.Encode() != "hostname=beat1" {
		t.Fatal("failed")
	}
}
