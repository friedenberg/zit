package typ

import "github.com/friedenberg/zit/src/charlie/string_expansion"

type Expander = string_expansion.Expander[Kennung, *Kennung]

var (
	ExpanderRight Expander
	ExpanderAll   Expander
)

func init() {
	ExpanderRight = string_expansion.MakeExpanderRight[Kennung, *Kennung](`-`)
	ExpanderAll = string_expansion.MakeExpanderAll[Kennung, *Kennung](`-`)
}
