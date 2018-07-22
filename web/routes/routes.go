package routes

import (
	"github.com/novakit/nova"
	"github.com/novakit/router"
	"github.com/novakit/view"
	"github.com/yankeguo/xlog"
	"github.com/yankeguo/xlog/web/modules"
)

// Route mount all routes on nova.Nova
func Route(n *nova.Nova) {
	r := router.Route(n)
	// query API
	r.Post("/api/query").Use(routeQuery)
}

// V short-hand for view.Extract(c)
func V(c *nova.Context) *view.View {
	return view.Extract(c)
}

// D short-hand for modules.Database(c)
func D(c *nova.Context) *xlog.Database {
	return modules.Database(c)
}

/*
func dateInfo(t time.Time) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	diff := t.Sub(today) / (time.Hour * 24)
	if diff == 0 {
		return "今天"
	} else if diff > 0 {
		return fmt.Sprintf("%d天后", diff)
	} else {
		return fmt.Sprintf("%d天前", -diff)
	}
}

func extractURLDate(c *nova.Context) (date time.Time, dateStr string, err error) {
	dateStr = router.PathParams(c).Get("date")
	date, err = time.Parse("2006-01-02", dateStr)
	return
}

func formatURLDate(t time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
}

func routeIndex(c *nova.Context) (err error) {
	http.Redirect(c.Res, c.Req, "/"+formatURLDate(time.Now()), http.StatusFound)
	return
}

func routeShow(c *nova.Context) (err error) {
	// variables
	v := view.Extract(c)
	d := modules.Database(c)
	// collection date
	var date time.Time
	var dateStr string
	if date, dateStr, err = extractURLDate(c); err != nil {
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
	if q, err = xlog.ParseForm(c.Req); err != nil {
		return
	}
	// results
	var results []xlog.Record
	if err = coll.Execute(q, &results); err != nil {
		return
	}
	// render
	v.Data["Query"] = q
	v.Data["Date"] = dateStr
	v.Data["DateInfo"] = dateInfo(date)
	v.Data["TotalCount"] = totalCount
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
	if date, _, err = extractURLDate(c); err != nil {
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

*/
