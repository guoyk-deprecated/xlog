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
	r.Post("/api/search").Use(routeSearch)
	// trends API
	r.Post("/api/trends").Use(routeTrends)
}

// V short-hand for view.Extract(c)
func V(c *nova.Context) *view.View {
	return view.Extract(c)
}

// D short-hand for modules.Database(c)
func D(c *nova.Context) *xlog.Database {
	return modules.Database(c)
}
