package organize_text

import (
	"flag"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/erworben_cli_print_options"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections_ptr"
	"github.com/friedenberg/zit/src/charlie/collections_value"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/hotel/sku"
	"github.com/friedenberg/zit/src/india/matcher"
	"github.com/friedenberg/zit/src/india/objekte_collections"
	"github.com/friedenberg/zit/src/india/sku_fmt"
	"github.com/friedenberg/zit/src/juliett/konfig"
)

type Flags struct {
	Options

	once           *sync.Once
	ExtraEtiketten collections_ptr.Flag[kennung.Etikett, *kennung.Etikett]
}

type Options struct {
	wasMade bool

	Konfig *konfig.Compiled

	commentMatchers   schnittstellen.SetLike[matcher.Matcher]
	rootEtiketten     kennung.EtikettSet
	Typ               kennung.Typ
	GroupingEtiketten kennung.Slice
	ExtraEtiketten    kennung.EtikettSet
	Transacted        schnittstellen.SetLike[*sku.Transacted]

	Expanders kennung.Abbr

	UsePrefixJoints            bool
	UseRightAlignedIndents     bool
	UseRefiner                 bool
	UseMetadateiHeader         bool
	IncludeEtikettenInBrackets bool

	PrintOptions       erworben_cli_print_options.PrintOptions
	organize           sku_fmt.Organize
	organizeNew        sku_fmt.OrganizeNew
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
			GroupingEtiketten: kennung.MakeSlice(),
			Transacted:        objekte_collections.MakeMutableSetMetadateiWithKennung(),
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
			GroupingEtiketten: kennung.MakeSlice(),
			Transacted:        objekte_collections.MakeMutableSetMetadateiWithKennung(),
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

	f.BoolVar(
		&o.IncludeEtikettenInBrackets,
		"include-etiketten",
		false,
		"include etiketten between brackets",
	)
}

func (o *Flags) GetOptions(
	printOptions erworben_cli_print_options.PrintOptions,
	q matcher.Query,
	organize *sku_fmt.Organize,
	organizeNew *sku_fmt.OrganizeNew,
) Options {
	o.once.Do(
		func() {
			o.Options.ExtraEtiketten = o.ExtraEtiketten.GetSetPtrLike()
		},
	)

	o.organize = *organize
	o.organizeNew = *organizeNew

	if q != nil {
		o.rootEtiketten = q.GetEtiketten()

		ks := collections_value.MakeMutableValueSet[matcher.Matcher](nil)

		if err := matcher.VisitAllMatchers(
			func(m matcher.Matcher) (err error) {
				switch m1 := m.(type) {
				case matcher.Negate:
					return ks.Add(m1)

				case *matcher.Negate:
					return ks.Add(m1)

				default:
					return
				}
			},
			// TODO-P1 modify sigil matcher to allow child traversal
			q,
		); err != nil {
			errors.PanicIfError(err)
		}

		o.commentMatchers = ks
	}

	o.PrintOptions = printOptions

	return o.Options
}

func (o Options) assignmentTreeConstructor() *AssignmentTreeConstructor {
	return &AssignmentTreeConstructor{
		Options: o,
	}
}

func (o Options) Factory() *Factory {
	return &Factory{
		Options: o,
	}
}

func (o Options) refiner() *Refiner {
	return &Refiner{
		Enabled:         o.UseRefiner,
		UsePrefixJoints: o.UsePrefixJoints,
	}
}
