package sku

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

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

type Resetter[T SkuLike, TPtr interface {
	schnittstellen.Ptr[T]
	SkuLikePtr
}] struct{}

func (_ Resetter[T, TPtr]) Reset(a TPtr) {
	a.Reset()
}

func (_ Resetter[T, TPtr]) ResetWith(
	a TPtr,
	b T,
) {
	err := a.SetFromSkuLike(b)
	errors.PanicIfError(err)
}

func (_ Resetter[T, TPtr]) ResetWithPtr(
	a TPtr,
	b TPtr,
) {
	err := a.SetFromSkuLike(b)
	errors.PanicIfError(err)
}
