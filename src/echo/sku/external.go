package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type External[T kennung.KennungLike[T], T1 kennung.KennungLikePtr[T]] struct {
	Kennung T
	Sha     sha.Sha
}

func (a *External[T, T1]) Reset(b *External[T, T1]) {
	if b == nil {
		a.Sha = sha.Sha{}
		T1(&a.Kennung).Reset(nil)
	} else {
		a.Sha = b.Sha
		T1(&a.Kennung).Reset(&b.Kennung)
	}
}

func (a External[T, T1]) Equals(b *External[T, T1]) (ok bool) {
	if a.Kennung.Equals(&b.Kennung) {
		return
	}

	if !a.Sha.Equals(b.Sha) {
		return
	}

	return true
}

func (o External[T, T1]) GetKey() string {
	return fmt.Sprintf("%s.%s", o.Kennung.Gattung(), o.Kennung)
}
