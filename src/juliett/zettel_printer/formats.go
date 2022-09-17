package zettel_printer

import (
	"io"
	"strings"
)

type Zettelish struct {
	zettelish
}

type zettelish struct {
	newZettelShaSyntax   bool
	includeBezeichnungen bool
	includeTyp           bool
	Id, Rev, Typ, Bez    io.WriterTo
}

func (zp *Printer) MakeZettelish() Zettelish {
	return Zettelish{
		zettelish{
			includeTyp:           zp.includeTyp,
			newZettelShaSyntax:   zp.newZettelShaSyntax,
			includeBezeichnungen: zp.includeBezeichnungen,
		},
	}
}

func (zi Zettelish) IdString(v string) Zettelish {
	zi.zettelish.Id = stringWriterTo(v)
	return zi
}

func (zi Zettelish) Id(v io.WriterTo) Zettelish {
	zi.zettelish.Id = v
	return zi
}

func (zi Zettelish) RevString(v string) Zettelish {
	zi.zettelish.Rev = stringWriterTo(v)
	return zi
}

func (zi Zettelish) Rev(v io.WriterTo) Zettelish {
	zi.zettelish.Rev = v
	return zi
}

func (zi Zettelish) TypString(v string) Zettelish {
	zi.zettelish.Typ = stringWriterTo(v)
	return zi
}

func (zi Zettelish) Typ(v io.WriterTo) Zettelish {
	zi.zettelish.Typ = v
	return zi
}

func (zi Zettelish) BezString(v string) Zettelish {
	zi.zettelish.Bez = stringWriterTo(v)
	return zi
}

func (zi Zettelish) Bez(v io.WriterTo) Zettelish {
	zi.zettelish.Bez = v
	return zi
}

func (zi zettelish) String() string {
	sb := &strings.Builder{}

	sb.WriteString("[")

	var err error

	if _, err = zi.Id.WriteTo(sb); err != nil {
		return err.Error()
	}

	if zi.newZettelShaSyntax && zi.Rev != nil {
		sb.WriteString("@")

		if _, err = zi.Rev.WriteTo(sb); err != nil {
			return err.Error()
		}
	}

	if zi.includeTyp && zi.Typ != nil {
		sb.WriteString(" !")

		if _, err = zi.Typ.WriteTo(sb); err != nil {
			return err.Error()
		}
	}

	if zi.includeBezeichnungen && zi.Bez != nil {
		//TODO use bez descriptors instead
		sb.WriteString(" ")

		if _, err = zi.Bez.WriteTo(sb); err != nil {
			return err.Error()
		}
	}

	sb.WriteString("]")

	return sb.String()
}
