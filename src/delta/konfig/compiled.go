package konfig

import (
	"fmt"
	"sort"

	"github.com/friedenberg/zit/src/bravo/collections"
)

type Compiled struct {
	ZettelFileExtension string
	DefaultTyp          *compiledTyp
	DefaultOrganizeExt  string
	EtikettenHidden     []string
	EtikettenToAddToNew []string
	ExtensionsToTypen   map[string]string
	TypenInline         collections.Set[string]
	Typen               collections.Set[*compiledTyp]
}

func MakeDefaultCompiled() Compiled {
	dt := &compiledTyp{
		Name:          "md",
		InlineAkte:    true,
		FileExtension: "md",
	}

	return Compiled{
		ZettelFileExtension: "md",
		DefaultTyp:          dt,
		DefaultOrganizeExt:  "md",
		ExtensionsToTypen:   make(map[string]string),
		TypenInline: collections.MakeSet[string](
			func(v string) string { return v },
			"md",
		),
		Typen: collections.MakeSet[*compiledTyp](
			func(v *compiledTyp) string {
				if v == nil {
					return ""
				}

				return v.Name
			},
			dt,
		),
	}
}

func makeCompiled(k toml) (kc Compiled, err error) {
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

	inlineTypen := kc.TypenInline.MutableCopy()
	typen := kc.Typen.MutableCopy()

	for tn, tv := range k.Typen {
		if tv.InlineAkte {
			inlineTypen.Add(tn)
		}

		if tv.FileExtension != "" {
			kc.ExtensionsToTypen[tv.FileExtension] = tn
		}

		ct := &compiledTyp{
			Name: tn,
		}

		ct.Apply(tv)
		typen.Add(ct)
	}

	typen.Each(
		func(ct *compiledTyp) (err error) {
			return
		},
	)

	kc.TypenInline = inlineTypen.Copy()
	kc.Typen = typen.Copy()

	return
}

func (c Compiled) GetZettelFileExtension() string {
	return fmt.Sprintf(".%s", c.ZettelFileExtension)
}
