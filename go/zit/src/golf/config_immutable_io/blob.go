package config_immutable_io

import (
	"encoding/gob"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/delta/config_immutable"
)

type blobV1Coder struct{}

func (blobV1Coder) DecodeFrom(
	subject typeWithConfigLoaded,
	r io.Reader,
) (n int64, err error) {
	subject.Object.ImmutableConfig = &config_immutable.TomlV1{}
	td := toml.NewDecoder(r)

	if err = td.Decode(subject.Object.ImmutableConfig); err != nil {
		if err == io.EOF {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (blobV1Coder) EncodeTo(
	subject typeWithConfigLoaded,
	w io.Writer,
) (n int64, err error) {
	te := toml.NewEncoder(w)

	if err = te.Encode(subject.Object.ImmutableConfig); err != nil {
		if err == io.EOF {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

type blobV0Coder struct{}

func (blobV0Coder) DecodeFrom(
	subject typeWithConfigLoaded,
	r io.Reader,
) (n int64, err error) {
	subject.Object.ImmutableConfig = &config_immutable.V0{}

	dec := gob.NewDecoder(r)

	if err = dec.Decode(subject.Object.ImmutableConfig); err != nil {
		if err == io.EOF {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (blobV0Coder) EncodeTo(
	subject typeWithConfigLoaded,
	w io.Writer,
) (n int64, err error) {
	dec := gob.NewEncoder(w)

	if err = dec.Encode(subject.Object.ImmutableConfig); err != nil {
		if err == io.EOF {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
