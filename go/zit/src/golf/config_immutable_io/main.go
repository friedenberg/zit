package config_immutable_io

import (
	"bytes"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/echo/triple_hyphen_io"
)

type Reader struct {
	*ConfigLoaded
}

func (s *Reader) ReadFromFile(
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

	if _, err = s.ReadFrom(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Reader) ReadFrom(r io.Reader) (n int64, err error) {
	thr := triple_hyphen_io.Reader{
		Metadata: metadata{ConfigLoaded: s.ConfigLoaded},
		Blob:     s.ConfigLoaded,
	}

	if n, err = thr.ReadFrom(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.ConfigLoaded.BlobStoreImmutableConfig = env_dir.MakeConfigFromImmutableBlobConfig(
		s.ConfigLoaded.ImmutableConfig.GetBlobStoreConfigImmutable(),
	)

	return
}
