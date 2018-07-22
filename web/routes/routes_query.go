package routes

import (
	"encoding/json"
	"github.com/novakit/nova"
	"github.com/yankeguo/xlog"
	"github.com/yankeguo/xlog/web/modules"
)

func routeQuery(c *nova.Context) (err error) {
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
	// execute query
	var ret xlog.Result
	if ret, err = db.Search(q); err != nil {
		return
	}
	// render result
	v := V(c)
	v.Data["query"] = q
	v.Data["result"] = ret
	v.DataAsJSON()
	return
}
