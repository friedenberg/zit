package sku

import (
	"fmt"

	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type External[T kennung.KennungLike[T], T1 kennung.KennungLikePtr[T]] struct {
	Kennung    T
	AkteSha    sha.Sha
	ObjekteSha sha.Sha
}

func (a *External[T, T1]) Transacted() (b Transacted[T, T1]) {
	b = Transacted[T, T1]{
		Kennung:    a.Kennung,
		ObjekteSha: a.ObjekteSha,
		AkteSha:    a.AkteSha,
	}

	return
}

func (a *External[T, T1]) Reset(b *External[T, T1]) {
	if b == nil {
		a.ObjekteSha = sha.Sha{}
		a.AkteSha = sha.Sha{}
		T1(&a.Kennung).Reset(nil)
	} else {
		a.ObjekteSha = b.ObjekteSha
		a.AkteSha = b.AkteSha
		T1(&a.Kennung).Reset(&b.Kennung)
	}
}

func (a External[T, T1]) Equals(b *External[T, T1]) (ok bool) {
	if a.Kennung.Equals(&b.Kennung) {
		return
	}

	if !a.ObjekteSha.Equals(b.ObjekteSha) {
		return
	}

	return true
}

func (o External[T, T1]) GetKey() string {
	return fmt.Sprintf("%s.%s", o.Kennung.GetGattung(), o.Kennung)
}
