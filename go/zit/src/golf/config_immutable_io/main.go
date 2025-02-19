package config_immutable_io

import (
	"bytes"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/echo/triple_hyphen_io"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
)

type typeWithConfigLoaded = *triple_hyphen_io.TypeWithObject[*ConfigLoaded]

var typedCoders = map[string]interfaces.Coder[typeWithConfigLoaded]{
	builtin_types.ImmutableConfigV1: blobV1Coder{},
	"":                              blobV0Coder{},
}

var coder = triple_hyphen_io.Coder[typeWithConfigLoaded]{
	Metadata: triple_hyphen_io.TypedMetadataCoder[*ConfigLoaded]{},
	Blob:     triple_hyphen_io.TypedCoders[*ConfigLoaded](typedCoders),
}

type Coder struct{}

func (coder Coder) DecodeFromFile(
	object *ConfigLoaded,
	p string,
) (err error) {
	var r io.Reader

	{
		var f *os.File

		if f, err = files.OpenExclusiveReadOnly(p); err != nil {
			if errors.IsNotExist(err) {
				err = nil
				r = bytes.NewBuffer(nil)
			} else {
				err = errors.Wrap(err)
				return
			}
		} else {
			defer errors.DeferredCloser(&err, f)

			r = f
		}
	}

	if _, err = coder.DecodeFrom(object, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (Coder) DecodeFrom(
	subject *ConfigLoaded,
	reader io.Reader,
) (n int64, err error) {
	if n, err = coder.DecodeFrom(
		&triple_hyphen_io.TypeWithObject[*ConfigLoaded]{
			Type:   &subject.Type,
			Object: subject,
		},
		reader,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	subject.BlobStoreImmutableConfig = env_dir.MakeConfigFromImmutableBlobConfig(
		subject.ImmutableConfig.GetBlobStoreConfigImmutable(),
	)

	return
}

func (Coder) EncodeTo(
	subject *ConfigLoaded,
	writer io.Writer,
) (n int64, err error) {
	if n, err = coder.EncodeTo(
		&triple_hyphen_io.TypeWithObject[*ConfigLoaded]{
			Type:   &subject.Type,
			Object: subject,
		},
		writer,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
