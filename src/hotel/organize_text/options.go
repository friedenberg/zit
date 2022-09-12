package organize_text

import (
	"flag"

	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/hinweis"
	"github.com/friedenberg/zit/src/golf/zettel_transacted"
)

type Options struct {
	hinweis.Abbr
	RootEtiketten          etikett.Set
	GroupingEtiketten      etikett.Slice
	ExtraEtiketten         etikett.Set
	Transacted             zettel_transacted.Set
	UsePrefixJoints        bool
	UseRightAlignedIndents bool
	UseAbbr                bool
	UseRefiner             bool
}

func MakeOptions() Options {
	return Options{
		GroupingEtiketten: etikett.NewSlice(),
		ExtraEtiketten:    etikett.MakeSet(),
	}
}

func (o *Options) AddToFlagSet(f *flag.FlagSet) {
	f.Var(&o.GroupingEtiketten, "group-by", "etikett prefixes to group zettels")
	f.Var(&o.ExtraEtiketten, "extras", "etiketten to always add to the organize text")
	f.BoolVar(&o.UsePrefixJoints, "prefix-joints", true, "split etiketten around hyphens")
	f.BoolVar(&o.UseRightAlignedIndents, "right-align", true, "right-align etiketten")
	f.BoolVar(&o.UseAbbr, "abbreviate-hinweisen", true, "abbreviate hinweisen")
	f.BoolVar(&o.UseRefiner, "refine", true, "refine the organize tree")
}

func (o Options) assignmentTreeConstructor() *AssignmentTreeConstructor {
	if !o.UseAbbr {
		o.Abbr = nil
	}

	return &AssignmentTreeConstructor{
		Options: o,
	}
}

func (o Options) refiner() *AssignmentTreeRefiner {
	return &AssignmentTreeRefiner{
		//TODO add to config
		Enabled:         o.UseRefiner,
		UsePrefixJoints: o.UsePrefixJoints,
	}
}
