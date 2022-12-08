package konfig

import (
	"fmt"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/delta/typ_toml"
	"github.com/friedenberg/zit/src/echo/sku"
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
		Sku: sku.Transacted[kennung.Typ, *kennung.Typ]{
			Kennung: kennung.MustTyp("md"),
		},
		Typ: typ_toml.Objekte{
			Akte: typ_toml.Typ{
				InlineAkte:    true,
				FileExtension: "md",
			},
		},
	}

	dt.generateSha()

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

func makeCompiled(
	kt objekteToml,
) (kc Compiled, sha sha.Sha, err error) {
	kc = MakeDefaultCompiled()

	for tn, tv := range kt.Konfig.Tags {
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

	kc.Typen = makeCompiledTypSet(typen)

	if sha, err = kc.recompile(kt.Sha); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Compiled) recompile(inSha sha.Sha) (outSha sha.Sha, err error) {
	shasTypen := sha.MakeMutableSet(inSha)

	if err = c.Typen.Each(
		func(ct *compiledTyp) (err error) {
			fe := ct.Typ.Akte.FileExtension

			if fe != "" {
				c.ExtensionsToTypen[fe] = ct.Sku.Kennung.String()
			}

			ct.ApplyExpanded(c)
			ct.generateSha()

			shasTypen.Add(ct.Sku.Sha)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	outSha = sha.ShaFromSet(shasTypen)

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
		return expandedActual[i].Sku.Kennung.Len() > expandedActual[j].Sku.Kennung.Len()
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
