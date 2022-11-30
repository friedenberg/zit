package zettel_external

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/echo/fd"
	"github.com/friedenberg/zit/src/india/zettel"
)

type Zettel struct {
	Objekte  zettel.Zettel
	Kennung  hinweis.Hinweis
	Sha      sha.Sha
	ZettelFD fd.FD
	AkteFD   fd.FD
}

type ExternalFormat interface {
	ReadExternalZettelFrom(*Zettel, io.Reader) (int64, error)
	WriteExternalZettelTo(Zettel, io.Writer) (int64, error)
}

func (e Zettel) String() string {
	return e.ExternalPathAndSha()
}

func (e Zettel) ExternalPathAndSha() string {
	if !e.ZettelFD.IsEmpty() {
		return fmt.Sprintf("[%s %s]", e.ZettelFD.Path, e.Sha)
	} else if !e.AkteFD.IsEmpty() {
		return fmt.Sprintf("[%s %s]", e.AkteFD.Path, e.Objekte.Akte)
	} else {
		return ""
	}
}
