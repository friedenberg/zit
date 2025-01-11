package repo_layout

import (
	"bytes"
	"encoding/gob"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/triple_hyphen_io"
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
		Metadata: &s.tipe,
		Blob:     &s.config,
	}

	if _, err = thr.ReadFrom(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *config) ReadFrom(r io.Reader) (n int64, err error) {
	switch s.tipe.String() {
	// case builtin_types.ImmutableConfigV1:
	// 	s.Config = &immutable_config.TodoP1{}

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
