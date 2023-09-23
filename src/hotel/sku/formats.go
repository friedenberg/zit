package sku

import (
	"github.com/friedenberg/zit/src/alfa/errors"
)

type (
	FuncMakeSkuLike func(string) (*Transacted, error)
)

func TryMakeSkuWithFormats(fms ...FuncMakeSkuLike) FuncMakeSkuLike {
	return func(line string) (sk *Transacted, err error) {
		em := errors.MakeMulti()

		for _, f := range fms {
			if sk, err = f(line); err == nil {
				return
			}

			em.Add(err)
		}

		return nil, em
	}
}
