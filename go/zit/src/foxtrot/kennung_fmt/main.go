package kennung_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

type Aligned interface {
	interfaces.StringFormatWriter[*kennung.Id]
	SetMaxKopfUndSchwanz(kop, schwanz int)
}
