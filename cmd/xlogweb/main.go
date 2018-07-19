package main

import (
	"net/http"

	"github.com/novakit/nova"
	"github.com/novakit/static"
	"github.com/novakit/view"
	"github.com/yankeguo/xlog"
	_ "github.com/yankeguo/xlog/web" // compiled binfs
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

	n := nova.New()
	n.Env = nova.Env(options.Env())
	n.Use(static.Handler(static.Options{BinFS: !options.Dev}))
	n.Use(view.Handler(view.Options{BinFS: !options.Dev}))
	n.Use(modules.Handler(options))
	routes.Route(n)
	http.ListenAndServe(options.Web.Addr(), n)
}
