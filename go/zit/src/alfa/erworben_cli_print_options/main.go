package erworben_cli_print_options

import (
	"flag"

	"github.com/friedenberg/zit/src/bravo/values"
)

type Abbreviations struct {
	Hinweisen bool `toml:"hinweisen"`
	Shas      bool `toml:"shas"`
}

type PrintOptions struct {
	Abbreviations             Abbreviations `toml:"abbreviations"`
	PrintIncludeTypen         bool          `toml:"print-include-typen"`
	PrintIncludeBezeichnungen bool          `toml:"print-include-bezeichnungen"`
	PrintTime                 bool          `toml:"print-time"`
	PrintEtikettenAlways      bool          `toml:"print-etiketten-always"`
	PrintEmptyShas            bool          `toml:"print-empty-shas"`
	PrintMatchedArchiviert    bool          `toml:"print-matched-archiviert"`
	PrintShas                 bool          `toml:"print-shas"`
	ZittishNewlines           bool          `toml:"-"`
}

func (a *PrintOptions) Merge(b PrintOptions, mask PrintOptions) {
	if mask.Abbreviations.Hinweisen {
		a.Abbreviations.Hinweisen = b.Abbreviations.Hinweisen
	}

	if mask.Abbreviations.Shas {
		a.Abbreviations.Shas = b.Abbreviations.Shas
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

	if mask.PrintMatchedArchiviert {
		a.PrintMatchedArchiviert = b.PrintMatchedArchiviert
	}

	if mask.PrintShas {
		a.PrintShas = b.PrintShas
	}

	a.ZittishNewlines = b.ZittishNewlines
}

func Default() PrintOptions {
	return PrintOptions{
		Abbreviations: Abbreviations{
			Hinweisen: true,
			Shas:      true,
		},
		PrintIncludeTypen:         true,
		PrintIncludeBezeichnungen: true,
		PrintTime:                 true,
		PrintEtikettenAlways:      true,
		PrintEmptyShas:            false,
		PrintMatchedArchiviert:    false,
		PrintShas:                 true,
	}
}

func boolVarWithMask(
	f *flag.FlagSet,
	name string,
	v *bool,
	m *bool,
	desc string,
) {
	f.Func(name,
		desc,
		func(value string) (err error) {
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
		&c.Abbreviations.Shas,
		&m.Abbreviations.Shas,
		"",
	)

	boolVarWithMask(
		f,
		"abbreviate-hinweisen",
		&c.Abbreviations.Hinweisen,
		&m.Abbreviations.Hinweisen,
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

	boolVarWithMask(
		f,
		"print-matched-archiviert",
		&c.PrintMatchedArchiviert,
		&m.PrintMatchedArchiviert,
		"",
	)

	boolVarWithMask(
		f,
		"print-shas",
		&c.PrintShas,
		&m.PrintShas,
		"",
	)
}
