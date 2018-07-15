package routes

import (
	"github.com/novakit/nova"
	"github.com/novakit/router"
)

// Route mount all routes on nova.Nova
func Route(n *nova.Nova) {
	r := router.Route(n)
	r.Get("/").Use(routeIndex)
	r.Get("/collections/:date").Use(routeCollectionsShow)
}
