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
	blob *ConfigLoaded,
	r io.Reader,
) (n int64, err error) {
	blob.ImmutableConfig = &config_immutable.TomlV1{}
	td := toml.NewDecoder(r)

	if err = td.Decode(blob.ImmutableConfig); err != nil {
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
	blob *ConfigLoaded,
	w io.Writer,
) (n int64, err error) {
	te := toml.NewEncoder(w)

	if err = te.Encode(blob.ImmutableConfig); err != nil {
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
	blob *ConfigLoaded,
	r io.Reader,
) (n int64, err error) {
	blob.ImmutableConfig = &config_immutable.V0{}

	dec := gob.NewDecoder(r)

	if err = dec.Decode(blob.ImmutableConfig); err != nil {
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
	blob *ConfigLoaded,
	w io.Writer,
) (n int64, err error) {
	dec := gob.NewEncoder(w)

	if err = dec.Encode(blob.ImmutableConfig); err != nil {
		if err == io.EOF {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
