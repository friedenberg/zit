package options_print

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
)

type Abbreviations struct {
	ZettelIds bool `toml:"hinweisen"`
	Shas      bool `toml:"shas"`
}

type Box struct {
	PrintIncludeDescription bool `toml:"print-include-description"`
	PrintTime               bool `toml:"print-time"`
	PrintTai                bool `toml:"-"`
	PrintTagsAlways         bool `toml:"print-etiketten-always"`
	PrintEmptyShas          bool `toml:"print-empty-shas"`
	PrintIncludeTypes       bool `toml:"print-include-typen"`
	DescriptionInBox        bool `toml:"-"`
	ExcludeFields           bool `toml:"-"`
	PrintState              bool `toml:"-"`
}

type General struct {
	Abbreviations Abbreviations `toml:"abbreviations"`
	Box
	PrintMatchedDormant bool `toml:"print-matched-archiviert"`
	PrintShas           bool `toml:"print-shas"`
	PrintFlush          bool `toml:"print-flush"`
	PrintUnchanged      bool `toml:"print-unchanged"`
	PrintColors         bool `toml:"print-colors"`
	PrintInventoryLists bool `toml:"print-bestandsaufnahme"`
	ZittishNewlines     bool `toml:"-"`
}

func (a *General) Merge(b General, mask General) {
	if mask.Abbreviations.ZettelIds {
		a.Abbreviations.ZettelIds = b.Abbreviations.ZettelIds
	}

	if mask.Abbreviations.Shas {
		a.Abbreviations.Shas = b.Abbreviations.Shas
	}

	if mask.PrintIncludeTypes {
		a.PrintIncludeTypes = b.PrintIncludeTypes
	}

	if mask.PrintIncludeDescription {
		a.PrintIncludeDescription = b.PrintIncludeDescription
	}

	if mask.PrintTime {
		a.PrintTime = b.PrintTime
	}

	if mask.PrintTagsAlways {
		a.PrintTagsAlways = b.PrintTagsAlways
	}

	if mask.PrintEmptyShas {
		a.PrintEmptyShas = b.PrintEmptyShas
	}

	if mask.PrintMatchedDormant {
		a.PrintMatchedDormant = b.PrintMatchedDormant
	}

	if mask.PrintShas {
		a.PrintShas = b.PrintShas
	}

	if mask.PrintFlush {
		a.PrintFlush = b.PrintFlush
	}

	if mask.PrintUnchanged {
		a.PrintUnchanged = b.PrintUnchanged
	}

	if mask.PrintColors {
		a.PrintColors = b.PrintColors
	}

	if mask.PrintInventoryLists {
		a.PrintInventoryLists = b.PrintInventoryLists
	}

	if mask.DescriptionInBox {
		a.DescriptionInBox = b.DescriptionInBox
	}

	a.ZittishNewlines = b.ZittishNewlines
}

func Default() General {
	return General{
		Abbreviations: Abbreviations{
			ZettelIds: true,
			Shas:      true,
		},
		Box: Box{
			PrintIncludeTypes:       true,
			PrintIncludeDescription: true,
			PrintTime:               true,
			PrintTagsAlways:         true,
			PrintEmptyShas:          false,
		},
		PrintMatchedDormant: false,
		PrintShas:           true,
		PrintFlush:          true,
		PrintUnchanged:      true,
		PrintColors:         true,
		PrintInventoryLists: true,
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

func (c *General) AddToFlags(f *flag.FlagSet, m *General) {
	boolVarWithMask(
		f,
		"print-typen",
		&c.PrintIncludeTypes,
		&m.PrintIncludeTypes,
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
		"abbreviate-zettel-ids",
		&c.Abbreviations.ZettelIds,
		&m.Abbreviations.ZettelIds,
		"",
	)

	boolVarWithMask(
		f,
		"print-description",
		&c.PrintIncludeDescription,
		&m.PrintIncludeDescription,
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
		&c.PrintTagsAlways,
		&m.PrintTagsAlways,
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
		&c.PrintMatchedDormant,
		&m.PrintMatchedDormant,
		"",
	)

	boolVarWithMask(
		f,
		"print-shas",
		&c.PrintShas,
		&m.PrintShas,
		"",
	)

	boolVarWithMask(
		f,
		"print-flush",
		&c.PrintFlush,
		&m.PrintFlush,
		"",
	)

	boolVarWithMask(
		f,
		"print-unchanged",
		&c.PrintUnchanged,
		&m.PrintUnchanged,
		"",
	)

	boolVarWithMask(
		f,
		"print-colors",
		&c.PrintColors,
		&m.PrintColors,
		"",
	)

	boolVarWithMask(
		f,
		"print-bestandsaufnahme",
		&c.PrintInventoryLists,
		&m.PrintInventoryLists,
		"",
	)

	boolVarWithMask(
		f,
		"boxed-description",
		&c.DescriptionInBox,
		&m.DescriptionInBox,
		"",
	)
}

func (c General) WithPrintShas(v bool) General {
	c.PrintShas = v
	return c
}

func (c General) WithDescriptionInBox(v bool) General {
	c.DescriptionInBox = v
	return c
}

func (c General) WithPrintTai(v bool) General {
	c.PrintTai = v
	return c
}

func (c General) WithExcludeFields(v bool) General {
	c.ExcludeFields = v
	return c
}

func (c General) WithPrintTime(v bool) General {
	c.PrintTime = v
	return c
}

func (c General) WithPrintState(v bool) General {
	c.PrintState = v
	return c
}
