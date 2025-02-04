package workspace_config_blobs

import (
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_mutable_blobs"
)

type V0 struct {
	Defaults config_mutable_blobs.DefaultsV1OmitEmpty `toml:"defaults,omitempty"`
	// FileExtensions file_extensions.V1    `toml:"file-extensions"`
	// PrintOptions   options_print.V0      `toml:"cli-output"`
	// Tools          options_tools.Options `toml:"tools"`

	Query string `toml:"query,omitempty"`
}

func (blob V0) GetWorkspaceConfig() Blob {
	return blob
}

func (blob V0) GetDefaults() config_mutable_blobs.Defaults {
	return blob.Defaults
}

func (blob V0) GetDefaultQueryGroup() string {
	return blob.Query
}

type blobV0Coder struct{}

func (blobV0Coder) DecodeFrom(
	subject TypeWithBlob,
	reader io.Reader,
) (n int64, err error) {
	blob := Blob(&V0{})
	subject.Object = &blob

	dec := toml.NewDecoder(reader)

	if err = dec.Decode(*subject.Object); err != nil {
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
	subject TypeWithBlob,
	writer io.Writer,
) (n int64, err error) {
	dec := toml.NewEncoder(writer)

	if err = dec.Encode(*subject.Object); err != nil {
		if err == io.EOF {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
