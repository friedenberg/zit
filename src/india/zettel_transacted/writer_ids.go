package zettel_transacted

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/id_set"
)

type WriterIds struct {
	Filter id_set.Filter
}

func (w WriterIds) WriteZettelTransacted(z *Transacted) (err error) {
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
	return z.Named.Kennung
}

func (z zettelFilterable) AkteEtiketten() kennung.EtikettSet {
	return z.Named.Stored.Objekte.Etiketten
}

func (z zettelFilterable) AkteSha() sha.Sha {
	return z.Named.Stored.Objekte.Akte
}

func (z zettelFilterable) SetAkteSha(v sha.Sha) {
	z.Named.Stored.Objekte.Akte = v
}

func (z zettelFilterable) ObjekteSha() sha.Sha {
	return z.Named.Stored.Sha
}

func (z zettelFilterable) SetObjekteSha(
	arf gattung.AkteReaderFactory,
	v string,
) (err error) {
	if err = z.Named.Stored.Sha.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (z zettelFilterable) AkteTyp() kennung.Typ {
	return z.Named.Stored.Objekte.Typ
}
