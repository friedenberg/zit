package value

import (
	"bytes"
	"crypto/sha256"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/bravo/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type Value struct {
	value string
	sha   sha.Sha
}

func (v Value) Sha() sha.Sha {
	return v.sha
}

func (v Value) String() string {
	return v.value
}

func (v *Value) SetString(s string) (err error) {
	s = strings.TrimSpace(s)

	hash := sha256.New()
	sr := strings.NewReader(s)

	if _, err = io.Copy(hash, sr); err != nil {
		err = errors.Error(err)
		return
	}

	v.sha = sha.FromHash(hash)
	v.value = s

	return
}

func (v Value) Buffer() *bytes.Buffer {
	return bytes.NewBufferString(v.String())
}
