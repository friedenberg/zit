package organize_text

import (
	"flag"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections_ptr"
	"code.linenisgreat.com/zit/go/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/juliett/konfig"
)

type Flags struct {
	Options

	once           *sync.Once
	ExtraEtiketten collections_ptr.Flag[ids.Tag, *ids.Tag]
}

type Options struct {
	wasMade bool

	Konfig *konfig.Compiled

	commentMatchers   interfaces.SetLike[sku.Query]
	rootEtiketten     ids.TagSet
	Typ               ids.Type
	GroupingEtiketten ids.TagSlice
	ExtraEtiketten    ids.TagSet
	Transacted        interfaces.SetLike[*sku.Transacted]

	Abbr ids.Abbr

	UsePrefixJoints        bool
	UseRightAlignedIndents bool
	UseRefiner             bool
	UseMetadateiHeader     bool

	PrintOptions       erworben_cli_print_options.PrintOptions
	skuFmt             sku_fmt.Organize
	stringFormatWriter interfaces.StringFormatWriter[*sku.Transacted]
}

func MakeFlags() Flags {
	return Flags{
		once: &sync.Once{},
		ExtraEtiketten: collections_ptr.MakeFlagCommas[ids.Tag](
			collections_ptr.SetterPolicyAppend,
		),

		Options: Options{
			wasMade:           true,
			GroupingEtiketten: ids.MakeTagSlice(),
			Transacted:        sku.MakeTransactedMutableSet(),
		},
	}
}

func MakeFlagsWithMetadatei(m object_metadata.Metadata) Flags {
	ui.Debug().Print(m.GetTags())

	return Flags{
		once: &sync.Once{},
		ExtraEtiketten: collections_ptr.MakeFlagCommas[ids.Tag](
			collections_ptr.SetterPolicyAppend,
		),

		Options: Options{
			rootEtiketten:     m.GetTags(),
			wasMade:           true,
			GroupingEtiketten: ids.MakeTagSlice(),
			Transacted:        sku.MakeTransactedMutableSet(),
		},
	}
}

func (o *Flags) AddToFlagSet(f *flag.FlagSet) {
	f.Var(&o.GroupingEtiketten, "group-by", "etikett prefixes to group zettels")

	f.Var(
		o.ExtraEtiketten,
		"extras",
		"etiketten to always add to the organize text",
	)

	f.BoolVar(
		&o.UsePrefixJoints,
		"prefix-joints",
		true,
		"split etiketten around hyphens",
	)

	f.BoolVar(
		&o.UseRightAlignedIndents,
		"right-align",
		true,
		"right-align etiketten",
	)

	f.BoolVar(&o.UseRefiner, "refine", true, "refine the organize tree")

	f.BoolVar(
		&o.UseMetadateiHeader,
		"metadatei-header",
		true,
		"metadatei header",
	)
}

func (o *Flags) GetOptions(
	printOptions erworben_cli_print_options.PrintOptions,
	q sku.QueryGroup,
	skuFmt *sku_fmt.Organize,
	abbr ids.Abbr,
) Options {
	o.once.Do(
		func() {
			o.Options.ExtraEtiketten = o.ExtraEtiketten.GetSetPtrLike()
		},
	)

	o.skuFmt = *skuFmt

	if q == nil {
		o.rootEtiketten = ids.MakeTagSet()
	} else {
		o.rootEtiketten = q.GetTags()

		// TODO handle negated
		// ks := collections_value.MakeMutableValueSet[sku.Query](nil)

		// if err := query.VisitAllMatchers(
		// 	func(m sku.Query) (err error) {
		// 		if e, ok := m.(*query.Exp); ok && e.Negated {
		// 			return ks.Add(e)
		// 		}

		// 		return
		// 	},
		// 	// TODO-P1 modify sigil matcher to allow child traversal
		// 	q,
		// ); err != nil {
		// 	errors.PanicIfError(err)
		// }

		// o.commentMatchers = ks
	}

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
