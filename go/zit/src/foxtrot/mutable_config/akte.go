package mutable_config

import (
	"reflect"

	"code.linenisgreat.com/zit/go/zit/src/bravo/erworben_tools"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/charlie/print_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type Defaults struct {
	Typ       ids.Type  `toml:"typ"`
	Etiketten []ids.Tag `toml:"etiketten"`
}

type Blob struct {
	Defaults        Defaults                              `toml:"defaults"`
	HiddenEtiketten []ids.Tag                             `toml:"hidden-etiketten"`
	FileExtensions  file_extensions.FileExtensions        `toml:"file-extensions"`
	RemoteScripts   map[string]script_config.RemoteScript `toml:"remote-scripts"`
	Actions         map[string]script_config.ScriptConfig `toml:"actions,omitempty"`
	PrintOptions    print_options.General                 `toml:"cli-output"`
	Tools           erworben_tools.Tools                  `toml:"tools"`
	Filters         map[string]string                     `toml:"filters"`
}

func (a Blob) GetFilters() map[string]string {
	return a.Filters
}

func (a Blob) Equals(b Blob) bool {
	todo.Change("don't use reflection for equality")

	if !reflect.DeepEqual(a.Defaults.Etiketten, b.Defaults.Etiketten) {
		return false
	}

	if !reflect.DeepEqual(a.HiddenEtiketten, b.HiddenEtiketten) {
		return false
	}

	if !a.Defaults.Typ.Equals(b.Defaults.Typ) {
		return false
	}

	if !reflect.DeepEqual(a.FileExtensions, b.FileExtensions) {
		return false
	}

	if !reflect.DeepEqual(a.RemoteScripts, b.RemoteScripts) {
		return false
	}

	if !reflect.DeepEqual(a.Actions, b.Actions) {
		return false
	}

	if !reflect.DeepEqual(a.PrintOptions, b.PrintOptions) {
		return false
	}

	if !reflect.DeepEqual(a.Filters, b.Filters) {
		return false
	}

	return true
}

func (a *Blob) Reset() {
	a.FileExtensions.Reset()
	a.Defaults.Typ = ids.Type{}
	a.Defaults.Etiketten = make([]ids.Tag, 0)
	a.HiddenEtiketten = make([]ids.Tag, 0)
	a.RemoteScripts = make(map[string]script_config.RemoteScript)
	a.Actions = make(map[string]script_config.ScriptConfig)
	a.PrintOptions = print_options.Default()
	a.Filters = make(map[string]string)
}

func (a *Blob) ResetWith(b *Blob) {
	a.FileExtensions.Reset()

	a.Defaults.Typ = b.Defaults.Typ

	a.Defaults.Etiketten = make([]ids.Tag, len(b.Defaults.Etiketten))
	copy(a.Defaults.Etiketten, b.Defaults.Etiketten)

	a.HiddenEtiketten = make([]ids.Tag, len(b.HiddenEtiketten))
	copy(a.HiddenEtiketten, b.HiddenEtiketten)

	a.RemoteScripts = b.RemoteScripts
	a.Actions = b.Actions
	a.PrintOptions = b.PrintOptions
	a.Filters = b.Filters
}
