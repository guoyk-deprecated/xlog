package types

import (
	"regexp"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
)

var (
	timestampFormat = "2006/01/02 15:04:05.000"
	linePattern     = regexp.MustCompile(`^\[(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3})\]`)
	cridPattern     = regexp.MustCompile(`CRID\[([0-9a-zA-Z]+)\]`)

	// IndexedFields fields needed to be indexed
	IndexedFields = []string{"timestamp", "hostname", "env", "project", "topic", "crid"}
)

// BeatInfo beat field in a event
type BeatInfo struct {
	Hostname string `json:"hostname"`
}

// BeatEntry a single event in redis LIST sent by filebeat
type BeatEntry struct {
	// Beat the beat object
	Beat BeatInfo `json:"beat"`
	// Message the message string, contains timestamp, crid
	Message string `json:"message"`
	// Source the source string, contains env, topic, project
	Source string `json:"source"`
}

// Convert decode a BeatEntry into a LogEntry
func (b BeatEntry) Convert(le *LogEntry) (ok bool) {
	if ok = decodeMessage(b.Message, le); !ok {
		return
	}
	if ok = decodeSource(b.Source, le); !ok {
		return
	}
	le.Hostname = b.Beat.Hostname
	return true
}

func decodeMessage(raw string, le *LogEntry) (ok bool) {
	var err error
	var match []string
	// trim message
	raw = strings.TrimSpace(raw)
	// search timestamp
	if match = linePattern.FindStringSubmatch(raw); len(match) != 2 {
		return
	}
	// parse timestamp
	if le.Timestamp, err = time.Parse(timestampFormat, match[1]); err != nil {
		return
	}
	// trim message
	le.Message = strings.TrimSpace(raw[len(match[0]):])
	// find crid
	if match = cridPattern.FindStringSubmatch(le.Message); len(match) == 2 {
		le.Crid = match[1]
	}
	return true
}

func decodeSource(raw string, le *LogEntry) (ok bool) {
	var cs []string
	// trim source
	raw = strings.TrimSpace(raw)
	if cs = strings.Split(raw, "/"); len(cs) < 3 {
		return
	}
	// assign fields
	le.Env = cs[len(cs)-3]
	le.Topic = cs[len(cs)-2]
	le.Project = cs[len(cs)-1]
	// sanitize dot seperated filename
	var ps []string
	if ps = strings.Split(le.Project, "."); len(ps) > 0 {
		le.Project = ps[0]
	}
	return true
}

// LogEntry a log document in mongodb
type LogEntry struct {
	Timestamp time.Time
	Hostname  string
	Env       string
	Project   string
	Topic     string
	Crid      string
	Message   string
}

// ToBSON convert to bson.M
func (l LogEntry) ToBSON() bson.M {
	return bson.M{
		"timestamp": l.Timestamp,
		"hostname":  l.Hostname,
		"env":       l.Env,
		"project":   l.Project,
		"topic":     l.Topic,
		"crid":      l.Crid,
		"message":   l.Message,
	}
}
