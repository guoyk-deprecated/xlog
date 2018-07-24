package xlog

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Options beatit options structure
type Options struct {
	// Redis redis options
	Redis RedisOptions `yaml:"redis"`
	// Mongo mongo options
	Mongo MongoOptions `yaml:"mongo"`
	// Web web options
	Web WebOptions `yaml:"web"`
	// Verbose verbose mode
	Verbose bool `yaml:"verbose"`
	// Dev dev mode
	Dev bool `yaml:"dev"`
}

// Env production / development
func (o Options) Env() string {
	if o.Dev {
		return "development"
	}
	return "production"
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
	// Tough no timeout
	Tough bool `yaml:"tough"`
}

// WebOptions web options
type WebOptions struct {
	// Host host to bind
	Host string `yaml:"host"`
	// Port port to listen
	Port string `yaml:"port"`
}

// Addr "host:port"
func (w WebOptions) Addr() string {
	return fmt.Sprintf("%s:%s", w.Host, w.Port)
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
	if len(opt.Web.Host) == 0 {
		opt.Web.Host = "127.0.0.1"
	}
	if len(opt.Web.Port) == 0 {
		opt.Web.Port = "3000"
	}
	return
}

// ParseOptionsFlag read options from command line flag
func ParseOptionsFlag(opt *Options) (err error) {
	var (
		file    string
		verbose bool
		dev     bool
	)
	flag.StringVar(&file, "c", "/etc/xlog.yml", "options file")
	flag.BoolVar(&verbose, "verbose", false, "verbose mode, overriding config file")
	flag.BoolVar(&dev, "dev", false, "dev mode, overrding config file")
	flag.Parse()

	// read options file
	if err = ReadOptionsFile(file, opt); err != nil {
		return
	}

	// override values
	if verbose {
		opt.Verbose = true
	}
	if dev {
		opt.Dev = true
	}
	return
}
