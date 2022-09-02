package hinweis

import (
	"crypto/sha256"
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type inner struct {
	Left, Right string
}

func (h inner) Kopf() string {
	return h.Left
}

func (h inner) Schwanz() string {
	return h.Right
}

func (h inner) String() string {
	return fmt.Sprintf("%s/%s", h.Left, h.Right)
}

func (h inner) Sha() sha.Sha {
	hash := sha256.New()
	sr := strings.NewReader(h.String())

	if _, err := io.Copy(hash, sr); err != nil {
		errors.PanicIfError(err)
	}

	return sha.FromHash(hash)
}
