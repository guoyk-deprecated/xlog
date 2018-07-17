package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/novakit/nova"
	"github.com/novakit/router"
	"github.com/novakit/view"
	"github.com/yankeguo/xlog"
	"github.com/yankeguo/xlog/web/modules"
)

// Route mount all routes on nova.Nova
func Route(n *nova.Nova) {
	r := router.Route(n)
	r.Get("/").Use(routeIndex)
	r.Get("/:date").Use(routeShow)
	r.Get("/:date/hints").Use(routeHints)
}

func routeIndex(c *nova.Context) (err error) {
	now := time.Now()
	http.Redirect(
		c.Res,
		c.Req,
		fmt.Sprintf("/%04d-%02d-%02d", now.Year(), now.Month(), now.Day()),
		http.StatusFound)
	return
}

func routeShow(c *nova.Context) (err error) {
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
	v.HTML("show")
	return
}

func routeHints(c *nova.Context) (err error) {
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
	ret := map[string]interface{}{}
	for _, field := range xlog.DistinctFields {
		sub := []string{}
		if err = d.Collection(date).Distinct(field, &sub); err != nil {
			return
		}
		ret[field] = sub
	}
	v.JSON(ret)
	return
}
