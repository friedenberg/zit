package objekte

import "github.com/friedenberg/zit/src/bravo/int_value"

type ObjekteWithIndex struct {
	Objekte
	Index int_value.IntValue
}

func (a ObjekteWithIndex) Equals(b ObjekteWithIndex) (ok bool) {
	if !a.Objekte.Equals(b.Objekte) {
		return
	}

	if !a.Index.Equals(b.Index) {
		return
	}

	return true
}
