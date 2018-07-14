package xlog

import "testing"

func TestReadOptionsFile(t *testing.T) {
	opt := &Options{}
	err := ReadOptionsFile("testdata/options.sample.yaml", opt)
	if err != nil {
		t.Fatal(err)
	}
	if opt.Redis.Key != "xlog" {
		t.Fatal("failed")
	}
	if opt.Redis.URLs[0] != "redis://127.0.0.1:6379" {
		t.Fatal("failed")
	}
	if opt.Mongo.DB != "xlog" {
		t.Fatal("failed")
	}
}
