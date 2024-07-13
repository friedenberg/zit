package age

import (
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"filippo.io/age"
)

type identity interface {
	age.Identity
	interfaces.Stringer
}

type Identity struct {
	identity
	age.Recipient

	path     string
	disabled bool
}

func (i *Identity) IsDisabled() bool {
	return i.disabled
}

func (i *Identity) IsEmpty() bool {
	return i.identity == nil
}

func (i *Identity) String() string {
	return i.path
}

func (i *Identity) SetFromX25519Identity(identity string) (err error) {
	var x *X25519Identity

	if x, err = age.ParseX25519Identity(identity); err != nil {
		err = errors.Wrapf(err, "Identity: %s", identity)
		return
	}

	i.SetX25519Identity(x)

	return
}

func (i *Identity) SetX25519Identity(x *age.X25519Identity) {
	i.disabled = false
	i.identity = x
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
	if i.IsDisabled() || !i.IsEmpty() {
		return
	}

	var x *X25519Identity

	if x, err = age.GenerateX25519Identity(); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.SetX25519Identity(x)

	return
}

func (i *Identity) WriteToPathIfNecessary(basePath string) (err error) {
	if i.IsDisabled() || i.IsEmpty() || basePath == "" {
		return
	}

	var f *os.File

	if f, err = files.CreateExclusiveWriteOnly(basePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	if _, err = io.WriteString(f, i.identity.String()); err != nil {
		err = errors.Wrap(err)
		return
	}

	i.path = basePath

	return
}
