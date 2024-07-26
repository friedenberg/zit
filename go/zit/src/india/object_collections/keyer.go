package object_collections

import "code.linenisgreat.com/zit/go/zit/src/hotel/sku"

type KeyerFD struct{}

func (k KeyerFD) GetKey(z sku.ExternalLike) string {
	if z == nil {
		return ""
	}

	return z.String()
}

type KeyerStored struct{}

func (k KeyerStored) GetKey(el sku.ExternalLike) string {
	if el == nil {
		return ""
	}

	z := el.GetSku()

	if z.GetObjectSha().IsNull() {
		return ""
	}

	return z.GetObjectSha().String()
}

type KeyerBlob struct{}

func (k KeyerBlob) GetKey(el sku.ExternalLike) string {
	if el == nil {
		return ""
	}

	z := el.GetSku()

	sh := z.GetBlobSha()

	if sh.IsNull() {
		return ""
	}

	return sh.String()
}
