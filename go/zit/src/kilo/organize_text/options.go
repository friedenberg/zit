package organize_text

import (
	"flag"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_ptr"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/juliett/config"
)

type Flags struct {
	Options

	once      *sync.Once
	ExtraTags collections_ptr.Flag[ids.Tag, *ids.Tag]
}

type Options struct {
	wasMade bool

	Config *config.Compiled

	commentMatchers interfaces.SetLike[sku.Query]
	rootTags        ids.TagSet
	Type            ids.Type
	GroupingTags    ids.TagSlice
	ExtraTags       ids.TagSet
	Skus            interfaces.SetLike[sku.ExternalLike]

	ObjectFactory

	Abbr ids.Abbr

	UsePrefixJoints        bool
	UseRightAlignedIndents bool
	UseRefiner             bool
	UseMetadataHeader      bool

	PrintOptions       erworben_cli_print_options.PrintOptions
	stringFormatReader catgut.StringFormatReader[sku.ExternalLike]
	stringFormatWriter catgut.StringFormatWriter[sku.ExternalLike]
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
			Skus:         sku.MakeExternalLikeMutableSet(),
		},
	}
}

func MakeFlagsWithMetadata(m object_metadata.Metadata) Flags {
	return Flags{
		once: &sync.Once{},
		ExtraTags: collections_ptr.MakeFlagCommas[ids.Tag](
			collections_ptr.SetterPolicyAppend,
		),

		Options: Options{
			rootTags:     m.GetTags(),
			wasMade:      true,
			GroupingTags: ids.MakeTagSlice(),
			Skus:         sku.MakeExternalLikeMutableSet(),
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

	f.BoolVar(
		&o.UseRightAlignedIndents,
		"right-align",
		true,
		"right-align tags",
	)

	f.BoolVar(&o.UseRefiner, "refine", true, "refine the organize tree")

	f.BoolVar(
		&o.UseMetadataHeader,
		"metadatei-header",
		true,
		"metadatei header",
	)
}

func (o *Flags) GetOptions(
	printOptions erworben_cli_print_options.PrintOptions,
	q sku.QueryGroup,
	skuFmt sku_fmt.ExternalLike,
	abbr ids.Abbr,
	of ObjectFactory,
) Options {
	o.once.Do(
		func() {
			o.Options.ExtraTags = o.ExtraTags.GetSetPtrLike()
		},
	)

	o.stringFormatReader = skuFmt.ReaderExternalLike
	o.stringFormatWriter = skuFmt.WriterExternalLike

	if q == nil {
		o.rootTags = ids.MakeTagSet()
	} else {
		o.rootTags = q.GetTags()
	}

	of.SetDefaultsIfNecessary()

	o.ObjectFactory = of
	o.PrintOptions = printOptions
	o.Abbr = abbr

	return o.Options
}

func (o Options) Make() (ot *Text, err error) {
	c := &constructor{
		Text: Text{
			Options: o,
		},
	}

	if ot, err = c.Make(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (o Options) refiner() *Refiner {
	return &Refiner{
		Enabled:         o.UseRefiner,
		UsePrefixJoints: o.UsePrefixJoints,
	}
}
