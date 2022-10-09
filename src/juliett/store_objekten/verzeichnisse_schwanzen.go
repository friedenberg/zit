package store_objekten

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/konfig"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
	"github.com/friedenberg/zit/src/hotel/zettel_verzeichnisse"
	"github.com/friedenberg/zit/src/india/store_verzeichnisse"
)

type verzeichnisseSchwanzen struct {
	headers [store_verzeichnisse.PageCount]*zettel_verzeichnisse.WriterSchwanzen
	*store_verzeichnisse.Zettelen
	ioFactory
}

func makeVerzeichnisseSchwanzen(
	k konfig.Konfig,
	st standort.Standort,
	iof ioFactory,
	p *zettel_verzeichnisse.Pool,
) (s *verzeichnisseSchwanzen, err error) {
	s = &verzeichnisseSchwanzen{
		ioFactory: iof,
	}

	for i, _ := range s.headers {
		s.headers[i] = zettel_verzeichnisse.MakeWriterSchwanzen()
	}

	s.Zettelen, err = store_verzeichnisse.MakeZettelen(
		k,
		st.DirVerzeichnisseZettelenNeueSchwanzen(),
		s,
		p,
	)

	return
}

func (s *verzeichnisseSchwanzen) Add(tz zettel_transacted.Zettel, v string) (err error) {
	var n int

	if n, err = s.PageForString(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	_ = s.headers[n].WriteZettelTransacted(&tz)

	if err = s.Zettelen.Add(tz, v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *verzeichnisseSchwanzen) ZettelVerzeichnisseWriter(
	n int,
) zettel_verzeichnisse.Writer {
	return s.headers[n]
}
