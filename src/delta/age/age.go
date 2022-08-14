package age

import (
	"io"

	"filippo.io/age"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
)

type Age interface {
	Recipient() Recipient
	Identity() Identity
	Decrypt(src io.Reader) (io.Reader, error)
	Encrypt(dst io.Writer) (io.WriteCloser, error)
}

type ages struct {
	recipient Recipient
	identity  Identity
}

func Make(basePath string) (a *ages, err error) {
	var contents string

	if contents, err = open_file_guard.ReadAllString(basePath); err != nil {
		return
	}

	var i *X25519Identity

	if i, err = age.ParseX25519Identity(contents); err != nil {
		return
	}

	a = &ages{
		recipient: i.Recipient(),
		identity:  i,
	}

	return
}

func (a ages) Recipient() Recipient {
	return a.recipient
}

func (a ages) Identity() Identity {
	return a.identity
}

func (a ages) Decrypt(src io.Reader) (io.Reader, error) {
	return age.Decrypt(src, a.Identity())
}

func (a ages) Encrypt(dst io.Writer) (io.WriteCloser, error) {
	return age.Encrypt(dst, a.Recipient())
}
