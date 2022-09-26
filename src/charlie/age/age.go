package age

import (
	"io"

	"filippo.io/age"
	"github.com/friedenberg/zit/src/bravo/files"
)

type Age struct {
	recipients []Recipient
	identities []Identity
}

func Make(basePath string) (a *Age, err error) {
	var contents string

	if contents, err = files.ReadAllString(basePath); err != nil {
		return
	}

	var i *X25519Identity

	if i, err = age.ParseX25519Identity(contents); err != nil {
		return
	}

	a = &Age{
		recipients: []Recipient{i.Recipient()},
		identities: []Identity{i},
	}

	return
}

func (a Age) Recipients() []Recipient {
	return a.recipients
}

func (a Age) Identities() []Identity {
	return a.identities
}

func (a Age) Decrypt(src io.Reader) (io.Reader, error) {
	i := a.Identities()

	switch len(i) {
	case 0:
		return src, nil

	default:
		return age.Decrypt(src, a.Identities()...)
	}
}

func (a Age) Encrypt(dst io.Writer) (io.WriteCloser, error) {
	r := a.Recipients()

	switch len(r) {
	case 0:
		return writeCloser{dst}, nil

	default:
		return age.Encrypt(dst, r...)
	}
}
