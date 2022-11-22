package konfig

import (
	"fmt"
	"sort"

	"github.com/friedenberg/zit/src/bravo/collections"
)

type Compiled struct {
	TypFileExtension    string
	ZettelFileExtension string
	DefaultTyp          *compiledTyp
	DefaultOrganizeExt  string
	EtikettenHidden     []string
	EtikettenToAddToNew []string
	ExtensionsToTypen   map[string]string
	Typen               collections.Set[*compiledTyp]
}

func MakeDefaultCompiled() Compiled {
	dt := &compiledTyp{
		Name:          collections.MakeStringValue("md"),
		InlineAkte:    true,
		FileExtension: "md",
	}

	return Compiled{
		TypFileExtension:    "typ",
		ZettelFileExtension: "zettel",
		DefaultTyp:          dt,
		DefaultOrganizeExt:  "md",
		ExtensionsToTypen:   make(map[string]string),
		Typen: collections.MakeSet[*compiledTyp](
			func(v *compiledTyp) string {
				if v == nil {
					return ""
				}

				return v.Name.String()
			},
			dt,
		),
	}
}

func makeCompiled(k tomlKonfig) (kc Compiled, err error) {
	kc = MakeDefaultCompiled()

	for tn, tv := range k.Tags {
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

	typen := kc.Typen.MutableCopy()

	for tn, tv := range k.Typen {
		if tv.FileExtension != "" {
			kc.ExtensionsToTypen[tv.FileExtension] = tn
		}

		ct := makeCompiledTyp(tn)
		ct.Apply(tv)
		typen.Add(ct)
	}

	kc.Typen = typen.Copy()

	typen.Each(
		func(ct *compiledTyp) (err error) {
			ct.ApplyExpanded(kc)
			return
		},
	)

	return
}

func (c Compiled) GetSortedTypenExpanded(v string) (expandedActual []*compiledTyp) {
	expandedMaybe := typExpander.Expand(v)
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
		return expandedActual[i].Name.Len() > expandedActual[j].Name.Len()
	})

	return
}

func (c Compiled) GetZettelFileExtension() string {
	return fmt.Sprintf(".%s", c.ZettelFileExtension)
}

func (k Compiled) GetTyp(n string) (ct *compiledTyp) {
	expandedActual := k.GetSortedTypenExpanded(n)

	if len(expandedActual) > 0 {
		ct = expandedActual[0]
	}

	return
}
