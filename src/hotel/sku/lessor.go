package sku

import "github.com/friedenberg/zit/src/alfa/schnittstellen"

type Lessor[T SkuLike, TPtr interface {
	schnittstellen.Ptr[T]
	SkuLikePtr
}] struct{}

func (_ Lessor[T, TPtr]) Less(a, b T) bool {
	return a.GetTai().Less(b.GetTai())
}

func (_ Lessor[T, TPtr]) LessPtr(a, b TPtr) bool {
	return a.GetTai().Less(b.GetTai())
}

type Equaler[T SkuLike, TPtr interface {
	schnittstellen.Ptr[T]
	SkuLikePtr
}] struct{}

func (_ Equaler[T, TPtr]) Equals(a, b T) bool {
	return a.EqualsSkuLike(b)
}

func (_ Equaler[T, TPtr]) EqualsPtr(a, b TPtr) bool {
	return a.EqualsSkuLike(b)
}
