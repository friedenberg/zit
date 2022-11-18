package etikett

import "github.com/friedenberg/zit/src/string_expansion"

type Expander = string_expansion.Expander[Etikett, *Etikett]

var (
	ExpanderRight Expander
	ExpanderAll   Expander
)

func init() {
	ExpanderRight = string_expansion.MakeExpanderRight[Etikett, *Etikett](`-`)
	ExpanderAll = string_expansion.MakeExpanderAll[Etikett, *Etikett](`-`)
}
