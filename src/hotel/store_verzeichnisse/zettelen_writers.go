package store_verzeichnisse

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

func (i *Zettelen) ZettelWriterSchwanzenOnly() Writer {
	return MakeWriter(
		func(z *Zettel) (err error) {
			if z.PageSelection.Reason != PageSelectionReasonHinweis {
				err = io.EOF
				return
			}

			ok := false

			if ok, err = i.IsSchwanz(z.Transacted); err != nil {
				err = errors.Wrap(err)
				return
			}

			if !ok {
				err = io.EOF
				return
			}

			return
		},
	)
}

//TODO add efficient parsing of hiding tags
func (i *Zettelen) ZettelWriterFilterHidden() Writer {
	return MakeWriter(
		func(z *Zettel) (err error) {
			if i.IncludeHidden {
				return
			}

			for _, p := range z.EtikettenExpandedSorted {
				for tn, tv := range i.Tags {
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
