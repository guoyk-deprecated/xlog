package outputs

import (
	"github.com/yankeguo/xlog"
	"io"
)

type Output interface {
	io.Closer
	Insert(xlog.RecordConvertible) error
}
