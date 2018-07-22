package inputs

import (
	"github.com/yankeguo/xlog"
	"io"
)

// Input abstract a input type
type Input interface {
	io.Closer

	// Next returns next RecordConvertible, if not present, return nil, nil
	Next() (r xlog.RecordConvertible, err error)

	// Recover recover a RecordConvertible
	Recover(r xlog.RecordConvertible) error
}
