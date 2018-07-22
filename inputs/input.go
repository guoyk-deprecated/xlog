package inputs

import (
	"io"
	"github.com/yankeguo/xlog"
)

// Input abstract a input type
type Input interface {
	io.Closer

	// Next returns next RecordConvertible, if not present, return nil, nil
	Next() (r xlog.RecordConvertible, err error)

	// Recover recover a RecordConvertible
	Recover(r xlog.RecordConvertible) error
}
