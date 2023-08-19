package sku_formats

import (
	"strings"

	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/hotel/sku"
)

func StringMetadatei(o sku.SkuLike) (str string) {
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
