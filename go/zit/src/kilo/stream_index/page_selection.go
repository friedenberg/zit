package stream_index

import (
	"crypto/sha256"
	"io"
	"math"
	"strconv"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

func (i *Index) PageForObjectId(h ids.IdLike) (n uint8, err error) {
	s := sha.FromStringer(h)
	return sha.PageIndexForSha(DigitWidth, s)
}

func (i *Index) PageForString(s string) (n uint8, err error) {
	sr := strings.NewReader(s)
	hash := sha256.New()

	if _, err = io.Copy(hash, sr); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh := sha.FromHash(hash)
	return i.PageForSha(sh)
}

func (i *Index) PageForSha(s interfaces.Sha) (n uint8, err error) {
	var n1 int64
	ss := s.String()[:DigitWidth]

	if n1, err = strconv.ParseInt(ss, 16, 64); err != nil {
		err = errors.Wrap(err)
		return
	}

	if n1 > math.MaxUint8 {
		err = errors.Errorf("page out of bounds: %d", n1)
		return
	}

	n = uint8(n1)

	return
}
