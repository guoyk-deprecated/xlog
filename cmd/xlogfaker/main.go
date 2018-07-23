package main

import (
	srand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"github.com/Pallinder/go-randomdata"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/yankeguo/xlog"
	"github.com/yankeguo/xlog/inputs"
	"log"
	"math/rand"
	"time"
)

var (
	options xlog.Options

	redisClients = make(map[int]*redis.Client)

	envs      = []string{"testenv1", "testenv2", "testenv3"}
	projects  = []string{"test-project-1", "test-project-2", "test-project-3"}
	topics    = []string{"test-topic-1", "test-topic-2", "test-topic-3"}
	hostnames = []string{"test.host1", "test.host2", "test.host3"}
)

func strSliceRandVal(ss []string) string {
	if len(ss) == 0 {
		return ""
	}
	return ss[rand.Intn(len(ss))]
}

func main() {
	// options flag
	var err error
	if err = xlog.ParseOptionsFlag(&options); err != nil {
		panic(err)
	}

	// check dev flag
	if !options.Dev {
		panic(errors.New("you must specify dev flag to use xlogfaker tool"))
	}

	// seed the math/rand with crypto/rand
	seed := make([]byte, 8)
	srand.Read(seed)
	rand.Seed(int64(binary.BigEndian.Uint64(seed)))

	for {
		// sleep a random time
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(10)+5) * 100)
		// choose a random redis url
		i := rand.Intn(len(options.Redis.URLs))
		// reuse redis client
		r := redisClients[i]
		// create if not existed
		if r == nil {
			ro, err := redis.ParseURL(options.Redis.URLs[i])
			if err != nil {
				panic(err)
			}
			r = redis.NewClient(ro)
			if err = r.Ping().Err(); err != nil {
				panic(err)
			}
			redisClients[i] = r
		}
		// send dirt randomly
		if rand.Intn(10) > 7 {
			if err = r.RPush(options.Redis.Key, randomdata.Address()).Err(); err != nil {
				panic(err)
			}
			continue
		}
		// crid
		crid := make([]byte, 8)
		rand.Read(crid)
		cridStr := hex.EncodeToString(crid)
		dateStr := time.Now().Format("2006/01/02 15:04:05.000")
		// create beat event
		be := inputs.BeatEvent{
			Beat: inputs.BeatInfo{
				Hostname: strSliceRandVal(hostnames),
			},
			Message: "[" + dateStr + "] CRID[" + cridStr + "] " + randomdata.Paragraph() + "\n" + randomdata.Paragraph(),
			Source:  "/var/log/" + strSliceRandVal(envs) + "/" + strSliceRandVal(topics) + "/" + strSliceRandVal(projects) + ".log",
		}
		// marshal beat event
		var buf []byte
		if buf, err = json.Marshal(&be); err != nil {
			panic(err)
		}
		// RPUSH
		if err = r.RPush(options.Redis.Key, string(buf)).Err(); err != nil {
			panic(err)
		}
		log.Println("sent", dateStr)
	}
}
