package zettel_external

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/foxtrot/zettel"
)

type Zettel struct {
	Named    zettel.Named
	ZettelFD FD
	AkteFD   FD
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
		return fmt.Sprintf("[%s %s]", e.ZettelFD.Path, e.Named.Stored.Sha)
	} else if !e.AkteFD.IsEmpty() {
		return fmt.Sprintf("[%s %s]", e.AkteFD.Path, e.Named.Stored.Objekte.Akte)
	} else {
		return ""
	}
}
