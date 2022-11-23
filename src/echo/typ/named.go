package typ

import (
	"crypto/sha256"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
)

type Named struct {
	Kennung
	Akte
}

func (t *Named) ObjekteSha() (s sha.Sha, err error) {
	hash := sha256.New()

	enc := MakeEncoderObjekte(hash)

	if _, err = enc.Encode(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	s = sha.FromHash(hash)

	return
}
