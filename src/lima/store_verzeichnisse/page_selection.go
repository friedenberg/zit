package store_verzeichnisse

import (
	"crypto/sha256"
	"io"
	"strconv"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/sha"
	"github.com/friedenberg/zit/src/foxtrot/hinweis"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/kilo/zettel"
)

func (i Zettelen) PageForHinweis(h hinweis.Hinweis) (n int, err error) {
	s := h.GetSha()
	return i.PageForSha(s)
}

func (i Zettelen) PageForTransacted(z zettel.Transacted) (n int, err error) {
	s := z.Sku.ObjekteSha
	return i.PageForSha(s)
}

func (i Zettelen) PageForEtikett(e kennung.Etikett) (n int, err error) {
	errors.TodoP4("does this actually work")
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
	ss := s.String()[:DigitWidth]

	if n1, err = strconv.ParseInt(ss, 16, 64); err != nil {
		err = errors.Wrap(err)
		return
	}

	n = int(n1)

	return
}
