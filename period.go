package xlog

import (
	"github.com/globalsign/mgo/bson"
	"time"
)

// Period time range
type Period struct {
	Beginning time.Time `json:"beginning,omitempty"` // time start
	End       time.Time `json:"end,omitempty"`       // time end
	Ascendant bool      `json:"ascendant" bson:"-"`  // timestamp ascendant, default to false
}

// TrendPeriods trend trend periods
func (p Period) TrendPeriods() (ret []Period) {
	ret = make([]Period, 0)
	// determine segment length
	var d time.Duration
	diff := p.End.Sub(p.Beginning)
	if diff < 0 {
		return
	}
	if diff > time.Hour*10 {
		// diff > 10 hours, segment = 1 hour
		d = time.Hour
	} else if diff > time.Hour*5 {
		// diff > 5 hours <= 10 hours, segment = 30 minutes
		d = time.Minute * 30
	} else if diff > time.Hour*2 {
		// diff > 2 hours <= 5 hours, segment = 15 minutes
		d = time.Minute * 15
	} else if diff > time.Hour {
		// diff < 2 hours > 1 hours, segment = 5 minutes
		d = time.Minute * 5
	} else if diff > time.Minute*30 {
		// diff < 1 hour > 30 minutes, segment = 2 minutes
		d = time.Minute * 2
	} else {
		d = time.Minute
	}
	// round up beginning
	rb := p.Beginning.Round(d)
	// if rounded value is greater than beginning, minus segment
	if rb.After(p.Beginning) {
		rb = rb.Add(-d)
	}
	// start loop
	var nb = rb
	for {
		ret = append(ret, Period{
			Beginning: nb,
			End:       nb.Add(d),
			Ascendant: p.Ascendant,
		})
		nb = nb.Add(d)
		if nb.After(p.End) {
			break
		}
	}
	return
}

// ToMatch to mongodb match
func (p Period) ToMatch() bson.M {
	return bson.M{
		"$gte": p.Beginning,
		"$lt":  p.End,
	}
}
