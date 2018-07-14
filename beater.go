package xlog

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis"
)

// ErrBeaterTimeout timeout error for beater, should be OK
var ErrBeaterTimeout = errors.New("beater: timeout")

// ErrBeaterMalformed malformed beat entry, should be ignored
var ErrBeaterMalformed = errors.New("beater: malformed beat entry")

// Beater a redis client wrapper for filebeat events
type Beater struct {
	Client *redis.Client

	// Key LIST key for BLPOP
	Key string
}

// DialBeater dial a redis and create *Beater
func DialBeater(url string, opt Options) (b *Beater, err error) {
	var ropt *redis.Options
	// parse redis.Options
	if ropt, err = redis.ParseURL(url); err != nil {
		return
	}
	// create redis client and ping
	var client = redis.NewClient(ropt)
	if err = client.Ping().Err(); err != nil {
		return
	}
	b = &Beater{Client: client, Key: opt.Redis.Key}
	return
}

// Close close the underlaying redis client
func (b *Beater) Close() error {
	return b.Client.Close()
}

// NextEvent fetch next event
func (b *Beater) NextEvent(be *BeatEntry) (err error) {
	var ret []string
	var raw string
	// clear
	*be = BeatEntry{}
	// blpop redis
	ret, err = b.Client.BLPop(time.Second*3, b.Key).Result()
	// convert redis.Nil to ErrBeaterTimeout
	if err != nil {
		if err == redis.Nil {
			err = ErrBeaterTimeout
		}
		return
	}
	// length 0 for timeout
	if len(ret) == 0 {
		err = ErrBeaterTimeout
		return
	}
	// length 1 for single key, 2 for multiple key, so ret[-1] should be fine
	raw = ret[len(ret)-1]
	// unmarshal json
	if err = json.Unmarshal([]byte(raw), be); err != nil {
		err = ErrBeaterMalformed
	}
	return
}
