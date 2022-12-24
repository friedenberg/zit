package zettel

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/golf/id_set"
)

// TODO-P2 move away from this and replace with compiled filter
type WriterIds struct {
	Filter id_set.Filter
}

func (w WriterIds) WriteZettelVerzeichnisse(z *Transacted) (err error) {
	z1 := zettelFilterable{Transacted: z}
	return w.Filter.Include(z1)
}

type zettelFilterable struct {
	*Transacted
}

func (z zettelFilterable) Gattung() gattung.Gattung {
	return gattung.Zettel
}

func (z zettelFilterable) Hinweis() hinweis.Hinweis {
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
	return z.Sku.Sha
}

func (z zettelFilterable) SetObjekteSha(
	arf gattung.AkteReaderFactory,
	v string,
) (err error) {
	if err = z.Sku.Sha.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (z zettelFilterable) AkteTyp() kennung.Typ {
	return z.Objekte.Typ
}
