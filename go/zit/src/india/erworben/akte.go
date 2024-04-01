package erworben

import (
	"reflect"

	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/erworben_tools"
	"code.linenisgreat.com/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/src/charlie/gattung"
	"code.linenisgreat.com/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/src/delta/file_extensions"
	"code.linenisgreat.com/zit/src/echo/kennung"
)

type Defaults struct {
	Typ       kennung.Typ       `toml:"typ"`
	Etiketten []kennung.Etikett `toml:"etiketten"`
}

type Akte struct {
	Defaults        Defaults                                `toml:"defaults"`
	HiddenEtiketten []kennung.Etikett                       `toml:"hidden-etiketten"`
	FileExtensions  file_extensions.FileExtensions          `toml:"file-extensions"`
	RemoteScripts   map[string]script_config.RemoteScript   `toml:"remote-scripts"`
	Actions         map[string]script_config.ScriptConfig   `toml:"actions,omitempty"`
	PrintOptions    erworben_cli_print_options.PrintOptions `toml:"cli-output"`
	Tools           erworben_tools.Tools                    `toml:"tools"`
	Filters         map[string]string                       `toml:"filters"`
}

func (Akte) GetGattung() schnittstellen.GattungLike {
	return gattung.Konfig
}

func (a Akte) GetFilters() map[string]string {
	return a.Filters
}

func (a Akte) Equals(b Akte) bool {
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

func (a *Akte) Reset() {
	a.FileExtensions.Reset()
	a.Defaults.Typ = kennung.Typ{}
	a.Defaults.Etiketten = make([]kennung.Etikett, 0)
	a.HiddenEtiketten = make([]kennung.Etikett, 0)
	a.RemoteScripts = make(map[string]script_config.RemoteScript)
	a.Actions = make(map[string]script_config.ScriptConfig)
	a.PrintOptions = erworben_cli_print_options.Default()
	a.Filters = make(map[string]string)
}

func (a *Akte) ResetWith(b *Akte) {
	a.FileExtensions.Reset()

	a.Defaults.Typ = b.Defaults.Typ

	a.Defaults.Etiketten = make([]kennung.Etikett, len(b.Defaults.Etiketten))
	copy(a.Defaults.Etiketten, b.Defaults.Etiketten)

	a.HiddenEtiketten = make([]kennung.Etikett, len(b.HiddenEtiketten))
	copy(a.HiddenEtiketten, b.HiddenEtiketten)

	a.RemoteScripts = b.RemoteScripts
	a.Actions = b.Actions
	a.PrintOptions = b.PrintOptions
	a.Filters = b.Filters
}
