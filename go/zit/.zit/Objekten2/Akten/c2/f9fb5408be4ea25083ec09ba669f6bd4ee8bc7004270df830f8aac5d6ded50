package config

import (
	"fmt"
	"sort"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/mutable_config_blobs"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func (c *compiled) GetZettelFileExtension() string {
	return fmt.Sprintf(".%s", c.GetFileExtensions().GetFileExtensionZettel())
}

func (kc *Compiled) GetImmutableConfig() interfaces.ImmutableConfig {
	return kc.immutable_config_private
}

func (kc *compiled) getType(k ids.IdLike) (ct *sku.Transacted) {
	if k.GetGenre() != genres.Type {
		return
	}

	if ct1, ok := kc.Types.Get(k.String()); ok {
		ct = ct1.CloneTransacted()
	}

	return
}

func (kc *compiled) getRepo(k ids.IdLike) (ct *sku.Transacted) {
	if k.GetGenre() != genres.Repo {
		return
	}

	if ct1, ok := kc.Repos.Get(k.String()); ok {
		ct = ct1.CloneTransacted()
	}

	return
}

// Returns the exactly matching Typ, or if it doesn't exist, returns the parent
// Typ or nil. (Parent Typ for `md-gdoc` would be `md`.)
func (kc *compiled) getApproximatedType(
	k ids.IdLike,
) (ct ApproximatedType) {
	if k.GetGenre() != genres.Type {
		return
	}

	expandedActual := kc.getSortedTypesExpanded(k.String())
	if len(expandedActual) > 0 {
		ct.HasValue = true
		ct.Type = expandedActual[0]

		if ids.Equals(ct.Type.GetObjectId(), k) {
			ct.IsActual = true
		}
	}

	return
}

func (kc *compiled) getTagOrRepoIdOrType(
	v string,
) (sk *sku.Transacted, err error) {
	var k ids.ObjectId

	if err = k.Set(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	switch k.GetGenre() {
	case genres.Tag:
		sk, _ = kc.getTag(&k)
	case genres.Repo:
		sk = kc.getRepo(&k)
	case genres.Type:
		sk = kc.getType(&k)

	default:
		err = genres.MakeErrUnsupportedGenre(&k)
		return
	}

	return
}

func (kc *compiled) getTag(
	k ids.IdLike,
) (ct *sku.Transacted, ok bool) {
	if k.GetGenre() != genres.Tag {
		return
	}

	expandedActual := kc.getSortedTagsExpanded(k.String())

	if len(expandedActual) > 0 {
		ct = expandedActual[0]
		ok = true
	}

	return
}

// TODO-P3 merge all the below
func (c *compiled) getSortedTypesExpanded(
	v string,
) (expandedActual []*sku.Transacted) {
	expandedMaybe := collections_value.MakeMutableValueSet[values.String](nil)

	sa := quiter.MakeFuncSetString(expandedMaybe)

	typeExpander.Expand(sa, v)
	expandedActual = make([]*sku.Transacted, 0)

	expandedMaybe.Each(
		func(v values.String) (err error) {
			c.lock.Lock()
			defer c.lock.Unlock()

			ct, ok := c.Types.Get(v.String())

			if !ok {
				return
			}

			expandedActual = append(expandedActual, ct)

			return
		},
	)

	sort.Slice(expandedActual, func(i, j int) bool {
		return len(
			expandedActual[i].GetObjectId().String(),
		) > len(
			expandedActual[j].GetObjectId().String(),
		)
	})

	return
}

func (c *compiled) getSortedTagsExpanded(
	v string,
) (expandedActual []*sku.Transacted) {
	c.lock.Lock()
	defer c.lock.Unlock()

	expandedMaybe := collections_value.MakeMutableValueSet[values.String](nil)
	sa := quiter.MakeFuncSetString(
		expandedMaybe,
	)
	typeExpander.Expand(sa, v)
	expandedActual = make([]*sku.Transacted, 0)

	expandedMaybe.Each(
		func(v values.String) (err error) {
			ct, ok := c.Tags.Get(v.String())

			if !ok {
				return
			}

			ct1 := sku.GetTransactedPool().Get()

			sku.Resetter.ResetWith(ct1, &ct.Transacted)

			expandedActual = append(expandedActual, ct1)

			return
		},
	)

	sort.Slice(expandedActual, func(i, j int) bool {
		return len(
			expandedActual[i].GetObjectId().String(),
		) > len(
			expandedActual[j].GetObjectId().String(),
		)
	})

	return
}

func (c *compiled) getImplicitTags(
	e *ids.Tag,
) ids.TagSet {
	s, ok := c.ImplicitTags[e.String()]

	if !ok || s == nil {
		return ids.MakeTagSet()
	}

	return s
}

func (kc *Compiled) Cli() mutable_config_blobs.Cli {
	return kc.cli
}
