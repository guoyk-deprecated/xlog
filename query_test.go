package xlog

import (
	"testing"
)

func TestQuery_ToURLQuery(t *testing.T) {
	q := &Query{}
	q.Hostname = "beat1"
	if q.ToURLQuery() != "hostname=beat1" {
		t.Fatal("failed")
	}
}
