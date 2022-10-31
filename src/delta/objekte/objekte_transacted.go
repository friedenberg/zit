package objekte

import "github.com/friedenberg/zit/src/charlie/ts"

type ObjekteTransacted struct {
	ObjekteWithIndex
	Schwanz ts.Time
}

func (a ObjekteTransacted) Equals(b ObjekteTransacted) (ok bool) {
	if !a.ObjekteWithIndex.Equals(b.ObjekteWithIndex) {
		return
	}

	if !a.Schwanz.Equals(b.Schwanz) {
		return
	}

	ok = true

	return
}

func (a ObjekteTransacted) Less(b ObjekteTransacted) (ok bool) {
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
