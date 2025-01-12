package repo_layout

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/triple_hyphen_io"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
)

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
		Metadata: metadata{Config: &s.Config},
		Blob:     &s.Config,
	}

	if _, err = thr.ReadFrom(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	s.Config.Blob = dir_layout.MakeConfigFromImmutableBlobConfig(
		s.Config.Config.GetBlobStoreImmutableConfig(),
	)

	return
}

type Config struct {
	ids.Type
	Config immutable_config.Config
	Blob   dir_layout.Config
}

type metadata struct {
	*Config
}

func (m metadata) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	if n, err = format.ReadLines(
		r,
		ohio.MakeLineReaderRepeat(
			ohio.MakeLineReaderKeyValues(
				map[string]interfaces.FuncSetString{
					"!": m.Type.Set,
				},
			),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (m metadata) WriteTo(w io.Writer) (n int64, err error) {
	var n1 int
	n1, err = fmt.Fprintf(w, "! %s\n", m.Type.StringSansOp())
	n += int64(n1)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Config) ReadFrom(r io.Reader) (n int64, err error) {
	switch c.Type.String() {
	case builtin_types.ImmutableConfigV1:
		c.Config = &immutable_config.TomlV1{}
		td := toml.NewDecoder(r)

		if err = td.Decode(c.Config); err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}

	case "":
		c.Config = &immutable_config.V0{}

		dec := gob.NewDecoder(r)

		if err = dec.Decode(c.Config); err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}

	default:
		err = errors.Errorf("unsupported config type: %q", c.Type)
		return
	}

	c.Blob = dir_layout.MakeConfigFromImmutableBlobConfig(
		c.Config.GetBlobStoreImmutableConfig(),
	)

	return
}

func (s *Config) WriteTo(w io.Writer) (n int64, err error) {
	switch s.Type.String() {
	case builtin_types.ImmutableConfigV1:
		te := toml.NewEncoder(w)

		if err = te.Encode(s.Config); err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}

	case "":
		dec := gob.NewEncoder(w)

		if err = dec.Encode(s.Config); err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}

	default:
		err = errors.Errorf("unsupported config type: %q", s.Type)
		return
	}

	return
}
