package options_print

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
)

type Abbreviations struct {
	Hinweisen bool `toml:"hinweisen"`
	Shas      bool `toml:"shas"`
}

type Box struct {
	PrintIncludeDescription bool `toml:"print-include-description"`
	PrintTime               bool `toml:"print-time"`
	PrintTai                bool `toml:"-"`
	PrintTagsAlways         bool `toml:"print-etiketten-always"`
	PrintEmptyShas          bool `toml:"print-empty-shas"`
	PrintIncludeTypen       bool `toml:"print-include-typen"`
	DescriptionInBox        bool `toml:"-"`
	ExcludeFields           bool `toml:"-"`
	PrintState              bool `toml:"-"`
}

type General struct {
	Abbreviations Abbreviations `toml:"abbreviations"`
	Box
	PrintMatchedArchiviert bool `toml:"print-matched-archiviert"`
	PrintShas              bool `toml:"print-shas"`
	PrintFlush             bool `toml:"print-flush"`
	PrintUnchanged         bool `toml:"print-unchanged"`
	PrintColors            bool `toml:"print-colors"`
	PrintBestandsaufnahme  bool `toml:"print-bestandsaufnahme"`
	ZittishNewlines        bool `toml:"-"`
}

func (a *General) Merge(b General, mask General) {
	if mask.Abbreviations.Hinweisen {
		a.Abbreviations.Hinweisen = b.Abbreviations.Hinweisen
	}

	if mask.Abbreviations.Shas {
		a.Abbreviations.Shas = b.Abbreviations.Shas
	}

	if mask.PrintIncludeTypen {
		a.PrintIncludeTypen = b.PrintIncludeTypen
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

	if mask.PrintMatchedArchiviert {
		a.PrintMatchedArchiviert = b.PrintMatchedArchiviert
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

	if mask.PrintBestandsaufnahme {
		a.PrintBestandsaufnahme = b.PrintBestandsaufnahme
	}

	if mask.DescriptionInBox {
		a.DescriptionInBox = b.DescriptionInBox
	}

	a.ZittishNewlines = b.ZittishNewlines
}

func Default() General {
	return General{
		Abbreviations: Abbreviations{
			Hinweisen: true,
			Shas:      true,
		},
		Box: Box{
			PrintIncludeTypen:       true,
			PrintIncludeDescription: true,
			PrintTime:               true,
			PrintTagsAlways:         true,
			PrintEmptyShas:          false,
		},
		PrintMatchedArchiviert: false,
		PrintShas:              true,
		PrintFlush:             true,
		PrintUnchanged:         true,
		PrintColors:            true,
		PrintBestandsaufnahme:  true,
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
		"abbreviate-zettel-ids",
		&c.Abbreviations.Hinweisen,
		&m.Abbreviations.Hinweisen,
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
		&c.PrintBestandsaufnahme,
		&m.PrintBestandsaufnahme,
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
