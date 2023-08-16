package erworben

import (
	"reflect"

	"github.com/friedenberg/zit/src/alfa/erworben_cli_print_options"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type Defaults struct {
	Typ       kennung.Typ       `toml:"typ"`
	Etiketten []kennung.Etikett `toml:"etiketten"`
}

type Akte struct {
	Defaults       Defaults                                `toml:"defaults"`
	FileExtensions FileExtensions                          `toml:"file-extensions"`
	RemoteScripts  map[string]script_config.RemoteScript   `toml:"remote-scripts"`
	Recipients     []string                                `toml:"recipients"`
	Actions        map[string]script_config.ScriptConfig   `toml:"actions,omitempty"`
	PrintOptions   erworben_cli_print_options.PrintOptions `toml:"cli-output"`
}

func (_ Akte) GetGattung() schnittstellen.GattungLike {
	return gattung.Konfig
}

func (a Akte) Equals(b Akte) bool {
	todo.Change("don't use reflection for equality")

	if !reflect.DeepEqual(a.Defaults.Etiketten, b.Defaults.Etiketten) {
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

	if !reflect.DeepEqual(a.Recipients, b.Recipients) {
		return false
	}

	if !reflect.DeepEqual(a.Actions, b.Actions) {
		return false
	}

	if !reflect.DeepEqual(a.PrintOptions, b.PrintOptions) {
		return false
	}

	return true
}

func (a *Akte) Reset() {
	a.FileExtensions.Reset()
	a.Defaults.Typ = kennung.Typ{}
	a.Defaults.Etiketten = make([]kennung.Etikett, 0)
	a.RemoteScripts = make(map[string]script_config.RemoteScript)
	// TODO-P4 should reuse
	a.Recipients = make([]string, 0)
	a.Actions = make(map[string]script_config.ScriptConfig)
	a.PrintOptions = erworben_cli_print_options.Default()
}

func (a *Akte) ResetWith(b Akte) {
	a.FileExtensions.Reset()
	// TODO-P4 should copy
	a.Defaults.Typ = b.Defaults.Typ
	a.Defaults.Etiketten = make([]kennung.Etikett, len(b.Defaults.Etiketten))
	copy(a.Defaults.Etiketten, b.Defaults.Etiketten)
	a.RemoteScripts = b.RemoteScripts
	a.Recipients = b.Recipients
	a.Actions = b.Actions
	a.PrintOptions = b.PrintOptions
}
