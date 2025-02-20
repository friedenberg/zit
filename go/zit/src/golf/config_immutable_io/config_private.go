package config_immutable_io

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/config_immutable"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type ConfigLoadedPrivate struct {
	ids.Type
	ImmutableConfig          config_immutable.ConfigPrivate
	BlobStoreImmutableConfig env_dir.Config // TODO extricate from env_dir
}

func (c *ConfigLoadedPrivate) GetType() ids.Type {
	return c.Type
}

func (c *ConfigLoadedPrivate) SetType(t ids.Type) {
	c.Type = t
}
