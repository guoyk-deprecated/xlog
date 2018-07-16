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
	// collection date
	var date time.Time
	dateStr := router.PathParams(c).Get("date")
	if date, err = time.Parse("2006-01-02", dateStr); err != nil {
		return
	}
	// collection
	coll := d.Collection(date)
	// toal count
	var totalCount int
	if err = coll.TotalCount(&totalCount); err != nil {
		return
	}
	// stats
	var stats xlog.CollectionStats
	if err = coll.Stats(&stats); err != nil {
		return
	}
	// parse query
	var q xlog.Query
	if q, err = xlog.ParseQuery(c.Req); err != nil {
		return
	}
	// query count
	var count int
	if err = coll.Count(q, &count); err != nil {
		return
	}
	// results
	var results []xlog.LogEntry
	if err = coll.Execute(q, &results); err != nil {
		return
	}
	// render
	v.Data["Query"] = q
	v.Data["Date"] = dateStr
	v.Data["TotalCount"] = totalCount
	v.Data["Count"] = count
	v.Data["Stats"] = stats
	v.Data["Results"] = results
	v.HTML("collections/show")
	return
}
