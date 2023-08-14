package erworben_cli_print_options

import (
	"flag"

	"github.com/friedenberg/zit/src/bravo/values"
)

type PrintOptions struct {
	PrintAbbreviatedHinweisen bool `toml:"-"`
	PrintAbbreviatedKennungen bool `toml:"-"`
	PrintAbbreviatedShas      bool `toml:"-"`
	PrintIncludeTypen         bool `toml:"print-include-typen"`
	PrintIncludeBezeichnungen bool `toml:"print-include-bezeichnungen"`
	PrintTime                 bool `toml:"print-time"`
	PrintEtikettenAlways      bool `toml:"print-etiketten-always"`
	PrintEmptyShas            bool `toml:"print-empty-shas"`
}

func (a *PrintOptions) Merge(b PrintOptions, mask PrintOptions) {
	if mask.PrintAbbreviatedHinweisen {
		a.PrintAbbreviatedHinweisen = b.PrintAbbreviatedHinweisen
	}

	if mask.PrintAbbreviatedKennungen {
		a.PrintAbbreviatedKennungen = b.PrintAbbreviatedKennungen
	}

	if mask.PrintAbbreviatedShas {
		a.PrintAbbreviatedShas = b.PrintAbbreviatedShas
	}

	if mask.PrintIncludeTypen {
		a.PrintIncludeTypen = b.PrintIncludeTypen
	}

	if mask.PrintIncludeBezeichnungen {
		a.PrintIncludeBezeichnungen = b.PrintIncludeBezeichnungen
	}

	if mask.PrintTime {
		a.PrintTime = b.PrintTime
	}

	if mask.PrintEtikettenAlways {
		a.PrintEtikettenAlways = b.PrintEtikettenAlways
	}

	if mask.PrintEmptyShas {
		a.PrintEmptyShas = b.PrintEmptyShas
	}
}

func Default() PrintOptions {
	return PrintOptions{
		PrintAbbreviatedHinweisen: true,
		PrintAbbreviatedKennungen: true,
		PrintAbbreviatedShas:      true,
		PrintIncludeTypen:         true,
		PrintIncludeBezeichnungen: true,
		PrintTime:                 true,
		PrintEtikettenAlways:      true,
		PrintEmptyShas:            false,
	}
}

func boolVarWithMask(
	f *flag.FlagSet,
	name string,
	v *bool,
	m *bool,
	desc string,
) {
	f.Func(name, desc, func(value string) (err error) {
		var bv values.Bool

		*m = true

		if err = bv.Set(value); err != nil {
			return
		}

		*v = bv.Bool()

		return
	},
	)
}

func (c *PrintOptions) AddToFlags(f *flag.FlagSet, m *PrintOptions) {
	boolVarWithMask(
		f,
		"print-typen",
		&c.PrintIncludeTypen,
		&m.PrintIncludeTypen,
		"",
	)

	// TODO-P4 combine below three options
	boolVarWithMask(
		f,
		"abbreviate-shas",
		&c.PrintAbbreviatedShas,
		&m.PrintAbbreviatedShas,
		"",
	)

	boolVarWithMask(
		f,
		"abbreviate-hinweisen",
		&c.PrintAbbreviatedHinweisen,
		&m.PrintAbbreviatedHinweisen,
		"",
	)

	boolVarWithMask(
		f,
		"abbreviate-kennungen",
		&c.PrintAbbreviatedKennungen,
		&m.PrintAbbreviatedKennungen,
		"",
	)

	boolVarWithMask(
		f,
		"print-bezeichnungen",
		&c.PrintIncludeBezeichnungen,
		&m.PrintIncludeBezeichnungen,
		"",
	)

	boolVarWithMask(
		f,
		"print-time",
		&c.PrintTime,
		&m.PrintTime,
		"",
	)

	boolVarWithMask(
		f,
		"print-etiketten",
		&c.PrintEtikettenAlways,
		&m.PrintEtikettenAlways,
		"",
	)

	boolVarWithMask(
		f,
		"print-empty-shas",
		&c.PrintEmptyShas,
		&m.PrintEmptyShas,
		"",
	)
}
