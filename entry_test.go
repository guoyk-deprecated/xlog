package xlog

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestBeatEmtry(t *testing.T) {
	buf, err := ioutil.ReadFile("testdata/beat_entry.sample.json")
	if err != nil {
		t.Fatal(err)
	}
	var be BeatEntry
	if err = json.Unmarshal(buf, &be); err != nil {
		t.Fatal(err)
	}
}
