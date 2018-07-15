package routes

import (
	"time"

	"github.com/novakit/nova"
	"github.com/novakit/router"
	"github.com/novakit/view"
	"github.com/yankeguo/xlog"
	"github.com/yankeguo/xlog/web/modules"
)

func routeCollectionsShow(c *nova.Context) (err error) {
	// variables
	v := view.Extract(c)
	d := modules.Database(c)
	// collection id
	date := router.PathParams(c).Get("date")
	// day containing year/month/day
	var day time.Time
	if day, err = time.Parse("2006-01-02", date); err != nil {
		return
	}
	// collection toal count
	var totalCount int
	if totalCount, err = d.Collection(day).Count(); err != nil {
		return
	}
	// query
	var queryCount int
	var results []xlog.LogEntry
	q := &xlog.Query{}
	if err = q.Decode(c.Req); err != nil {
		return
	}
	if err = q.Count(d, day, &queryCount); err != nil {
		return
	}
	if err = q.Execute(d, day, &results); err != nil {
		return
	}
	// render
	v.Data["Query"] = q
	v.Data["Results"] = results
	v.Data["DatabaseName"] = d.DB.Name
	v.Data["CollectionTotalCount"] = totalCount
	v.Data["QueryCount"] = queryCount
	v.Data["CollectionDate"] = date
	v.Data["CollectionName"] = d.CollectionPrefix + date
	v.HTML("collections/show")
	return
}
