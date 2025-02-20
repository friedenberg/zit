package repo_blobs

import (
	"crypto"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/charlie/repo_signing"
	"code.linenisgreat.com/zit/go/zit/src/echo/triple_hyphen_io"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
)

type Blob interface {
	GetRepoBlob() Blob
	GetPublicKey() repo_signing.PublicKey
	// TODO
	// GetSupportedConnectionTypes() []connection_type.Type
}

type BlobMutable interface {
	Blob
	SetPublicKey(crypto.PublicKey)
}

type TypeWithBlob = triple_hyphen_io.TypedStruct[*Blob]

var typedCoders = map[string]interfaces.Coder[*TypeWithBlob]{
	builtin_types.RepoTypeLocalPath:   coderToml[TomlLocalPathV0]{},
	builtin_types.RepoTypeXDGDotenvV0: coderToml[TomlXDGV0]{},
	builtin_types.RepoTypeUri:         coderToml[TomlUriV0]{},
	"":                                coderToml[TomlUriV0]{},
}

var Coder = interfaces.Coder[*TypeWithBlob](triple_hyphen_io.CoderTypeMap[*Blob](typedCoders))

type coderToml[T Blob] struct {
	Blob T
}

func (coder coderToml[T]) DecodeFrom(
	subject *TypeWithBlob,
	reader io.Reader,
) (n int64, err error) {
	decoder := toml.NewDecoder(reader)

	if err = decoder.Decode(&coder.Blob); err != nil {
		if err == io.EOF {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	blob := Blob(coder.Blob)
	subject.Struct = &blob

	return
}

func (coderToml[_]) EncodeTo(
	subject *TypeWithBlob,
	writer io.Writer,
) (n int64, err error) {
	encoder := toml.NewEncoder(writer)

	if err = encoder.Encode(subject.Struct); err != nil {
		if err == io.EOF {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
