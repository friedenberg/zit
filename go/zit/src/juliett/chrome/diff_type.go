package chrome

import (
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
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

func (c *Chrome) getDiff(kinder, mutter *sku.Transacted) (dt diff, err error) {
	dt.diffType = diffTypeIgnore

	if mutter == nil {
		if dt, err = c.getDiffKinderOnly(kinder); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if metadatei.EqualerSansTaiIncludeVirtual.Equals(
		&kinder.Metadatei,
		&mutter.Metadatei,
	) {
		return
	}

	kees := etiketten(kinder)
	mees := etiketten(mutter)

	kinderHasChrome := kees.ContainsKey("%chrome")
	mutterHasChrome := mees.ContainsKey("%chrome")

	switch {
	case kinderHasChrome && mutterHasChrome:
		dt.diffType = diffTypeChange

	case kinderHasChrome:
		dt.diffType = diffTypeNew

	case mutterHasChrome:
		dt.diffType = diffTypeDelete
	}

	return
}

func (c *Chrome) getDiffKinderOnly(kinder *sku.Transacted) (dt diff, err error) {
	dt.diffType = diffTypeIgnore

	if !kinder.GetTyp().Equals(c.typ) {
		return
	}

	ees := etiketten(kinder)

	if !ees.ContainsKey("%chrome") {
		return
	}

	dt.diffType = diffTypeNew

	return
}
