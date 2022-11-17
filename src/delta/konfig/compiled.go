package konfig

import (
	"fmt"
	"sort"

	"github.com/friedenberg/zit/src/bravo/collections"
)

type Compiled struct {
	ZettelFileExtension string
	DefaultTyp          string
	EtikettenHidden     []string
	EtikettenToAddToNew []string
	TypenExtensions     map[string]string
	TypenInline         collections.Set[string]
}

func MakeDefaultCompiled() Compiled {
	return Compiled{
		ZettelFileExtension: "md",
		DefaultTyp:          "md",
		TypenExtensions:     make(map[string]string),
		TypenInline: collections.MakeSet[string](
			func(v string) string { return v },
			"md",
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

	for tn, tv := range k.Typen {
		if tv.InlineAkte {
			inlineTypen.Add(tn)
		}

		if tv.FileExtension != "" {
			kc.TypenExtensions[tv.FileExtension] = tn
		}
	}

	kc.TypenInline = inlineTypen.Copy()

	return
}

func (c Compiled) GetZettelFileExtension() string {
	return fmt.Sprintf(".%s", c.ZettelFileExtension)
}
