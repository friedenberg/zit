package zettel_external

import "github.com/friedenberg/zit/src/delta/ts"

type FD struct {
	Path    string
	ModTime ts.Time
}

func (f FD) String() string {
	return f.Path
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
