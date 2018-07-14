package routes

import (
	"github.com/novakit/nova"
	"github.com/novakit/router"
	"github.com/novakit/view"
)

func routeCollectionsShow(c *nova.Context) (err error) {
	v := view.Extract(c)
	id := router.PathParams(c).Get("id")
	v.Text(id)
	return
}
