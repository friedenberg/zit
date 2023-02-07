package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/objekte"
)

func init() {
	errors.TodoP2("move away from this and replace with compiled filter")
}

type WriterIds struct {
	Filter kennung.Filter
}

func (w WriterIds) WriteTransactedLike(maybeZ objekte.TransactedLike) (err error) {
	if z, ok := maybeZ.(*Transacted); ok {
		return w.WriteZettelTransacted(z)
	}

	return
}

func (w WriterIds) WriteZettelTransacted(z *Transacted) (err error) {
	z1 := zettelFilterable{Transacted: z}
	return w.Filter.Include(z1)
}

type zettelFilterable struct {
	*Transacted
}

func (z zettelFilterable) GetGattung() schnittstellen.Gattung {
	return gattung.Zettel
}

func (z zettelFilterable) Hinweis() kennung.Hinweis {
	return z.Sku.Kennung
}

func (z zettelFilterable) AkteEtiketten() kennung.EtikettSet {
	return z.Objekte.Etiketten
}

func (z zettelFilterable) AkteSha() sha.Sha {
	return z.Objekte.Akte
}

func (z zettelFilterable) SetAkteSha(v sha.Sha) {
	z.Objekte.Akte = v
}

func (z zettelFilterable) ObjekteSha() sha.Sha {
	return z.Sku.ObjekteSha
}

func (z zettelFilterable) SetObjekteSha(
	arf schnittstellen.AkteReaderFactory,
	v string,
) (err error) {
	if err = z.Sku.ObjekteSha.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (z zettelFilterable) AkteTyp() kennung.Typ {
	return z.Objekte.Typ
}
