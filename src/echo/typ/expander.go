package typ

import "github.com/friedenberg/zit/src/charlie/string_expansion"

type Expander = string_expansion.Expander[Typ, *Typ]

var (
	ExpanderRight Expander
	ExpanderAll   Expander
)

func init() {
	ExpanderRight = string_expansion.MakeExpanderRight[Typ, *Typ](`-`)
	ExpanderAll = string_expansion.MakeExpanderAll[Typ, *Typ](`-`)
}
