package sku

import (
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type ObjectIdKeyer[
	T ids.ObjectIdGetter,
] struct{}

func (sk ObjectIdKeyer[T]) GetKey(e T) string {
	return e.GetObjectId().String()
}

type ExternalObjectIdKeyer[
	T interface {
		ids.ObjectIdGetter
		ExternalObjectIdGetter
		TransactedGetter
		ExternalLikeGetter
	},
] struct{}

func (ExternalObjectIdKeyer[T]) GetKey(el T) string {
	eoid := el.GetExternalObjectId()

	if !eoid.IsEmpty() {
		return eoid.String()
	}

	oid := el.GetSkuExternal().GetObjectId()

	if !oid.IsEmpty() {
		return oid.String()
	}

	desc := el.GetSkuExternal().Metadata.Description.String()

	if desc != "" {
		return desc
	}

	panic(fmt.Sprintf("empty key for external like: %#v", el))
}
