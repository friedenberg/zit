package organize_text

import (
	"flag"

	"github.com/friedenberg/zit/src/charlie/gattung"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/juliett/konfig"
	zettel_pkg "github.com/friedenberg/zit/src/kilo/zettel"
)

type Options struct {
	Konfig konfig.Compiled
	Abbr   gattung.FuncAbbrIdMitKorper

	RootEtiketten     kennung.EtikettSet
	Typ               kennung.Typ
	GroupingEtiketten kennung.Slice
	ExtraEtiketten    kennung.EtikettSet
	Transacted        zettel_pkg.MutableSet

	UsePrefixJoints        bool
	UseRightAlignedIndents bool
	UseRefiner             bool
	UseMetadateiHeader     bool

	wasMade bool
}

func MakeOptions() Options {
	return Options{
		wasMade:           true,
		RootEtiketten:     kennung.MakeEtikettSet(),
		GroupingEtiketten: kennung.MakeSlice(),
		ExtraEtiketten:    kennung.MakeEtikettSet(),
		Transacted:        zettel_pkg.MakeMutableSetHinweis(0),
	}
}

func (o *Options) AddToFlagSet(f *flag.FlagSet) {
	f.Var(&o.GroupingEtiketten, "group-by", "etikett prefixes to group zettels")
	f.Var(&o.ExtraEtiketten, "extras", "etiketten to always add to the organize text")
	f.BoolVar(&o.UsePrefixJoints, "prefix-joints", true, "split etiketten around hyphens")
	f.BoolVar(&o.UseRightAlignedIndents, "right-align", true, "right-align etiketten")
	f.BoolVar(&o.UseRefiner, "refine", true, "refine the organize tree")
	f.BoolVar(&o.UseMetadateiHeader, "metadatei-header", true, "metadatei header")
}

func (o Options) assignmentTreeConstructor() *AssignmentTreeConstructor {
	if !o.Konfig.PrintAbbreviatedHinweisen {
		o.Abbr = nil
	}

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
