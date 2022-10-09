package zettel_verzeichnisse

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/charlie/konfig"
)

//TODO add efficient parsing of hiding tags
func MakeWriterKonfig(k konfig.Konfig) Writer {
	return MakeWriter(
		func(z *Zettel) (err error) {
			if k.IncludeHidden {
				return
			}

			for _, p := range z.EtikettenExpandedSorted {
				for tn, tv := range k.Tags {
					if !tv.Hide {
						continue
					}

					if strings.HasPrefix(p, tn) {
						err = io.EOF
						return
					}
				}
			}

			return
		},
	)
}
