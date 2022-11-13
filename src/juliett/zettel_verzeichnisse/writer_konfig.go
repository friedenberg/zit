package zettel_verzeichnisse

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/delta/konfig"
)

// TODO add efficient parsing of hiding tags
func MakeWriterKonfig(k konfig.Konfig) collections.WriterFunc[*Zettel] {
	if k.IncludeHidden {
		return collections.MakeWriterNoop[*Zettel]()
	}

	return func(z *Zettel) (err error) {
		for _, p := range z.EtikettenSorted {
			for _, t := range k.Compiled.EtikettenHidden {
				if strings.HasPrefix(p, t) {
					errors.Log().Printf("eliding: %s", z.Transacted.Named.Hinweis)
					err = io.EOF
					return
				}
			}
		}

		return
	}
}
