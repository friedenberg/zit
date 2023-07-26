package organize_text

import (
	"flag"
	"sync"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections2"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/golf/sku"
	"github.com/friedenberg/zit/src/india/konfig"
	"github.com/friedenberg/zit/src/india/objekte_collections"
)

type Flags struct {
	Options

	once           *sync.Once
	ExtraEtiketten collections2.Flag[kennung.Etikett, *kennung.Etikett]
}

type Options struct {
	wasMade bool

	Konfig konfig.Compiled

	RootEtiketten     kennung.EtikettSet
	Typ               kennung.Typ
	GroupingEtiketten kennung.Slice
	ExtraEtiketten    kennung.EtikettSet
	Transacted        schnittstellen.SetLike[sku.SkuLike]

	Expanders kennung.Abbr

	UsePrefixJoints        bool
	UseRightAlignedIndents bool
	UseRefiner             bool
	UseMetadateiHeader     bool
}

func MakeFlags() Flags {
	return Flags{
		once: &sync.Once{},
		ExtraEtiketten: collections2.MakeFlagCommas[kennung.Etikett](
			collections2.SetterPolicyAppend,
		),

		Options: Options{
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
}

func (o *Flags) GetOptions() Options {
	o.once.Do(
		func() {
			o.Options.ExtraEtiketten = o.ExtraEtiketten.GetSetPtrLike()
		},
	)

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
