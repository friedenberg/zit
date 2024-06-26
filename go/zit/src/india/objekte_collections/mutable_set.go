package objekte_collections

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
)

type MutableSet = schnittstellen.MutableSetLike[*store_fs.External]

type KeyerHinweis struct{}

func (k KeyerHinweis) GetKey(z *store_fs.External) string {
	if z == nil {
		return ""
	}

	return z.GetKennung().String()
}

func MakeMutableSetUniqueHinweis(zs ...*store_fs.External) MutableSet {
	return collections_value.MakeMutableValueSet[*store_fs.External](
		KeyerHinweis{},
		zs...)
}

type KeyerFD struct{}

func (k KeyerFD) GetKey(z *store_fs.External) string {
	if z == nil {
		return ""
	}

	return z.String()
}

func MakeMutableSetUniqueFD(zs ...*store_fs.External) MutableSet {
	return collections_value.MakeMutableValueSet[*store_fs.External](
		KeyerFD{},
		zs...)
}

type KeyerStored struct{}

func (k KeyerStored) GetKey(z *store_fs.External) string {
	if z == nil {
		return ""
	}

	if z.GetObjekteSha().IsNull() {
		return ""
	}

	return z.GetObjekteSha().String()
}

func MakeMutableSetUniqueStored(zs ...*store_fs.External) MutableSet {
	return collections_value.MakeMutableValueSet[*store_fs.External](
		KeyerStored{},
		zs...)
}

type KeyerAkte struct{}

func (k KeyerAkte) GetKey(z *store_fs.External) string {
	if z == nil {
		return ""
	}

	sh := z.GetAkteSha()

	if sh.IsNull() {
		return ""
	}

	return sh.String()
}

func MakeMutableSetUniqueAkte(zs ...*store_fs.External) MutableSet {
	return collections_value.MakeMutableValueSet[*store_fs.External](
		KeyerAkte{},
		zs...)
}
