package erworben_cli_print_options

import "flag"

type Options struct {
	PrintAbbreviatedHinweisen bool
	PrintAbbreviatedKennungen bool
	PrintAbbreviatedShas      bool
	PrintIncludeTypen         bool
	PrintIncludeBezeichnungen bool
	PrintTime                 bool
	PrintEtikettenAlways      bool
}

func (c *Options) AddToFlags(f *flag.FlagSet) {
	f.BoolVar(&c.PrintIncludeTypen, "print-typen", true, "")
	// TODO-P4 combine below three options
	f.BoolVar(&c.PrintAbbreviatedShas, "abbreviate-shas", true, "")
	f.BoolVar(&c.PrintAbbreviatedHinweisen, "abbreviate-hinweisen", true, "")
	f.BoolVar(&c.PrintAbbreviatedKennungen, "abbreviate-kennungen", true, "")
	f.BoolVar(&c.PrintIncludeBezeichnungen, "print-bezeichnungen", true, "")
	f.BoolVar(&c.PrintTime, "print-time", true, "")
	f.BoolVar(&c.PrintEtikettenAlways, "print-etiketten", false, "")
}
