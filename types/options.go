package types

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Options beatit options structure
type Options struct {
	// Redis redis options
	Redis RedisOptions `yaml:"redis"`
	// Mongo mongo options
	Mongo MongoOptions `yaml:"mongo"`
}

// RedisOptions redis options
type RedisOptions struct {
	// URLs redis urls
	URLs []string `yaml:"urls"`
	// Key redis key for BLPOP
	Key string `yaml:"key"`
}

// MongoOptions mongo options
type MongoOptions struct {
	// URL mongo url
	URL string `yaml:"url"`
	// DB mongo db name
	DB string `yaml:"db"`
	// Collection mongo collection name
	Collection string `yaml:"collection"`
}

// ReadOptionsFile load options from file
func ReadOptionsFile(file string, opt *Options) (err error) {
	var buf []byte
	if buf, err = ioutil.ReadFile(file); err != nil {
		return
	}
	if err = yaml.Unmarshal(buf, opt); err != nil {
		return
	}
	if len(opt.Redis.URLs) == 0 {
		err = errors.New("no redis urls in config")
		return
	}
	for _, u := range opt.Redis.URLs {
		if len(u) == 0 {
			err = errors.New("found invalid redis url")
			return
		}
	}
	if len(opt.Redis.Key) == 0 {
		err = errors.New("no redis key in config")
		return
	}
	if len(opt.Mongo.URL) == 0 {
		err = errors.New("no mongo url in config")
		return
	}
	if len(opt.Mongo.DB) == 0 {
		err = errors.New("no mongo db in config")
		return
	}
	if len(opt.Mongo.Collection) == 0 {
		err = errors.New("no mongo collection in config")
		return
	}
	return
}
