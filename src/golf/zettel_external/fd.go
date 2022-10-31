package zettel_external

import (
	"path"

	"github.com/friedenberg/zit/src/charlie/ts"
)

type FD struct {
	Path    string
	ModTime ts.Time
}

func (f FD) String() string {
	return f.Path
}

func (e FD) Ext() string {
	return path.Ext(e.Path)
}

func (f FD) IsEmpty() bool {
	if f.Path == "" {
		return true
	}

	// if f.ModTime.IsZero() {
	// 	return true
	// }

	return false
}
