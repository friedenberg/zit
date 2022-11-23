package sku

import "github.com/friedenberg/zit/src/delta/ts"

type Transacted struct {
	Indexed
	Schwanz ts.Time
}

func (a Transacted) Equals(b Transacted) (ok bool) {
	if !a.Indexed.Equals(b.Indexed) {
		return
	}

	if !a.Schwanz.Equals(b.Schwanz) {
		return
	}

	ok = true

	return
}

func (a Transacted) Less(b Transacted) (ok bool) {
	if a.Schwanz.Less(b.Schwanz) {
		ok = true
		return
	}

	if a.Schwanz.Equals(b.Schwanz) && a.Index.Less(b.Index) {
		ok = true
		return
	}

	return
}
