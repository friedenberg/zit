package typ

import "github.com/friedenberg/zit/src/charlie/kennung"

type Expander = kennung.Expander[Kennung, *Kennung]

var (
	ExpanderRight Expander
	ExpanderAll   Expander
)

func init() {
	ExpanderRight = kennung.MakeExpanderRight[Kennung, *Kennung](`-`)
	ExpanderAll = kennung.MakeExpanderAll[Kennung, *Kennung](`-`)
}
