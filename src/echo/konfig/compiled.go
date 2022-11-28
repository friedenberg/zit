package konfig

import (
	"fmt"
	"sort"

	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/echo/sku"
	"github.com/friedenberg/zit/src/typ_toml"
)

type Compiled struct {
	ZettelFileExtension string
	DefaultOrganizeExt  string
	EtikettenHidden     []string
	EtikettenToAddToNew []string

	//Typen
	ExtensionsToTypen map[string]string
	TypFileExtension  string
	DefaultTyp        *compiledTyp
	Typen             compiledTypSet
}

func MakeDefaultCompiled() (c Compiled) {
	dt := &compiledTyp{
		Sku: sku.Sku2[kennung.Typ, *kennung.Typ]{
			Kennung: kennung.MustTyp("md"),
		},
		Typ: typ_toml.Typ{
			InlineAkte:    true,
			FileExtension: "md",
		},
	}

	typen := collections.MakeMutableSet[*compiledTyp](
		c.Typen.Key,
		dt,
	)

	c = Compiled{
		TypFileExtension:    "typ",
		ZettelFileExtension: "zettel",
		DefaultTyp:          dt,
		DefaultOrganizeExt:  "md",
		ExtensionsToTypen:   make(map[string]string),
		Typen:               makeCompiledTypSet(typen),
	}

	return
}

func makeCompiled(kt tomlKonfig) (kc Compiled, err error) {
	kc = MakeDefaultCompiled()

	for tn, tv := range kt.Tags {
		switch {
		case tv.Hide:
			kc.EtikettenHidden = append(kc.EtikettenHidden, tn)

		case tv.AddToNewZettels:
			kc.EtikettenToAddToNew = append(kc.EtikettenToAddToNew, tn)
		}
	}

	sort.Slice(kc.EtikettenHidden, func(i, j int) bool {
		return kc.EtikettenHidden[i] < kc.EtikettenHidden[j]
	})

	sort.Slice(kc.EtikettenToAddToNew, func(i, j int) bool {
		return kc.EtikettenToAddToNew[i] < kc.EtikettenToAddToNew[j]
	})

	typen := collections.MakeMutableSet[*compiledTyp](
		kc.Typen.Key,
		kc.DefaultTyp,
	)

	for tn, tv := range kt.Typen {
		if tv.FileExtension != "" {
			kc.ExtensionsToTypen[tv.FileExtension] = tn
		}

		ct := makeCompiledTyp(tn)
		ct.Typ.Apply(&tv)
		typen.Add(ct)
	}

	kc.Typen = makeCompiledTypSet(typen)

	typen.Each(
		func(ct *compiledTyp) (err error) {
			ct.ApplyExpanded(kc)
			return
		},
	)

	return
}

func (c Compiled) GetSortedTypenExpanded(v string) (expandedActual []*compiledTyp) {
	expandedMaybe := collections.MakeMutableValueSet[collections.StringValue, *collections.StringValue]()
	typExpander.Expand(expandedMaybe, v)
	expandedActual = make([]*compiledTyp, 0)

	expandedMaybe.Each(
		func(v collections.StringValue) (err error) {
			ct, ok := c.Typen.Get(v.String())

			if !ok {
				return
			}

			expandedActual = append(expandedActual, ct)

			return
		},
	)

	sort.Slice(expandedActual, func(i, j int) bool {
		return expandedActual[i].Sku.Kennung.Less(expandedActual[j].Sku.Kennung)
	})

	return
}

func (c Objekte) GetZettelFileExtension() string {
	return fmt.Sprintf(".%s", c.Akte.ZettelFileExtension)
}

func (kc Objekte) GetTyp(n string) (ct *compiledTyp) {
	expandedActual := kc.Akte.GetSortedTypenExpanded(n)

	if len(expandedActual) > 0 {
		ct = expandedActual[0]
	}

	return
}
