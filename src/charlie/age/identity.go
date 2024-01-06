package age

import (
	"io"
	"os"

	"filippo.io/age"
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
)

type Identity struct {
	age.Identity
	age.Recipient

	path     string
	disabled bool
}

func (i *Identity) IsDisabled() bool {
	return i.disabled
}

func (i *Identity) IsEmpty() bool {
	return i.Identity == nil
}

func (i *Identity) String() string {
	return i.path
}

func (i *Identity) SetFromX25519Identity(identity string) (err error) {
	var x *X25519Identity

	if x, err = age.ParseX25519Identity(identity); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.SetX25519Identity(x)

	return
}

func (i *Identity) SetX25519Identity(x *age.X25519Identity) {
	i.disabled = false
	i.Identity = x
	i.Recipient = x.Recipient()
}

func (i *Identity) SetFromPath(path string) (err error) {
	i.path = path

	var contents string

	if contents, err = files.ReadAllString(path); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = i.SetFromX25519Identity(contents); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (i *Identity) Set(path_or_identity string) (err error) {
	switch {
	case path_or_identity == "":

	case path_or_identity == "disabled" || path_or_identity == "none":
		i.disabled = true
		// no-op

	case files.Exists(path_or_identity):
		if err = i.SetFromPath(path_or_identity); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		if err = i.SetFromX25519Identity(path_or_identity); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (i *Identity) GenerateIfNecessary(basePath string) (err error) {
	if i.disabled || basePath == "" {
		return
	}

	var x *X25519Identity

	if x, err = age.GenerateX25519Identity(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var f *os.File

	if f, err = files.CreateExclusiveWriteOnly(basePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if _, err = io.WriteString(f, x.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.SetX25519Identity(x)
	i.path = basePath

	return
}
