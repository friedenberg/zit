package repo_layout

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/triple_hyphen_io"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
)

type config struct {
	tipe ids.Type
	immutable_config.Config
	storeVersion      immutable_config.StoreVersion
	compressionType   immutable_config.CompressionType
	lockInternalFiles bool
}

func (s *Layout) loadImmutableConfig() (err error) {
	var r io.Reader

	{
		var f *os.File

		if f, err = files.OpenExclusiveReadOnly(s.FileConfigPermanent()); err != nil {
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

	thr := triple_hyphen_io.Reader{
		Metadata: metadataReader{config: &s.config},
		Blob:     &s.config,
	}

	if _, err = thr.ReadFrom(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type metadataReader struct {
	*config
}

func (m metadataReader) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	if n, err = format.ReadLines(
		r,
		ohio.MakeLineReaderRepeat(
			ohio.MakeLineReaderKeyValues(
				map[string]interfaces.FuncSetString{
					"!": m.tipe.Set,
				},
			),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *config) ReadFrom(r io.Reader) (n int64, err error) {
	switch s.tipe.String() {
	case builtin_types.ImmutableConfigV1:
		s.Config = &immutable_config.TomlV1{}
		td := toml.NewDecoder(r)

		if err = td.Decode(s.Config); err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}

	case "":
		s.Config = &immutable_config.V0{}

		dec := gob.NewDecoder(r)

		if err = dec.Decode(s.Config); err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}

	default:
		err = errors.Errorf("unsupported config type: %q", s.tipe)
		return
	}

	s.storeVersion = immutable_config.MakeStoreVersion(s.GetStoreVersion())
	s.compressionType = s.GetCompressionType()
	s.lockInternalFiles = s.GetLockInternalFiles()

	return
}
