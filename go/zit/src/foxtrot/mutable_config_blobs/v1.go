package mutable_config_blobs

import (
	"code.linenisgreat.com/zit/go/zit/src/bravo/options_tools"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type V1 struct {
	Defaults       DefaultsV1                     `toml:"defaults"`
	FileExtensions file_extensions.FileExtensions `toml:"file-extensions"`
	PrintOptions   options_print.General          `toml:"cli-output"`
	Tools          options_tools.Options          `toml:"tools"`
}

func (a *V1) Reset() {
	a.FileExtensions.Reset()
	a.Defaults.Type = ids.Type{}
	a.Defaults.Tags = make([]ids.Tag, 0)
	a.PrintOptions = options_print.Default()
}

func (a *V1) ResetWith(b *V1) {
	a.FileExtensions.Reset()

	a.Defaults.Type = b.Defaults.Type

	a.Defaults.Tags = make([]ids.Tag, len(b.Defaults.Tags))
	copy(a.Defaults.Tags, b.Defaults.Tags)

	a.PrintOptions = b.PrintOptions
}

func (a V1) GetFilters() map[string]string {
	return make(map[string]string)
}
