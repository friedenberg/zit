package zettel_verzeichnisse

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/juliett/konfig_compiled"
)

// TODO add efficient parsing of hiding tags
func MakeWriterKonfig(k konfig_compiled.Compiled) collections.WriterFunc[*Verzeichnisse] {
	if k.IncludeHidden {
		return collections.MakeWriterNoop[*Verzeichnisse]()
	}

	return func(z *Verzeichnisse) (err error) {
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
