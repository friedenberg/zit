package sku

type wrapper struct {
	SkuLikePtr
}

func (a wrapper) Less(b wrapper) bool {
	return a.SkuLikePtr.GetTai().Less(b.SkuLikePtr.GetTai())
}

func (a wrapper) Equals(b wrapper) bool {
	return a.SkuLikePtr.EqualsSkuLike(b.SkuLikePtr)
}

func (a *wrapper) Reset() {
	a.SkuLikePtr.Reset()
}

func (a *wrapper) ResetWith(b wrapper) {
	a.SkuLikePtr = b.SkuLikePtr
}
