package sku

type wrapper struct {
	SkuLikePtr
}

func (a *wrapper) Reset() {
	a.SkuLikePtr.Reset()
}

func (a *wrapper) ResetWith(b wrapper) {
	a.SkuLikePtr = b.SkuLikePtr
}

type skuLessor struct{}

func (_ skuLessor) Less(a, b wrapper) bool {
	return a.SkuLikePtr.GetTai().Less(b.SkuLikePtr.GetTai())
}

func (_ skuLessor) LessPtr(a, b *wrapper) bool {
	return a.SkuLikePtr.GetTai().Less(b.SkuLikePtr.GetTai())
}

type skuEqualer struct{}

func (_ skuEqualer) Equals(a, b wrapper) bool {
	return a.SkuLikePtr.EqualsSkuLike(b.SkuLikePtr)
}

func (_ skuEqualer) EqualsPtr(a, b *wrapper) bool {
	return a.SkuLikePtr.EqualsSkuLike(b.SkuLikePtr)
}

type skuResetter struct{}

func (_ skuResetter) Reset(a *wrapper) {
	a.SkuLikePtr.Reset()
}

func (_ skuResetter) ResetWith(a *wrapper, b wrapper) {
	a.SkuLikePtr = b.SkuLikePtr
}

func (_ skuResetter) ResetWithPtr(a *wrapper, b *wrapper) {
	a.SkuLikePtr = b.SkuLikePtr
}
