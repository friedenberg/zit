package age

import (
	"io"

	"filippo.io/age"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
)

type Age struct {
	recipients []Recipient
	identities []Identity
}

func MakeFromIdentity(path_or_identity string) (a *Age, err error) {
	switch {
	case path_or_identity == "":
		a = &Age{}

	case files.Exists(path_or_identity):
		if a, err = MakeFromIdentityFile(path_or_identity); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		if a, err = MakeFromIdentityString(path_or_identity); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func MakeFromIdentityFile(basePath string) (a *Age, err error) {
	var contents string

	if contents, err = files.ReadAllString(basePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	if a, err = MakeFromIdentityString(contents); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeFromIdentityString(contents string) (a *Age, err error) {
	var i *X25519Identity

	if i, err = age.ParseX25519Identity(contents); err != nil {
		err = errors.Wrap(err)
		return
	}

	a = &Age{
		recipients: []Recipient{i.Recipient()},
		identities: []Identity{i},
	}

	return
}

// func (a *Age) AddBech32PivYubikeyEC256(bech string) (err error) {
// 	var r *age.PivYubikeyEC256Recipient

// 	if r, err = age.ParseBech32PivYubikeyEC256Recipient(bech); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	var i *age.PivYubikeyEC256Identity

// 	if i, err = age.ReadPivYubikeyEC256Identity(r); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	a.recipients = append(a.recipients, r)
// 	a.identities = append(a.identities, i)

// 	return
// }

func (a *Age) Recipients() []Recipient {
  if a == nil {
    return nil
  }

	return a.recipients
}

func (a *Age) Identities() []Identity {
  if a == nil {
    return nil
  }

	return a.identities
}

func (a *Age) Decrypt(src io.Reader) (out io.Reader, err error) {
	i := a.Identities()

	switch len(i) {
	case 0:
		out = src
		return

	default:
		if out, err = age.Decrypt(src, a.Identities()...); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (a *Age) Encrypt(dst io.Writer) (out io.WriteCloser, err error) {
	r := a.Recipients()

	switch len(r) {
	case 0:
		out = writeCloser{dst}
		return

	default:
		if out, err = age.Encrypt(dst, r...); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (a *Age) Close() error {
	if a == nil {
		return nil
	}

	for _, i := range a.identities {
		if c, ok := i.(io.Closer); ok {
			err := c.Close()
			if err != nil {
				return err
			}
		}
	}

	return nil
}
