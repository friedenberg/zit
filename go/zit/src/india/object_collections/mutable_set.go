package object_collections

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
)

type MutableSet = interfaces.MutableSetLike[*store_fs.External]

type KeyerHinweis struct{}

func (k KeyerHinweis) GetKey(z *store_fs.External) string {
	if z == nil {
		return ""
	}

	return z.GetObjectId().String()
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

	if z.GetObjectSha().IsNull() {
		return ""
	}

	return z.GetObjectSha().String()
}

func MakeMutableSetUniqueStored(zs ...*store_fs.External) MutableSet {
	return collections_value.MakeMutableValueSet[*store_fs.External](
		KeyerStored{},
		zs...)
}

type KeyerBlob struct{}

func (k KeyerBlob) GetKey(z *store_fs.External) string {
	if z == nil {
		return ""
	}

	sh := z.GetBlobSha()

	if sh.IsNull() {
		return ""
	}

	return sh.String()
}

func MakeMutableSetUniqueBlob(zs ...*store_fs.External) MutableSet {
	return collections_value.MakeMutableValueSet[*store_fs.External](
		KeyerBlob{},
		zs...)
}
