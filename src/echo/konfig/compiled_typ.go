package konfig

import (
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/typ_toml"
)

var (
	typExpander kennung.Expander
)

func init() {
	typExpander = kennung.MakeExpanderRight(`-`)
}

type compiledTyp struct {
	Sha  sha.Sha
	Name collections.StringValue
	typ_toml.Typ
}

func makeCompiledTyp(n string) *compiledTyp {
	return &compiledTyp{
		Name: collections.MakeStringValue(n),
		Typ: typ_toml.Typ{
			ExecCommand:    &typ_toml.ScriptConfig{},
			Actions:        make(map[string]*typ_toml.Action),
			EtikettenRules: make(map[string]typ_toml.EtikettRule),
		},
	}
}

func (ct *compiledTyp) ApplyExpanded(c Compiled) {
	expandedActual := c.GetSortedTypenExpanded(ct.Name.String())

	for _, ex := range expandedActual {
		ct.Merge(&ex.Typ)
	}
}
