package objekte

import (
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

type ExternalKeyer[
	T Akte[T],
	T1 AktePtr[T],
	T2 kennung.KennungLike[T2],
	T3 kennung.KennungLikePtr[T2],
] struct{}

func (_ ExternalKeyer[T, T1, T2, T3]) Key(e *sku.External[T2, T3]) string {
	if e == nil {
		return ""
	}

	return e.GetKennung().String()
}
