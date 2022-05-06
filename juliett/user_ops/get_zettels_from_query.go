package user_ops

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/delta/umwelt"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/india/store_with_lock"
)

type GetZettelsFromQuery struct {
	Umwelt *umwelt.Umwelt
}

func (c GetZettelsFromQuery) Run(query stored_zettel.NamedFilter) (result ZettelResults, err error) {
	var store store_with_lock.Store

	if store, err = store_with_lock.New(c.Umwelt); err != nil {
		err = errors.Error(err)
		return
	}

	defer errors.PanicIfError(store.Flush)

	if result.SetNamed, err = store.Zettels().Query(query); err != nil {
		err = _Error(err)
		return
	}

	return
}
