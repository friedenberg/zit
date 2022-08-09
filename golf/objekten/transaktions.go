package objekten

import (
	"io"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/id"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/echo/transaktion"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
)

func (s Store) storedZettelFromSha(sh sha.Sha) (sz stored_zettel.Stored, err error) {
	var or io.ReadCloser

	if or, err = s.ReadCloser(id.Path(sh, s.Umwelt().DirObjektenZettelen())); err != nil {
		err = ErrNotFound{Id: sh}
		return
	}

	defer or.Close()

	f := zettel_formats.Objekte{}

	if _, err = f.ReadFrom(&sz.Zettel, or); err != nil {
		err = errors.Error(err)
		return
	}

	sz.Sha = sh

	return
}

func (s Store) transactedZettelFromTransaktionObjekte(t transaktion.Transaktion, o transaktion.Objekte) (tz stored_zettel.Transacted, err error) {
	ok := false

	var h *hinweis.Hinweis

	if h, ok = o.Id.(*hinweis.Hinweis); !ok {
		err = errors.Wrapped(err, "transacktion.Objekte Id was not hinweis but was %s", o.Id)
		return
	}

	tz.Hinweis = *h

	if tz.Stored, err = s.storedZettelFromSha(o.Sha); err != nil {
		err = errors.Wrapped(err, "failed to find zettel objekte for hinweis: %s", tz.Hinweis)
		return
	}

	tz.Time = t.Time

	return
}
