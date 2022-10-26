package konfig

import (
	"fmt"
	"sort"
)

type Compiled struct {
	ZettelFileExtension string
	DefaultTyp          string
	EtikettenHidden     []string
	EtikettenToAddToNew []string
	//TODO add typen extensions
}

func makeCompiled(k toml) (kc Compiled, err error) {
	kc.ZettelFileExtension = "md"
	kc.DefaultTyp = "md"

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

	return
}

func (c Compiled) GetZettelFileExtension() string {
	return fmt.Sprintf(".%s", c.ZettelFileExtension)
}
