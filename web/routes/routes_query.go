package routes

import (
	"github.com/novakit/nova"
)

func routeQuery(c *nova.Context) (err error) {
	v := V(c)
	v.DataAsJSON()
	return
}
