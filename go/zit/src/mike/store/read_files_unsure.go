package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type UnsureMatchMap struct {
	sku.UnsureMatchType
	lookup map[sha.Bytes]sku.CheckedOutLikeMutableSet
}

func (s *Store) QueryUnsure(
	qg sku.ExternalQuery,
	o sku.UnsureMatchOptions,
	f sku.IterMatching,
) (err error) {
	matchMaps := o.MakeMatchMap()

	if err = s.cwdFiles.QueryUnsure(
		qg,
		sku.MakeUnsureMatchMapsCollector(matchMaps),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO create a new query group for all of history
	qg.QueryGroup.SetIncludeHistory()

	if matchMaps.Len() == 0 {
		return
	}

	if err = s.QueryWithKasten(
		qg,
		sku.MakeUnsureMatchMapsMatcher(
			matchMaps,
			f,
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
