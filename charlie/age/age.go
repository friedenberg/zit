package age

import (
	"io"
)

type Age interface {
	Recipient() _AgeRecipient
	Identity() _AgeIdentity
	Decrypt(src io.Reader) (io.Reader, error)
	Encrypt(dst io.Writer) (io.WriteCloser, error)
}

type age struct {
	recipient _AgeRecipient
	identity  _AgeIdentity
}

func Make(basePath string) (a *age, err error) {
	var contents string

	if contents, err = _ReadStringAll(basePath); err != nil {
		return
	}

	var i *_AgeX25519Identity

	if i, err = _ParseX25519Identity(contents); err != nil {
		return
	}

	a = &age{
		recipient: i.Recipient(),
		identity:  i,
	}

	return
}

func (a age) Recipient() _AgeRecipient {
	return a.recipient
}

func (a age) Identity() _AgeIdentity {
	return a.identity
}

func (a age) Decrypt(src io.Reader) (io.Reader, error) {
	return _AgeDecrypt(src, a.Identity())
}

func (a age) Encrypt(dst io.Writer) (io.WriteCloser, error) {
	return _AgeEncrypt(dst, a.Recipient())
}
