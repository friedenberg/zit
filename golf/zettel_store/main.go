package zettel_store

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/delta/objekte"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
	"github.com/friedenberg/zit/hotel/zettels"
)

type ZettelStore struct {
	zettels.Zettels
}

func (s ZettelStore) writeZettelObjekte(z zettel.Zettel) (err error) {
	var w *objekte.Mover

	if w, err = objekte.NewWriterMover(s.Age(), s.Umwelt().DirZettelen()); err != nil {
		err = errors.Error(err)
		return
	}

	defer w.Close()

	f := zettel_formats.Objekte{}

	if _, err = f.WriteTo(z, w); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s ZettelStore) Create(in zettel.Zettel) (z stored_zettel.Named, err error) {
	if z, err = s.Zettels.Create(in); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.writeZettelObjekte(in); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s ZettelStore) CreateWithHinweis(in zettel.Zettel, h hinweis.Hinweis) (z stored_zettel.Named, err error) {
	if z, err = s.Zettels.CreateWithHinweis(in, h); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.writeZettelObjekte(in); err != nil {
		err = errors.Error(err)
		return
	}

	return
}

func (s ZettelStore) Update(z stored_zettel.Named) (stored stored_zettel.Named, err error) {
	if stored, err = s.Zettels.Update(z); err != nil {
		err = errors.Error(err)
		return
	}

	if err = s.writeZettelObjekte(stored.Zettel); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
