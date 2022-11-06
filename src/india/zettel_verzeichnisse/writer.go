package zettel_verzeichnisse

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/foxtrot/zettel_named"
	"github.com/friedenberg/zit/src/hotel/zettel_transacted"
)

type Writer interface {
	WriteZettelVerzeichnisse(z *Zettel) (err error)
}

func MakeWriterZettelTransacted(
	wf collections.WriterFunc[*zettel_transacted.Zettel],
) collections.WriterFunc[*Zettel] {
	return func(z *Zettel) (err error) {
		if err = wf(&z.Transacted); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func MakeWriterZettelNamed(
	wf collections.WriterFunc[*zettel_named.Zettel],
) collections.WriterFunc[*Zettel] {
	return func(z *Zettel) (err error) {
		if err = wf(&z.Transacted.Named); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}
