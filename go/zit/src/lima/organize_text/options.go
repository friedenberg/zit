package organize_text

import (
	"flag"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_ptr"
	"code.linenisgreat.com/zit/go/zit/src/charlie/options_print"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/box_format"
)

type Flags struct {
	Options

	once      *sync.Once
	ExtraTags collections_ptr.Flag[ids.Tag, *ids.Tag]
}

type Options struct {
	wasMade bool

	Config interface {
		interfaces.MutableConfigDryRun
		interfaces.ConfigGetFilters
	}

	Metadata

	commentMatchers interfaces.SetLike[sku.Query]
	GroupingTags    ids.TagSlice
	ExtraTags       ids.TagSet
	Skus            sku.SkuTypeSet

	sku.ObjectFactory

	Abbr ids.Abbr

	UsePrefixJoints   bool
	UseRefiner        bool
	UseMetadataHeader bool

	PrintOptions options_print.V0
	fmtBox       *box_format.BoxCheckedOut
}

func MakeFlags() Flags {
	return Flags{
		once: &sync.Once{},
		ExtraTags: collections_ptr.MakeFlagCommas[ids.Tag](
			collections_ptr.SetterPolicyAppend,
		),

		Options: Options{
			wasMade:      true,
			GroupingTags: ids.MakeTagSlice(),
			Skus:         sku.MakeSkuTypeSetMutable(),
			Metadata:     NewMetadata(ids.RepoId{}),
		},
	}
}

func MakeFlagsWithMetadata(m Metadata) Flags {
	if m.TagSet == nil {
		m.TagSet = ids.MakeTagSet()
	}

	return Flags{
		once: &sync.Once{},
		ExtraTags: collections_ptr.MakeFlagCommas[ids.Tag](
			collections_ptr.SetterPolicyAppend,
		),

		Options: Options{
			Metadata:     m,
			wasMade:      true,
			GroupingTags: ids.MakeTagSlice(),
			Skus:         sku.MakeSkuTypeSetMutable(),
		},
	}
}

func (o *Flags) AddToFlagSet(f *flag.FlagSet) {
	f.Var(&o.GroupingTags, "group-by", "tag prefixes to group zettels")

	f.Var(
		o.ExtraTags,
		"extras",
		"tags to always add to the organize text",
	)

	f.BoolVar(
		&o.UsePrefixJoints,
		"prefix-joints",
		true,
		"split tags around hyphens",
	)

	f.BoolVar(&o.UseRefiner, "refine", true, "refine the organize tree")

	f.BoolVar(
		&o.UseMetadataHeader,
		"metadata-header",
		true,
		"metadata header",
	)
}

func (o *Flags) GetOptionsWithMetadata(
	printOptions options_print.V0,
	skuFmt *box_format.BoxCheckedOut,
	abbr ids.Abbr,
	of sku.ObjectFactory,
	m Metadata,
) Options {
	o.once.Do(
		func() {
			o.Options.ExtraTags = o.ExtraTags.GetSetPtrLike()
		},
	)

	o.fmtBox = skuFmt

	of.SetDefaultsIfNecessary()

	o.ObjectFactory = of
	o.PrintOptions = printOptions
	o.Abbr = abbr
	o.Metadata = m

	return o.Options
}

func (o *Flags) GetOptions(
	printOptions options_print.V0,
	q TagSetGetter,
	skuBoxFormat *box_format.BoxCheckedOut,
	abbr ids.Abbr, // TODO move Abbr as required arg
	of sku.ObjectFactory,
) Options {
	m := o.Metadata

	if q != nil {
		m.TagSet = q.GetTags()
	}

	if m.prototype == nil {
		panic("Metadata not initalized")
	}

	return o.GetOptionsWithMetadata(
		printOptions,
		skuBoxFormat,
		abbr,
		of,
		m,
	)
}

func (o Options) Make() (ot *Text, err error) {
	c := &constructor{
		Text: Text{
			Options: o,
		},
	}

	ot = &c.Text

	c.all = MakePrefixSet(0)
	c.Assignment = newAssignment(0)
	c.IsRoot = true

	if c.TagSet == nil {
		c.TagSet = ids.MakeTagSet()
	}

	if err = c.Options.Skus.Each(c.all.AddSku); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.preparePrefixSetsAndRootsAndExtras(); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.populate(); err != nil {
		err = errors.Wrap(err)
		return
	}

	c.Metadata.Type = c.Options.Type

	if err = ot.Refine(); err != nil {
		err = errors.Wrap(err)
		return
	}

	ot.SortChildren()

	return
}

func (o Options) refiner() *Refiner {
	return &Refiner{
		Enabled:         o.UseRefiner,
		UsePrefixJoints: o.UsePrefixJoints,
	}
}
