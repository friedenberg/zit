package sku

import "github.com/friedenberg/zit/src/bravo/int_value"

type Indexed struct {
	Sku
	Index int_value.IntValue
}

func (a Indexed) Equals(b Indexed) (ok bool) {
	if !a.Sku.Equals(b.Sku) {
		return
	}

	if !a.Index.Equals(b.Index) {
		return
	}

	return true
}
