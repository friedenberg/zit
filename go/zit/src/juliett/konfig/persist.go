package konfig

import (
	"encoding/gob"
	"os"
	"sort"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/delta/typ_akte"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func (kc *Compiled) recompile(
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) (err error) {
	if err = kc.recompileEtiketten(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = kc.recompileTypen(tagp); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (kc *Compiled) recompileEtiketten() (err error) {
	kc.DefaultEtiketten = kennung.MakeEtikettSet(kc.Defaults.Etiketten...)

	kc.ImplicitEtiketten = make(implicitEtikettenMap)

	if err = kc.Etiketten.Each(
		func(ke *ketikett) (err error) {
			var e kennung.Etikett

			if err = e.Set(ke.Transacted.GetKennung().String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = kc.AccumulateImplicitEtiketten(e); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = kc.ApplyToSku(&ke.Transacted); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sort.Slice(kc.EtikettenHiddenStringsSlice, func(i, j int) bool {
		return kc.EtikettenHiddenStringsSlice[i] < kc.EtikettenHiddenStringsSlice[j]
	})

	kc.EtikettenHidden = kennung.MakeEtikettSet(
		kc.HiddenEtiketten...,
	)

	return
}

func (kc *Compiled) recompileTypen(
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
) (err error) {
	inlineTypen := collections_value.MakeMutableValueSet[values.String](nil)

	defer func() {
		kc.InlineTypen = inlineTypen.CloneSetLike()
	}()

	if err = kc.Typen.Each(
		func(ct *sku.Transacted) (err error) {
			var ta *typ_akte.V0

			if ta, err = tagp.GetAkte(ct.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer tagp.PutAkte(ta)

			fe := ta.FileExtension

			if fe == "" {
				fe = ct.GetKennung().String()
			}

			// TODO-P2 enforce uniqueness
			kc.ExtensionsToTypen[fe] = ct.GetKennung().String()
			kc.TypenToExtensions[ct.GetKennung().String()] = fe

			if ta.InlineAkte {
				inlineTypen.Add(values.MakeString(ct.Kennung.String()))
			}

			if err = kc.ApplyToSku(ct); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}
	return
}

func (kc *Compiled) HasChanges() bool {
	kc.lock.Lock()
	defer kc.lock.Unlock()

	return kc.hasChanges
}

func (kc *Compiled) SetHasChanges(v bool) {
	kc.lock.Lock()
	defer kc.lock.Unlock()

	kc.hasChanges = true
}

func (kc *Compiled) loadKonfigErworben(s standort.Standort) (err error) {
	var f *os.File

	p := s.FileKonfigCompiled()

	if kc.UseKonfigErworbenFile {
		p = s.FileKonfigErworben()
	}

	if f, err = files.Open(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	dec := gob.NewDecoder(f)

	if err = dec.Decode(&kc.compiled); err != nil {
		if errors.IsEOF(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	return
}

func (kc *Compiled) Flush(
	s standort.Standort,
	tagp schnittstellen.AkteGetterPutter[*typ_akte.V0],
  printerHeader schnittstellen.FuncIter[string],
) (err error) {
	if !kc.hasChanges || kc.DryRun {
		return
	}

	if err = printerHeader("recompiling konfig"); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = kc.recompile(tagp); err != nil {
		err = errors.Wrap(err)
		return
	}

	p := s.FileKonfigCompiled()

	if kc.UseKonfigErworbenFile {
		p = s.FileKonfigErworben()
	}

	var f *os.File

	if f, err = files.OpenExclusiveWriteOnlyTruncate(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	dec := gob.NewEncoder(f)

	if err = dec.Encode(&kc.compiled); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = printerHeader("recompiled konfig"); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
