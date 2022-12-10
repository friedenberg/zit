package zettel_verzeichnisse

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/hotel/konfig_compiled"
)

// TODO add efficient parsing of hiding tags
func MakeWriterKonfig(k konfig_compiled.Compiled) collections.WriterFunc[*Zettel] {
	if k.IncludeHidden {
		return collections.MakeWriterNoop[*Zettel]()
	}

	return func(z *Zettel) (err error) {
		for _, p := range z.EtikettenSorted {
			for _, t := range k.EtikettenHidden {
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
