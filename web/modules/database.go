package modules

import (
	"github.com/novakit/nova"
	"github.com/yankeguo/xlog"
)

// DatabaseContextKey context key for database
const DatabaseContextKey = "database"

// Database extract database from nova.Context
func Database(c *nova.Context) (db *xlog.Database) {
	db, _ = c.Values[DatabaseContextKey].(*xlog.Database)
	return
}
