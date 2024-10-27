package config

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/expansion"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/tag_paths"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

// TODO
func (k *Compiled) ApplyDormantAndRealizeTags(
	sk *sku.Transacted,
) (err error) {
	ui.Log().Print("applying konfig to:", sk)
	mp := &sk.Metadata

	mp.Cache.SetExpandedTags(ids.ExpandMany(
		mp.GetTags(),
		expansion.ExpanderRight,
	))

	g := genres.Must(sk.GetGenre())
	isTag := g == genres.Tag

	// if g.HasParents() {
	// 	k.SetHasChanges(fmt.Sprintf("adding etikett with parents: %s", sk))
	// }

	var tag ids.Tag

	// TODO better solution for "realizing" etiketten against Konfig.
	// Specifically, making this less fragile and dependent on remembering to do
	// ApplyToSku for each Sku. Maybe a factory?
	mp.Cache.TagPaths.Reset()
	mp.GetTags().Each(mp.Cache.TagPaths.AddTagOld)

	if isTag {
		ks := sk.ObjectId.String()

		if err = tag.Set(ks); err != nil {
			return
		}

		sk.Metadata.Cache.TagPaths.AddSelf(catgut.MakeFromString(ks))

		ids.ExpandOneInto(
			tag,
			ids.MakeTag,
			expansion.ExpanderRight,
			mp.Cache.GetExpandedTagsMutable(),
		)
	}

	if err = k.addSuperTags(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = k.addImplicitTags(sk); err != nil {
		err = errors.Wrap(err)
		return
	}

	sk.SetDormant(k.dormant.ContainsSku(sk))

	return
}

func (k *Compiled) addSuperTags(
	sk *sku.Transacted,
) (err error) {
	g := sk.GetGenre()

	var expanded []string
	var ks string

	switch g {
	case genres.Tag, genres.Type, genres.Repo:
		ks = sk.ObjectId.String()

		expansion.ExpanderRight.Expand(
			func(v string) (err error) {
				expanded = append(expanded, v)
				return
			},
			ks,
		)

	default:
		return
	}

	for _, ex := range expanded {
		if ex == ks || ex == "" {
			continue
		}

		var ek *sku.Transacted

		if ek, err = k.getTagOrRepoIdOrType(ex); err != nil {
			err = errors.Wrapf(err, "Expanded: %q", ex)
			return
		}

		if ek == nil {
			// this is ok because currently, konfig is applied twice. However, this
			// is fragile as the order in which this method is called is
			// non-deterministic and the `GetTag` call may request an Tag we
			// have not processed yet
			continue
		}

		if ek.Metadata.Cache.TagPaths.Paths.Len() <= 1 {
			ui.Log().Print(ks, ex, ek.Metadata.Cache.TagPaths)
			continue
		}

		prefix := catgut.MakeFromString(ex)

		a := &sk.Metadata.Cache.TagPaths
		b := &ek.Metadata.Cache.TagPaths

		ui.Log().Print("a", a)
		ui.Log().Print("b", b)

		ui.Log().Print("prefix", prefix)

		if err = a.AddSuperFrom(b, prefix); err != nil {
			err = errors.Wrap(err)
			return
		}

		ui.Log().Print("a after", a)
	}

	return
}

func (k *Compiled) addImplicitTags(
	sk *sku.Transacted,
) (err error) {
	mp := &sk.Metadata
	ie := ids.MakeTagMutableSet()

	addImpEts := func(e *ids.Tag) (err error) {
		p1 := tag_paths.MakePathWithType()
		p1.Type = tag_paths.TypeIndirect
		p1.Add(catgut.MakeFromString(e.String()))

		impl := k.getImplicitTags(e)

		if impl.Len() == 0 {
			sk.Metadata.Cache.TagPaths.AddPathWithType(p1)
			return
		}

		if err = impl.EachPtr(
			quiter.MakeChain(
				ie.AddPtr,
				func(e1 *ids.Tag) (err error) {
					p2 := p1.Clone()
					p2.Add(catgut.MakeFromString(e1.String()))
					sk.Metadata.Cache.TagPaths.AddPathWithType(p2)
					return
				},
			),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	mp.GetTags().EachPtr(addImpEts)

	typKonfig := k.getApproximatedType(mp.GetType()).ApproximatedOrActual()

	if typKonfig != nil {
		typKonfig.GetTags().EachPtr(ie.AddPtr)
		typKonfig.GetTags().EachPtr(addImpEts)
	}

	mp.Cache.SetImplicitTags(ie)

	return
}
