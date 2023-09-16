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

type lessor struct{}

func (_ lessor) Less(a, b wrapper) bool {
	return a.SkuLikePtr.GetTai().Less(b.SkuLikePtr.GetTai())
}

func (_ lessor) LessPtr(a, b *wrapper) bool {
	return a.SkuLikePtr.GetTai().Less(b.SkuLikePtr.GetTai())
}

type equaler struct{}

func (_ equaler) Equals(a, b wrapper) bool {
	return a.SkuLikePtr.EqualsSkuLike(b.SkuLikePtr)
}

func (_ equaler) EqualsPtr(a, b *wrapper) bool {
	return a.SkuLikePtr.EqualsSkuLike(b.SkuLikePtr)
}

type resetter struct{}

func (_ resetter) Reset(a *wrapper) {
	a.SkuLikePtr.Reset()
}

func (_ resetter) ResetWith(a *wrapper, b wrapper) {
	a.SkuLikePtr = b.SkuLikePtr
}

func (_ resetter) ResetWithPtr(a *wrapper, b *wrapper) {
	a.SkuLikePtr = b.SkuLikePtr
}
