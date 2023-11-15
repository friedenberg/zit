package store_verzeichnisse

import (
	"crypto/sha256"
	"io"
	"strconv"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/echo/kennung"
)

func (i Store) PageForKennung(h kennung.Kennung) (n int, err error) {
	s := sha.FromStringer(h)
	return i.PageForSha(s)
}

func (i Store) PageForString(s string) (n int, err error) {
	sr := strings.NewReader(s)
	hash := sha256.New()

	if _, err = io.Copy(hash, sr); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := sha.FromHash(hash)
	return i.PageForSha(sh)
}

func (i Store) PageForSha(s schnittstellen.ShaLike) (n int, err error) {
	var n1 int64
	ss := s.String()[:DigitWidth]

	if n1, err = strconv.ParseInt(ss, 16, 64); err != nil {
		err = errors.Wrap(err)
		return
	}

	n = int(n1)

	return
}
