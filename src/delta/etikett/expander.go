package etikett

import kennung "github.com/friedenberg/zit/src/charlie/kennung"

type Expander = kennung.Expander[Etikett, *Etikett]

var (
	ExpanderRight Expander
	ExpanderAll   Expander
)

func init() {
	ExpanderRight = kennung.MakeExpanderRight[Etikett, *Etikett](`-`)
	ExpanderAll = kennung.MakeExpanderAll[Etikett, *Etikett](`-`)
}
