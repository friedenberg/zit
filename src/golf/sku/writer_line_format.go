package sku

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/delta/format"
	"github.com/friedenberg/zit/src/delta/kennung"
)

func StringMetadateiSansTai(o SkuLike) (str string) {
	sb := &strings.Builder{}

	sb.WriteString(o.GetGattung().GetGattungString())

	sb.WriteString(" ")
	sb.WriteString(o.GetKennungLike().String())

	sb.WriteString(" ")
	sb.WriteString(o.GetAkteSha().String())

	m := o.GetMetadatei()

	t := m.GetTyp()

	if !t.IsEmpty() {
		sb.WriteString(" ")
		sb.WriteString(kennung.FormattedString(m.GetTyp()))
	}

	es := m.GetEtiketten()

	if es.Len() > 0 {
		sb.WriteString(" ")
		sb.WriteString(
			iter.StringDelimiterSeparated[kennung.Etikett](
				m.GetEtiketten(),
				" ",
			),
		)
	}

	b := m.GetBezeichnung()

	if !b.IsEmpty() {
		sb.WriteString(" ")
		sb.WriteString("\"" + b.String() + "\"")
	}

	return sb.String()
}

func StringMetadatei(o SkuLike) (str string) {
	sb := &strings.Builder{}

	sb.WriteString(o.GetTai().String())

	sb.WriteString(" ")
	sb.WriteString(o.GetGattung().GetGattungString())

	sb.WriteString(" ")
	sb.WriteString(o.GetKennungLike().String())

	sb.WriteString(" ")
	sb.WriteString(o.GetAkteSha().String())

	m := o.GetMetadatei()

	t := m.GetTyp()

	if !t.IsEmpty() {
		sb.WriteString(" ")
		sb.WriteString(kennung.FormattedString(m.GetTyp()))
	}

	es := m.GetEtiketten()

	if es.Len() > 0 {
		sb.WriteString(" ")
		sb.WriteString(
			iter.StringDelimiterSeparated[kennung.Etikett](
				m.GetEtiketten(),
				" ",
			),
		)
	}

	b := m.GetBezeichnung()

	if !b.IsEmpty() {
		sb.WriteString(" ")
		sb.WriteString("\"" + b.String() + "\"")
	}

	return sb.String()
}

func String(o SkuLike) (str string) {
	str = fmt.Sprintf(
		"%s %s %s %s %s",
		o.GetTai(),
		o.GetGattung(),
		o.GetKennungLike(),
		o.GetObjekteSha(),
		o.GetAkteSha(),
	)

	return
}

func MakeWriterLineFormat(
	lf *format.LineWriter,
) schnittstellen.FuncIter[SkuLike] {
	return func(o SkuLike) (err error) {
		lf.WriteFormat("%s", o)

		return
	}
}
