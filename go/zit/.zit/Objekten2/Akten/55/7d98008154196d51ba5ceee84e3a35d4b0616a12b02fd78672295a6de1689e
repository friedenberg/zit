package age

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"filippo.io/age"
)

type (
	Recipient       = age.Recipient
	X25519Identity  = age.X25519Identity
	X25519Recipient = age.X25519Recipient
)

type Age struct {
	recipients []Recipient
	identities []age.Identity
}

func (a *Age) AddIdentity(
	identity Identity,
) (err error) {
	if identity.IsDisabled() || identity.IsEmpty() {
		return
	}

	a.recipients = append(a.recipients, identity)
	a.identities = append(a.identities, identity)

	return
}

func (a *Age) AddIdentityOrGenerateIfNecessary(
	identity Identity,
	path string,
) (err error) {
	if err = identity.GenerateIfNecessary(path); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = identity.WriteToPathIfNecessary(path); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = a.AddIdentity(identity); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func MakeFromIdentity(identity Identity, path string) (a *Age, err error) {
	a = &Age{}
	err = a.AddIdentityOrGenerateIfNecessary(identity, path)
	return
}

func MakeFromIdentityPathOrString(path_or_identity string) (a *Age, err error) {
	var i Identity

	if err = i.Set(path_or_identity); err != nil {
		err = errors.Wrap(err)
		return
	}

	return MakeFromIdentity(i, "")
}

func MakeFromIdentityFile(basePath string) (a *Age, err error) {
	var i Identity

	if err = i.SetFromPath(basePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	return MakeFromIdentity(i, "")
}

func MakeFromIdentityString(contents string) (a *Age, err error) {
	var i Identity

	if err = i.SetFromX25519Identity(contents); err != nil {
		err = errors.Wrap(err)
		return
	}

	return MakeFromIdentity(i, "")
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

func (a *Age) Identities() []age.Identity {
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
