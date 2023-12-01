package sha

import (
	"crypto/sha256"
	"fmt"
	"hash"
	"io"
	"path/filepath"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func MakeHashWriter() (h hash.Hash) {
	h = sha256.New()
	return
}

func Make(s schnittstellen.ShaLike) Sha {
	switch st := s.(type) {
	case Sha:
		return st

	case *Sha:
		return *st

	default:
		panic(fmt.Sprintf("wrong type: %T", st))
	}
}

func Must(v string) (s Sha) {
	s = Sha{}

	errors.PanicIfError(s.Set(v))

	return
}

func MakeSha(v string) (s Sha, err error) {
	s = Sha{}

	if err = s.Set(v); err != nil {
		err = errors.Wrap(err)
	}

	return
}

func MakeShaFromPath(p string) (s Sha, err error) {
	schwanz := filepath.Base(p)
	kopf := filepath.Base(filepath.Dir(p))

	switch {
	case schwanz == string(filepath.Separator) || kopf == string(filepath.Separator):
		fallthrough

	case schwanz == "." || kopf == ".":
		err = errors.Errorf(
			"path cannot be turned into a kopf/schwanz pair: '%s/%s'",
			kopf,
			schwanz,
		)

		return
	}

	if s, err = MakeSha(fmt.Sprintf("%s%s", kopf, schwanz)); err != nil {
		err = errors.Wrapf(err, "kopf: %q, schwanz: %q", kopf, schwanz)
		return
	}

	return
}

func FromFormatString(f string, vs ...interface{}) Sha {
	return FromString(fmt.Sprintf(f, vs...))
}

func FromString(s string) Sha {
	hash := sha256.New()
	sr := strings.NewReader(s)

	if _, err := io.Copy(hash, sr); err != nil {
		errors.PanicIfError(err)
	}

	return FromHash(hash)
}

func FromStringer(v schnittstellen.Stringer) Sha {
	return FromString(v.String())
}

func FromHash(h hash.Hash) (s Sha) {
	s = Sha{}
	s.Reset()

	b := h.Sum(s.data.AvailableBuffer())
	s.data.Write(b)

	return
}
