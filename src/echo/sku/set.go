package sku

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/int_value"
)

type MutableSet struct {
	objekten map[string][]Sku
	count    int
}

func MakeMutableSet() MutableSet {
	return MutableSet{
		objekten: make(map[string][]Sku),
	}
}

func (os *MutableSet) Len() int {
	return os.count
}

type SkuLike interface {
	GetKey() string
	SetTransactionIndex(int)
	Sku() Sku
}

func (os *MutableSet) Add2(o SkuLike) {
	os.count++
	k := o.GetKey()
	s, _ := os.objekten[k]
	o.SetTransactionIndex(len(s))
	s = append(s, o.Sku())
	os.objekten[k] = s

	return
}

func (os *MutableSet) Add(o Sku) (i int) {
	os.count++
	k := o.GetKey()
	s, _ := os.objekten[k]
	i = len(s)
	s = append(s, o)
	os.objekten[k] = s

	return
}

func (os MutableSet) Get(k string) []Sku {
	return os.objekten[k]
}

func (os MutableSet) EachWithIndex(w WriterWithIndexFunc) (err error) {
	for _, oss := range os.objekten {
		for i, o := range oss {
			o1 := Indexed{
				Sku:   o,
				Index: int_value.Make(i),
			}

			if err = w(&o1); err != nil {
				switch {
				case errors.IsEOF(err):
					err = nil
					return

				default:
					err = errors.Wrap(err)
					return
				}
			}
		}
	}

	return
}

func (os MutableSet) Each(w WriterFunc) (err error) {
	return os.EachWithIndex(
		func(o *Indexed) (err error) {
			return w(&o.Sku)
		},
	)
}
