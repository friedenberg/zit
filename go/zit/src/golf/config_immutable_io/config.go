package config_immutable_io

import (
	"encoding/gob"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/delta/config_immutable"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
)

type ConfigLoaded struct {
	ids.Type
	ImmutableConfig          config_immutable.Config
	BlobStoreImmutableConfig env_dir.Config
}

func (c *ConfigLoaded) ReadFrom(r io.Reader) (n int64, err error) {
	switch c.Type.String() {
	case builtin_types.ImmutableConfigV1:
		c.ImmutableConfig = &config_immutable.TomlV1{}
		td := toml.NewDecoder(r)

		if err = td.Decode(c.ImmutableConfig); err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}

	case "":
		c.ImmutableConfig = &config_immutable.V0{}

		dec := gob.NewDecoder(r)

		if err = dec.Decode(c.ImmutableConfig); err != nil {
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

	c.BlobStoreImmutableConfig = env_dir.MakeConfigFromImmutableBlobConfig(
		c.ImmutableConfig.GetBlobStoreConfigImmutable(),
	)

	return
}

func (s *ConfigLoaded) WriteTo(w io.Writer) (n int64, err error) {
	switch s.Type.String() {
	case builtin_types.ImmutableConfigV1:
		te := toml.NewEncoder(w)

		if err = te.Encode(s.ImmutableConfig); err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
		}

	case "":
		dec := gob.NewEncoder(w)

		if err = dec.Encode(s.ImmutableConfig); err != nil {
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
