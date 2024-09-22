package config

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
	"code.linenisgreat.com/zit/go/zit/src/delta/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
)

func (kc *Compiled) recompile(
	tagp interfaces.BlobGetterPutter[*type_blobs.V0],
) (err error) {
	if err = kc.recompileTags(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = kc.recompileTypes(tagp); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (kc *Compiled) recompileTags() (err error) {
	kc.DefaultTags = ids.MakeTagSet(kc.Defaults.Etiketten...)

	kc.ImplicitTags = make(implicitTagMap)

	if err = kc.compiled.Tags.Each(
		func(ke *tag) (err error) {
			var e ids.Tag

			if err = e.Set(ke.String()); err != nil {
				err = errors.Wrapf(err, "Sku: %s", sku_fmt.String(&ke.Transacted))
				return
			}

			if err = kc.AccumulateImplicitTags(e); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = kc.ApplyDormantAndRealizeTags(&ke.Transacted); err != nil {
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

func (kc *Compiled) recompileTypes(
	tagp interfaces.BlobGetterPutter[*type_blobs.V0],
) (err error) {
	inlineTypen := collections_value.MakeMutableValueSet[values.String](nil)

	defer func() {
		kc.InlineTypes = inlineTypen.CloneSetLike()
	}()

	if err = kc.Types.Each(
		func(ct *sku.Transacted) (err error) {
			var ta *type_blobs.V0

			if ta, err = tagp.GetBlob(ct.GetBlobSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer tagp.PutBlob(ta)

			fe := ta.FileExtension

			if fe == "" {
				fe = ct.GetObjectId().String()
			}

			// TODO-P2 enforce uniqueness
			kc.ExtensionsToTypes[fe] = ct.GetObjectId().String()
			kc.TypesToExtensions[ct.GetObjectId().String()] = fe

			if ta.InlineBlob {
				inlineTypen.Add(values.MakeString(ct.ObjectId.String()))
			}

			if err = kc.ApplyDormantAndRealizeTags(ct); err != nil {
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

func (kc *Compiled) loadMutableConfig(s fs_home.Home) (err error) {
	var f *os.File

	p := s.FileKonfigErworben()

	if f, err = files.Open(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

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
	s fs_home.Home,
	tagp interfaces.BlobGetterPutter[*type_blobs.V0],
	printerHeader interfaces.FuncIter[string],
) (err error) {
	if !kc.HasChanges() || kc.DryRun {
		return
	}

	wg := iter.MakeErrorWaitGroupParallel()
	wg.Do(func() (err error) {
		if err = kc.flushMutableConfig(s, tagp, printerHeader); err != nil {
			err = errors.Wrap(err)
			return
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

func (kc *Compiled) flushMutableConfig(
	s fs_home.Home,
	tagp interfaces.BlobGetterPutter[*type_blobs.V0],
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

	p := s.FileKonfigErworben()

	var f *os.File

	if f, err = files.OpenCreateWriteOnlyTruncate(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, f)

	enc := gob.NewEncoder(f)

	if err = enc.Encode(&kc.compiled); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = printerHeader("recompiled konfig"); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
