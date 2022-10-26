package zettel_verzeichnisse

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/konfig"
)

// TODO add efficient parsing of hiding tags
func MakeWriterKonfig(k konfig.Konfig) Writer {
	if k.IncludeHidden {
		return WriterNoop{}
	}

	return MakeWriter(
		func(z *Zettel) (err error) {
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
		},
	)
}
