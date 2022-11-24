package typ

import "github.com/friedenberg/zit/src/delta/kennung"

type Expander = kennung.Expander

var (
	ExpanderRight Expander
	ExpanderAll   Expander
)

func init() {
	ExpanderRight = kennung.MakeExpanderRight(`-`)
	ExpanderAll = kennung.MakeExpanderAll(`-`)
}
