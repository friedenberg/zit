package objekte

import (
	"encoding/base64"
	"io"
	"strings"
)

func EncodeBase64(age _Age, in string) (out string, sha _Sha, err error) {
	r1 := strings.NewReader(in)
	sb := &strings.Builder{}

	r3 := base64.NewEncoder(base64.StdEncoding, sb)

	var r2 *writer

	if r2, err = NewWriter(age, r3); err != nil {
		err = _Error(err)
		return
	}

	if _, err = io.Copy(r2, r1); err != nil {
		err = _Error(err)
		return
	}

	r2.Close()
	r3.Close()

	sha = r2.Sha()
	out = sb.String()

	return
}

func DecodeBase64(age _Age, in string) (out string, sha _Sha, err error) {
	r1 := strings.NewReader(in)
	r2 := base64.NewDecoder(base64.StdEncoding, r1)

	var r3 *reader

	if r3, err = NewReader(age, r2); err != nil {
		err = _Error(err)
		return
	}

	sb := &strings.Builder{}

	if _, err = io.Copy(sb, r3); err != nil {
		err = _Error(err)
		return
	}

	r3.Close()

	sha = r3.Sha()
	out = sb.String()

	return
}
