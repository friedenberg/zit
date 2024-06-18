package sku_fmt

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

func StringMetadateiSansTai(o *sku.Transacted) (str string) {
	sb := &strings.Builder{}

	sb.WriteString(o.GetGattung().GetGattungString())

	sb.WriteString(" ")
	sb.WriteString(o.GetKennung().String())

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
				" ",
				m.GetEtiketten(),
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
