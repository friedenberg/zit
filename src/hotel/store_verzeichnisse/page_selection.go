package store_verzeichnisse

import (
	"crypto/sha256"
	"io"
	"strconv"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/etikett"
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

func (i Zettelen) PageForEtikett(e etikett.Etikett) (n int, err error) {
  //TODO does this actually work?
	return i.PageForRune(rune(e.String()[0]))
}

func (i Zettelen) PageForRune(r rune) (n int, err error) {
	return i.PageForString(string([]rune{r}))
}

func (i Zettelen) PageForString(s string) (n int, err error) {
	sr := strings.NewReader(s)
	hash := sha256.New()

	if _, err = io.Copy(hash, sr); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := sha.FromHash(hash)
	return i.PageForSha(sh)
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
