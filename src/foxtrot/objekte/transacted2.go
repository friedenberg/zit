package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/metadatei_io"
	"github.com/friedenberg/zit/src/echo/sku"
)

type Transacted2[T Objekte2, T1 ObjektePtr[T], T2 Identifier2[T2], T3 IdentifierPtr[T2]] struct {
	Objekte T
	Sku     sku.Sku2[T2, T3]
}

func (t Transacted2[T, T1, T2, T3]) Kennung() T3 {
	return &t.Sku.Kennung
}

func (t Transacted2[T, T1, T2, T3]) Sha() sha.Sha {
	return t.Sku.Sha
}

func (t Transacted2[T, T1, T2, T3]) AkteSha() sha.Sha {
	return t.Objekte.AkteSha()
}

func (t *Transacted2[T, T1, T2, T3]) SetSha(
	arf metadatei_io.AkteReaderFactory,
	v string,
) (err error) {
	if err = t.Sku.Sha.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	//TODO provide opportunity for Objekte to bootstrap from arf

	return
}

func (t Transacted2[T, T1, T2, T3]) Gattung() gattung.Gattung {
	return t.Sku.Kennung.Gattung()
}

func (zt Transacted2[T, T1, T2, T3]) IsNew() bool {
	return zt.Sku.Kopf == zt.Sku.Schwanz && zt.Sku.TransactionIndex == 0
}

//TODO
func (a Transacted2[T, T1, T2, T3]) Equals(
	b Transacted2[T, T1, T2, T3],
) bool {
	return false
}

func (a *Transacted2[T, T1, T2, T3]) Reset(
	b *Transacted2[T, T1, T2, T3],
) {
	a.Sku.Reset(&b.Sku)
	//TODO
	// a.Kennung.Reset(b.Kennung)
	// a.Objekte.Reset(b.Objekte)
}