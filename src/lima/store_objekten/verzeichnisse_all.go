package store_objekten

import (
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/echo/konfig"
	"github.com/friedenberg/zit/src/juliett/zettel_verzeichnisse"
	"github.com/friedenberg/zit/src/kilo/store_verzeichnisse"
)

type verzeichnisseAll struct {
	*store_verzeichnisse.Zettelen
	ioFactory
}

func makeVerzeichnisseAll(
	k konfig.Konfig,
	st standort.Standort,
	iof ioFactory,
	p zettel_verzeichnisse.Pool,
) (s *verzeichnisseAll, err error) {
	s = &verzeichnisseAll{
		ioFactory: iof,
	}

	s.Zettelen, err = store_verzeichnisse.MakeZettelen(
		k,
		st.DirVerzeichnisseZettelenNeue(),
		s,
		p,
	)

	return
}
