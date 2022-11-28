package konfig

import (
	"crypto/sha256"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/toml"
	"github.com/friedenberg/zit/src/bravo/gattung"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/metadatei_io"
	"github.com/friedenberg/zit/src/echo/sku"
	"github.com/friedenberg/zit/src/objekte_format"
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
	Typ typ_toml.Objekte
}

func makeCompiledTyp(n string) *compiledTyp {
	return &compiledTyp{
		Sku: sku.Sku2[kennung.Typ, *kennung.Typ]{
			Kennung: kennung.MustTyp(n),
		},
		Typ: typ_toml.Objekte{
			Akte: typ_toml.Typ{
				ExecCommand:    &typ_toml.ScriptConfig{},
				Actions:        make(map[string]*typ_toml.Action),
				EtikettenRules: make(map[string]typ_toml.EtikettRule),
			},
		},
	}
}

func (ct compiledTyp) Gattung() gattung.Gattung {
	return gattung.Typ
}

func (ct compiledTyp) AkteSha() sha.Sha {
	return ct.Typ.Sha
}

func (ct *compiledTyp) SetAkteSha(v sha.Sha) {
	ct.Typ.Sha = v
}

func (ct compiledTyp) ObjekteSha() sha.Sha {
	return ct.Sku.Sha
}

func (ct *compiledTyp) SetObjekteSha(
	arf metadatei_io.AkteReaderFactory,
	v string,
) (err error) {
	if err = ct.Sku.Sha.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (ct *compiledTyp) generateSha() {
	hash := sha256.New()
	enc := toml.NewEncoder(hash)

	if err := enc.Encode(&ct.Typ); err != nil {
		panic(err)
	}

	ct.Typ.Sha = sha.FromHash(hash)

	f := objekte_format.MakeFormat(metadatei_io.NopAkteFactory())

	if _, err := f.WriteFormat(io.Discard, ct); err != nil {
		panic(err)
	}
}

func (ct *compiledTyp) ApplyExpanded(c Compiled) {
	expandedActual := c.GetSortedTypenExpanded(ct.Sku.Kennung.String())

	for _, ex := range expandedActual {
		ct.Typ.Akte.Merge(&ex.Typ.Akte)
	}
}
