package browser

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

//go:generate stringer -type=diffType
type diffType int

const (
	diffTypeIgnore = diffType(iota)
	diffTypeNew
	diffTypeChange
	diffTypeDelete
)

type diff struct {
	diffType
}

func (c *Store) getDiff(kinder, mutter *sku.Transacted) (dt diff, err error) {
	dt.diffType = diffTypeIgnore

	if mutter == nil {
		if dt, err = c.getDiffKinderOnly(kinder); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if object_metadata.EqualerSansTaiIncludeVirtual.Equals(
		&kinder.Metadata,
		&mutter.Metadata,
	) {
		return
	}

	kees := etiketten(kinder)
	mees := etiketten(mutter)

	kinderHasBrowser := kees.ContainsKey("%browser")
	mutterHasBrowser := mees.ContainsKey("%browser")

	switch {
	case kinderHasBrowser && mutterHasBrowser:
		dt.diffType = diffTypeChange

	case kinderHasBrowser:
		dt.diffType = diffTypeNew

	case mutterHasBrowser:
		dt.diffType = diffTypeDelete
	}

	return
}

func (c *Store) getDiffKinderOnly(kinder *sku.Transacted) (dt diff, err error) {
	dt.diffType = diffTypeIgnore

	if !kinder.GetType().Equals(c.typ) {
		return
	}

	ees := etiketten(kinder)

	if !ees.ContainsKey("%browser") {
		return
	}

	dt.diffType = diffTypeNew

	return
}

func etiketten(sk *sku.Transacted) ids.TagSet {
	return ids.ExpandMany(sk.Metadata.GetTags(), expansion.ExpanderRight)
}
