package store_verzeichnisse

import (
	"strconv"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

func (i Zettelen) PageForHinweis(h hinweis.Hinweis) (n int, err error) {
	s := h.Sha()
	return i.PageForSha(s)
}

func (i Zettelen) PageForTransacted(z zettel_transacted.Zettel) (n int, err error) {
	s := z.Named.Stored.Sha
	return i.PageForSha(s)
}

func (i Zettelen) PageForSha(s sha.Sha) (n int, err error) {
	var n1 int64
	ss := s.String()[:digitWidth]

	if n1, err = strconv.ParseInt(ss, 16, 64); err != nil {
		err = errors.Wrap(err)
		return
	}

	n = int(n1)

	return
}
