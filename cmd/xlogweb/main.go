package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/novakit/nova"
	"github.com/novakit/view"
	"github.com/yankeguo/xlog"
	"github.com/yankeguo/xlog/web/modules"
	"github.com/yankeguo/xlog/web/routes"
)

var (
	options xlog.Options
)

func main() {
	var err error
	if err = xlog.ParseOptionsFlag(&options); err != nil {
		panic(err)
	}

	addr := fmt.Sprintf("%s:%s", options.Web.Host, options.Web.Port)

	n := nova.New()
	if !options.Dev {
		n.Env = nova.Production
	}
	n.Use(view.Handler(view.Options{BinFS: !options.Dev}))
	n.Use(modules.Handler(options))
	routes.Route(n)
	log.Println("listening at", addr)
	http.ListenAndServe(addr, n)
}
