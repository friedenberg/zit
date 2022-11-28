package konfig

import (
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/sku"
	"github.com/friedenberg/zit/src/typ_toml"
)

var (
	typExpander kennung.Expander
)

func init() {
	typExpander = kennung.MakeExpanderRight(`-`)
}

type compiledTyp struct {
	Sku sku.Sku2[kennung.Typ, *kennung.Typ]
	Typ typ_toml.Typ
}

func makeCompiledTyp(n string) *compiledTyp {
	return &compiledTyp{
		Sku: sku.Sku2[kennung.Typ, *kennung.Typ]{
			Kennung: kennung.MustTyp(n),
		},
		Typ: typ_toml.Typ{
			ExecCommand:    &typ_toml.ScriptConfig{},
			Actions:        make(map[string]*typ_toml.Action),
			EtikettenRules: make(map[string]typ_toml.EtikettRule),
		},
	}
}

func (ct *compiledTyp) ApplyExpanded(c Compiled) {
	expandedActual := c.GetSortedTypenExpanded(ct.Sku.Kennung.String())

	for _, ex := range expandedActual {
		ct.Typ.Merge(&ex.Typ)
	}
}
