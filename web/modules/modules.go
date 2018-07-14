package modules

import (
	"github.com/novakit/nova"
	"github.com/yankeguo/xlog"
)

// Handler create nova.HandlerFunc, injects modules
func Handler(opt xlog.Options) nova.HandlerFunc {
	var err error
	var db *xlog.Database
	// panic if failed to dial database
	if db, err = xlog.DialDatabase(opt); err != nil {
		panic(err)
	}
	return func(c *nova.Context) (err error) {
		c.Values[DatabaseContextKey] = db
		c.Next()
		return
	}
}
