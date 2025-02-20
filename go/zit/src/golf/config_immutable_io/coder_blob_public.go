package config_immutable_io

import (
	"encoding/gob"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/delta/config_immutable"
)

type blobV1CoderPublic struct{}

func (blobV1CoderPublic) DecodeFrom(
	subject typeWithConfigLoadedPublic,
	r io.Reader,
) (n int64, err error) {
	subject.Struct.ImmutableConfig = &config_immutable.TomlV1Public{}
	td := toml.NewDecoder(r)

	if err = td.Decode(subject.Struct.ImmutableConfig); err != nil {
		if err == io.EOF {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (blobV1CoderPublic) EncodeTo(
	subject typeWithConfigLoadedPublic,
	w io.Writer,
) (n int64, err error) {
	te := toml.NewEncoder(w)

	if err = te.Encode(subject.Struct.ImmutableConfig); err != nil {
		if err == io.EOF {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

type blobV0CoderPublic struct{}

func (blobV0CoderPublic) DecodeFrom(
	subject typeWithConfigLoadedPublic,
	r io.Reader,
) (n int64, err error) {
	subject.Struct.ImmutableConfig = &config_immutable.V0Public{}

	dec := gob.NewDecoder(r)

	if err = dec.Decode(subject.Struct.ImmutableConfig); err != nil {
		if err == io.EOF {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (blobV0CoderPublic) EncodeTo(
	subject typeWithConfigLoadedPublic,
	w io.Writer,
) (n int64, err error) {
	dec := gob.NewEncoder(w)

	if err = dec.Encode(subject.Struct.ImmutableConfig); err != nil {
		if err == io.EOF {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
