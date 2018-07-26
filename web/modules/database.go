package modules

import (
	"github.com/novakit/nova"
	"github.com/yankeguo/xlog/outputs"
)

// DatabaseContextKey context key for database
const DatabaseContextKey = "database"

// MongoDB extract database from nova.Context
func Database(c *nova.Context) (db *outputs.MongoDB) {
	db, _ = c.Values[DatabaseContextKey].(*outputs.MongoDB)
	return
}
