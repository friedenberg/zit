package organize_text

import (
	"flag"
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/charlie/collections_ptr"
	"code.linenisgreat.com/zit/src/charlie/collections_value"
	"code.linenisgreat.com/zit/src/charlie/erworben_cli_print_options"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/src/juliett/konfig"
	"code.linenisgreat.com/zit/src/juliett/query"
)

type Flags struct {
	Options

	once           *sync.Once
	ExtraEtiketten collections_ptr.Flag[kennung.Etikett, *kennung.Etikett]
}

type Options struct {
	wasMade bool

	Konfig *konfig.Compiled

	commentMatchers   schnittstellen.SetLike[sku.Query]
	rootEtiketten     kennung.EtikettSet
	Typ               kennung.Typ
	GroupingEtiketten kennung.EtikettSlice
	ExtraEtiketten    kennung.EtikettSet
	Transacted        schnittstellen.SetLike[*sku.Transacted]

	Abbr kennung.Abbr

	UsePrefixJoints        bool
	UseRightAlignedIndents bool
	UseRefiner             bool
	UseMetadateiHeader     bool

	PrintOptions       erworben_cli_print_options.PrintOptions
	skuFmt             sku_fmt.Organize
	stringFormatWriter schnittstellen.StringFormatWriter[*sku.Transacted]
}

func MakeFlags() Flags {
	return Flags{
		once: &sync.Once{},
		ExtraEtiketten: collections_ptr.MakeFlagCommas[kennung.Etikett](
			collections_ptr.SetterPolicyAppend,
		),

		Options: Options{
			wasMade:           true,
			GroupingEtiketten: kennung.MakeEtikettSlice(),
			Transacted:        sku.MakeTransactedMutableSet(),
		},
	}
}

func MakeFlagsWithMetadatei(m metadatei.Metadatei) Flags {
	return Flags{
		once: &sync.Once{},
		ExtraEtiketten: collections_ptr.MakeFlagCommas[kennung.Etikett](
			collections_ptr.SetterPolicyAppend,
		),

		Options: Options{
			rootEtiketten:     m.GetEtiketten(),
			wasMade:           true,
			GroupingEtiketten: kennung.MakeEtikettSlice(),
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
	q *query.Group,
	skuFmt *sku_fmt.Organize,
	abbr kennung.Abbr,
) Options {
	o.once.Do(
		func() {
			o.Options.ExtraEtiketten = o.ExtraEtiketten.GetSetPtrLike()
		},
	)

	o.skuFmt = *skuFmt

	if q == nil {
		o.rootEtiketten = kennung.MakeEtikettSet()
	} else {
		o.rootEtiketten = q.GetEtiketten()

		ks := collections_value.MakeMutableValueSet[sku.Query](nil)

		if err := query.VisitAllMatchers(
			func(m sku.Query) (err error) {
				if e, ok := m.(*query.Exp); ok && e.Negated {
					return ks.Add(e)
				}

				return
			},
			// TODO-P1 modify sigil matcher to allow child traversal
			q,
		); err != nil {
			errors.PanicIfError(err)
		}

		o.commentMatchers = ks
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
