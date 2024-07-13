package user_ops

import (
	"io"
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/objekte_mode"
	"code.linenisgreat.com/zit/go/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/objekte_collections"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/kilo/zettel"
	"code.linenisgreat.com/zit/go/zit/src/november/umwelt"
)

type ZettelFromExternalAkte struct {
	*umwelt.Umwelt
	// TODO switch to using ObjekteOptions
	ProtoZettel zettel.ProtoZettel
	Filter      script_value.ScriptValue
	Delete      bool
	Dedupe      bool
}

func (c ZettelFromExternalAkte) Run(
	fdSet fd.Set,
) (results sku.TransactedMutableSet, err error) {
	if err = c.Lock(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, c.Unlock)

	toCreate := objekte_collections.MakeMutableSetUniqueAkte()
	toDelete := fd.MakeMutableSet()

	results = sku.MakeTransactedMutableSet()

	fds := make(map[sha.Bytes][]*fd.FD)

	for _, fd := range iter.SortedValues(fdSet) {
		if err = c.processOneFD(fd, fds); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	for _, fdsForSha := range fds {
		sort.Slice(fdsForSha, func(i, j int) bool {
			return fdsForSha[i].String() < fdsForSha[j].String()
		})

		var z *store_fs.External

		if z, err = c.zettelForAkte(fdsForSha); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = toCreate.Add(z); err != nil {
			err = errors.Wrap(err)
			return
		}

		if !c.Delete {
			return
		}

		for _, fd := range fdsForSha {
			if err = toDelete.Add(fd); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	// if c.Dedupe {
	// 	matcher := objekte_collections.MakeMutableMatchSet(toCreate)

	// 	if err = c.GetStore().Query(
	// 		qg,
	// 		iter.MakeChain(
	// 			matcher.Match,
	// 			func(sk *sku.Transacted) (err error) {
	// 				z := &sku.Transacted{}

	// 				if err = z.SetFromSkuLike(sk); err != nil {
	// 					err = errors.Wrap(err)
	// 					return
	// 				}

	// 				return results.Add(z)
	// 			},
	// 		),
	// 	); err != nil {
	// 		err = errors.Wrap(err)
	// 		return
	// 	}
	// }

	if err = results.Each(
		func(z *sku.Transacted) (err error) {
			if c.ProtoZettel.Apply(z, gattung.Zettel) {
				if err = c.GetStore().CreateOrUpdateFromTransacted(
					z,
					objekte_mode.ModeApplyProto,
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	sortedToCreated := iter.Elements[*store_fs.External](toCreate)

	sort.Slice(
		sortedToCreated,
		func(i, j int) bool {
			return sortedToCreated[i].GetAkteFD().String() < sortedToCreated[j].GetAkteFD().String()
		},
	)

	for _, z := range sortedToCreated {
		if z.Metadatei.IsEmpty() {
			return
		}

		if err = c.GetStore().CreateOrUpdateFromTransacted(
			z.GetSku(),
			objekte_mode.ModeApplyProto,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		// TODO switch to using ObjekteOptions
		if c.ProtoZettel.Apply(z, gattung.Zettel) {
			if err = c.GetStore().CreateOrUpdateFromTransacted(
				z.GetSku(),
				objekte_mode.ModeEmpty,
			); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		results.Add(z.GetSku())
	}

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	// TODO move to umwelt
	dp := c.Umwelt.PrinterFDDeleted()

	err = toDelete.Each(
		func(f *fd.FD) (err error) {
			if err = c.Umwelt.Standort().Delete(f.GetPath()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = dp(f); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *ZettelFromExternalAkte) processOneFD(
	f *fd.FD,
	fds map[sha.Bytes][]*fd.FD,
) (err error) {
	var r io.Reader

	if r, err = c.Filter.Run(f.GetPath()); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, &c.Filter)

	var akteWriter sha.WriteCloser

	if akteWriter, err = c.Standort().BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, akteWriter)

	if _, err = io.Copy(akteWriter, r); err != nil {
		err = errors.Wrap(err)
		return
	}

	f.SetShaLike(akteWriter.GetShaLike())

	key := sha.Make(f.GetShaLike()).GetBytes()
	existing := fds[key]
	existing = append(existing, f)
	fds[key] = existing

	return
}

func (c *ZettelFromExternalAkte) zettelForAkte(
	akteFDs []*fd.FD,
) (z *store_fs.External, err error) {
	akteFD := akteFDs[0]
	z = store_fs.GetExternalPool().Get()

	z.FDs.Akte.ResetWith(akteFD)

	if err = z.Transacted.Kennung.SetWithKennung(&kennung.Hinweis{}); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.ProtoZettel.ApplyWithAkteFD(z, akteFD); err != nil {
		err = errors.Wrap(err)
		return
	}

	z.SetAkteSha(akteFD.GetShaLike())

	return
}
