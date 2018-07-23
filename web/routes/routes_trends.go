package routes

import (
	"encoding/json"
	"github.com/novakit/nova"
	"github.com/novakit/view"
	"github.com/yankeguo/xlog"
	"github.com/yankeguo/xlog/web/modules"
)

func routeTrends(c *nova.Context) (err error) {
	// decode query
	q := xlog.Query{}
	d := json.NewDecoder(c.Req.Body)
	if err = d.Decode(&q); err != nil {
		return
	}
	// fix query for timestamp and offset
	q = q.Validated()
	// extract db
	db := modules.Database(c)
	// execute trends
	var ts []xlog.Trend
	if ts, err = db.Trends(q); err != nil {
		return
	}
	// render
	v := view.Extract(c)
	v.Data["trends"] = ts
	v.DataAsJSON()
	return
}
