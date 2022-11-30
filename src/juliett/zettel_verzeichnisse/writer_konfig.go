package zettel_verzeichnisse

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/echo/konfig"
)

// TODO add efficient parsing of hiding tags
func MakeWriterKonfig(k konfig.Konfig) collections.WriterFunc[*Zettel] {
	if k.IncludeHidden {
		return collections.MakeWriterNoop[*Zettel]()
	}

	return func(z *Zettel) (err error) {
		for _, p := range z.EtikettenSorted {
			for _, t := range k.Transacted.Objekte.Akte.EtikettenHidden {
				if strings.HasPrefix(p, t) {
					errors.Log().Printf("eliding: %s", z.Transacted.Sku.Kennung)
					err = io.EOF
					return
				}
			}
		}

		return
	}
}
