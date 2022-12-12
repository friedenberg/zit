package konfig_compiled

import (
	"encoding/gob"
	"fmt"
	"os"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/delta/collections"
	"github.com/friedenberg/zit/src/echo/standort"
	"github.com/friedenberg/zit/src/etikett"
	"github.com/friedenberg/zit/src/foxtrot/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/india/typ"
)

var (
	typExpander kennung.Expander
)

func init() {
	typExpander = kennung.MakeExpanderRight(`-`)
}

type Compiled struct {
	cli
	compiled
}

// TODO-P4 rename this
func (a *Compiled) ResetWithInner(b compiled) {
	a.compiled = b
}

type cli = konfig.Cli

type compiled struct {
	hasChanges bool

	Sku sku.Transacted[kennung.Konfig, *kennung.Konfig]

	konfig.Toml

	//Etiketten
	EtikettenHidden     []string
	EtikettenToAddToNew []string
	Etiketten           etikettSet

	//Typen
	ExtensionsToTypen map[string]string
	DefaultTyp        typ.Transacted
	Typen             typSet
}

func Make(
	s standort.Standort,
	kcli konfig.Cli,
) (c *Compiled, err error) {
	c = &Compiled{
		cli: kcli,
		compiled: compiled{
			ExtensionsToTypen: make(map[string]string),
		},
	}

	var f *os.File

	if f, err = files.Open(s.FileKonfigCompiled()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	dec := gob.NewDecoder(f)

	if err = dec.Decode(&c.compiled); err != nil {
		if errors.IsEOF(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	return
}

func (kc Compiled) Cli() konfig.Cli {
	return kc.cli
}

func (kc *compiled) recompile() (err error) {
	kc.hasChanges = true

	if err = kc.Etiketten.Each(
		func(ct *etikett.Transacted) (err error) {
			tn := ct.Sku.Kennung.String()
			tv := ct.Objekte.Akte

			switch {
			case tv.Hide:
				kc.EtikettenHidden = append(kc.EtikettenHidden, tn)

			case tv.AddToNewZettels:
				kc.EtikettenToAddToNew = append(kc.EtikettenToAddToNew, tn)
			}

			kc.applyExpandedEtikett(ct)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sort.Slice(kc.EtikettenHidden, func(i, j int) bool {
		return kc.EtikettenHidden[i] < kc.EtikettenHidden[j]
	})

	sort.Slice(kc.EtikettenToAddToNew, func(i, j int) bool {
		return kc.EtikettenToAddToNew[i] < kc.EtikettenToAddToNew[j]
	})

	if err = kc.Typen.Each(
		func(ct *typ.Transacted) (err error) {
			fe := ct.Objekte.Akte.FileExtension

			if fe != "" {
				kc.ExtensionsToTypen[fe] = ct.Sku.Kennung.String()
			}

			kc.applyExpandedTyp(ct)

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (kc *compiled) Flush(s standort.Standort) (err error) {
	if !kc.hasChanges {
		return
	}

	if err = kc.recompile(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var f *os.File

	if f, err = files.OpenExclusiveWriteOnlyTruncate(
		s.FileKonfigCompiled(),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	dec := gob.NewEncoder(f)

	if err = dec.Encode(kc); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c compiled) GetZettelFileExtension() string {
	return fmt.Sprintf(".%s", c.FileExtensions.Zettel)
}

//TODO-P3 merge all the below
func (c compiled) GetSortedTypenExpanded(v string) (expandedActual []*typ.Transacted) {
	expandedMaybe := collections.MakeMutableValueSet[collections.StringValue, *collections.StringValue]()
	typExpander.Expand(expandedMaybe, v)
	expandedActual = make([]*typ.Transacted, 0)

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

func (c compiled) GetSortedEtikettenExpanded(
	v string,
) (expandedActual []*etikett.Transacted) {
	expandedMaybe := collections.MakeMutableValueSet[collections.StringValue, *collections.StringValue]()
	typExpander.Expand(expandedMaybe, v)
	expandedActual = make([]*etikett.Transacted, 0)

	expandedMaybe.Each(
		func(v collections.StringValue) (err error) {
			ct, ok := c.Etiketten.Get(v.String())

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

func (kc compiled) GetTyp(k kennung.Typ) (ct *typ.Transacted) {
	expandedActual := kc.GetSortedTypenExpanded(k.String())

	if len(expandedActual) > 0 {
		ct = expandedActual[0]
	}

	return
}

func (kc compiled) GetEtikett(k kennung.Etikett) (ct *etikett.Transacted) {
	expandedActual := kc.GetSortedEtikettenExpanded(k.String())

	if len(expandedActual) > 0 {
		ct = expandedActual[0]
	}

	return
}

func (k *compiled) SetTransacted(
	kt *konfig.Transacted,
) {
	k.hasChanges = true
	k.Sku = kt.Sku
	k.Toml = kt.Objekte.Akte
}

func (k *compiled) AddTyp(
	ct *typ.Transacted,
) {
  if ct.Objekte.Akte.Actions == nil {
    panic(errors.Errorf("actions were nil"))
  }

  if ct.Objekte.Akte.EtikettenRules == nil {
    panic(errors.Errorf("etiketten rules were nil"))
  }

	k.hasChanges = true
	m := k.Typen.Elements()
	m = append(m, ct)
	k.Typen = makeCompiledTypSetFromSlice(m)

	return
}

func (k *compiled) AddEtikett(
	ct *etikett.Transacted,
) {
	k.hasChanges = true
	m := k.Etiketten.Elements()
	m = append(m, ct)
	k.Etiketten = makeCompiledEtikettSetFromSlice(m)

	return
}

func (c *compiled) applyExpandedTyp(ct *typ.Transacted) {
	expandedActual := c.GetSortedTypenExpanded(ct.Sku.Kennung.String())

	for _, ex := range expandedActual {
		ct.Objekte.Akte.Merge(&ex.Objekte.Akte)
	}
}

func (c *compiled) applyExpandedEtikett(ct *etikett.Transacted) {
	expandedActual := c.GetSortedEtikettenExpanded(ct.Sku.Kennung.String())

	for _, ex := range expandedActual {
		ct.Objekte.Akte.Merge(&ex.Objekte.Akte)
	}
}
