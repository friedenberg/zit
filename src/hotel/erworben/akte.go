package erworben

import (
	"reflect"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/bravo/script_config"
	"github.com/friedenberg/zit/src/bravo/todo"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type Akte struct {
	DefaultTyp     kennung.Typ                           `toml:"default-typ"`
	FileExtensions FileExtensions                        `toml:"file-extensions"`
	RemoteScripts  map[string]script_config.RemoteScript `toml:"remote-scripts"`
	Recipients     []string                              `toml:"recipients"`
	Actions        map[string]script_config.ScriptConfig `toml:"actions,omitempty"`
}

func (_ Akte) GetGattung() schnittstellen.GattungLike {
	return gattung.Konfig
}

func (a Akte) Equals(b Akte) bool {
	todo.Change("don't use reflection for equality")

	if !a.DefaultTyp.Equals(b.DefaultTyp) {
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

	return true
}

func (a *Akte) Reset() {
	a.FileExtensions.Reset()
	a.DefaultTyp = kennung.Typ{}
	a.RemoteScripts = make(map[string]script_config.RemoteScript)
	// TODO-P4 should reuse
	a.Recipients = make([]string, 0)
	a.Actions = make(map[string]script_config.ScriptConfig)
}

func (a *Akte) ResetWith(b Akte) {
	a.FileExtensions.Reset()
	// TODO-P4 should copy
	a.DefaultTyp = b.DefaultTyp
	a.RemoteScripts = b.RemoteScripts
	a.Recipients = b.Recipients
	a.Actions = b.Actions
}
