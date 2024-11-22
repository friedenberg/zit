package sku

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type ObjectIdKeyer[T ids.ObjectIdGetter] struct{}

func (sk ObjectIdKeyer[T]) GetKey(e T) string {
	return e.GetObjectId().String()
}

type ExternalObjectIdKeyer[T ExternalObjectIdGetter] struct{}

func (ExternalObjectIdKeyer[T]) GetKey(el T) string {
	return el.GetExternalObjectId().String()
}

type DescriptionKeyer[T ExternalLikeGetter] struct{}

func (DescriptionKeyer[T]) GetKey(el T) string {
	return el.GetSkuExternal().Metadata.Description.String()
}
