package value

import (
	"bytes"
	"crypto/sha256"
	"io"
	"strings"
)

type Value struct {
	value string
	sha   _Sha
}

func (v Value) Sha() _Sha {
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
		err = _Error(err)
		return
	}

	v.sha = _MakeShaFromHash(hash)
	v.value = s

	return
}

func (v Value) Buffer() *bytes.Buffer {
	return bytes.NewBufferString(v.String())
}
