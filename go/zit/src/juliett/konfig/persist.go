package konfig

import (
	"encoding/gob"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/typ_akte"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (kc *Compiled) recompile(
	tagp interfaces.BlobGetterPutter[*typ_akte.V0],
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

	if err = kc.compiled.Etiketten.Each(
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

			if err = kc.ApplySchlummerndAndRealizeEtiketten(&ke.Transacted); err != nil {
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

func (kc *Compiled) recompileTypen(
	tagp interfaces.BlobGetterPutter[*typ_akte.V0],
) (err error) {
	inlineTypen := collections_value.MakeMutableValueSet[values.String](nil)

	defer func() {
		kc.InlineTypen = inlineTypen.CloneSetLike()
	}()

	if err = kc.Typen.Each(
		func(ct *sku.Transacted) (err error) {
			var ta *typ_akte.V0

			if ta, err = tagp.GetBlob(ct.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer tagp.PutBlob(ta)

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

			if err = kc.ApplySchlummerndAndRealizeEtiketten(ct); err != nil {
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

func (kc *Compiled) HasChanges() (ok bool) {
	kc.lock.Lock()
	defer kc.lock.Unlock()

	ok = len(kc.compiled.changes) > 0

	if ok {
		ui.Log().Print(kc.compiled.changes)
	}

	return
}

func (kc *Compiled) GetChanges() (out []string) {
	kc.lock.Lock()
	defer kc.lock.Unlock()

	out = make([]string, len(kc.changes))
	copy(out, kc.changes)

	return
}

func (kc *compiled) SetHasChanges(reason string) {
	kc.lock.Lock()
	defer kc.lock.Unlock()

	kc.setHasChanges(reason)
}

func (kc *compiled) setHasChanges(reason string) {
	ui.Log().FunctionName(1)
	kc.changes = append(kc.changes, reason)
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
	tagp interfaces.BlobGetterPutter[*typ_akte.V0],
	printerHeader interfaces.FuncIter[string],
) (err error) {
	if !kc.HasChanges() || kc.DryRun {
		return
	}

	wg := iter.MakeErrorWaitGroupParallel()
	wg.Do(func() (err error) {
		if err = kc.flushErworben(s, tagp, printerHeader); err != nil {
			if errors.IsNotExist(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}
		}

		return
	})

	if err = wg.GetError(); err != nil {
		err = errors.Wrap(err)
		return
	}

	kc.changes = kc.changes[:0]

	return
}

func (kc *Compiled) flushErworben(
	s standort.Standort,
	tagp interfaces.BlobGetterPutter[*typ_akte.V0],
	printerHeader interfaces.FuncIter[string],
) (err error) {
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
