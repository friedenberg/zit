package organize_text

import (
	"flag"

	"github.com/friedenberg/zit/src/charlie/kennung"
	"github.com/friedenberg/zit/src/delta/hinweis"
	"github.com/friedenberg/zit/src/delta/konfig"
	"github.com/friedenberg/zit/src/echo/typ"
	"github.com/friedenberg/zit/src/india/zettel_transacted"
)

type Options struct {
	konfig.Konfig
	hinweis.Abbr

	RootEtiketten     kennung.Set
	Typ               typ.Kennung
	GroupingEtiketten kennung.Slice
	ExtraEtiketten    kennung.Set
	Transacted        zettel_transacted.MutableSet

	UsePrefixJoints        bool
	UseRightAlignedIndents bool
	UseRefiner             bool
	UseMetadateiHeader     bool

	wasMade bool
}

func MakeOptions() Options {
	return Options{
		wasMade:           true,
		RootEtiketten:     kennung.MakeSet(),
		GroupingEtiketten: kennung.MakeSlice(),
		ExtraEtiketten:    kennung.MakeSet(),
		Transacted:        zettel_transacted.MakeMutableSetHinweis(0),
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
	if !o.PrintAbbreviatedHinweisen {
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
