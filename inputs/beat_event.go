package inputs

import (
	"regexp"
	"strings"
	"time"
	"github.com/yankeguo/xlog"
)

var (
	timestampLayout = "2006/01/02 15:04:05.000"
	linePattern     = regexp.MustCompile(`^\[(\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}\.\d{3})\]`)
	cridPattern     = regexp.MustCompile(`CRID\[([0-9a-zA-Z\-]+)\]`)
)

// BeatInfo beat info field
type BeatInfo struct {
	Hostname string `json:"hostname"` // hostname
}

// BeatEvent a single event in redis LIST sent by filebeat
type BeatEvent struct {
	Beat    BeatInfo `json:"beat"`    // contains hostname
	Message string   `json:"message"` // contains timestamp, crid
	Source  string   `json:"source"`  // contains env, topic, project
}

// ToRecord implements xlog.RecordConvertible
func (b BeatEvent) ToRecord() (r xlog.Record, err xlog.RecordConversionError) {
	// assign hostname
	r.Hostname = b.Beat.Hostname
	// decode message field
	if ok := decodeBeatMessage(b.Message, &r); !ok {
		err = xlog.NewRecordConversionError("invalid source")
		return
	}
	// decode source field
	if ok := decodeBeatSource(b.Source, &r); !ok {
		err = xlog.NewRecordConversionError("invalid message")
		return
	}
	return
}

func decodeBeatMessage(raw string, r *xlog.Record) bool {
	var err error
	var match []string
	// trim message
	raw = strings.TrimSpace(raw)
	// search timestamp
	if match = linePattern.FindStringSubmatch(raw); len(match) != 2 {
		return false
	}
	// parse timestamp
	if r.Timestamp, err = time.Parse(timestampLayout, match[1]); err != nil {
		return false
	}
	// trim message
	r.Message = strings.TrimSpace(raw[len(match[0]):])
	// find crid
	if match = cridPattern.FindStringSubmatch(r.Message); len(match) == 2 {
		r.Crid = match[1]
	} else {
		r.Crid = "-"
	}
	return true
}

func decodeBeatSource(raw string, r *xlog.Record) bool {
	var cs []string
	// trim source
	raw = strings.TrimSpace(raw)
	if cs = strings.Split(raw, "/"); len(cs) < 3 {
		return false
	}
	// assign fields
	r.Env, r.Topic, r.Project = cs[len(cs)-3], cs[len(cs)-2], cs[len(cs)-1]
	// sanitize dot separated filename
	var ss []string
	if ss = strings.Split(r.Project, "."); len(ss) > 0 {
		r.Project = ss[0]
	}
	return true
}
