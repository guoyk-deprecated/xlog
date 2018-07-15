package routes

import (
	"fmt"
	"net/http"
	"time"

	"github.com/novakit/nova"
)

func routeIndex(c *nova.Context) (err error) {
	now := time.Now()
	http.Redirect(
		c.Res,
		c.Req,
		fmt.Sprintf("/collections/%04d-%02d-%02d", now.Year(), now.Month(), now.Day()),
		http.StatusFound)
	return
}
