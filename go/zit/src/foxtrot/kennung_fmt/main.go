package kennung_fmt

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
)

type Aligned interface {
	schnittstellen.StringFormatWriter[*kennung.Kennung2]
	SetMaxKopfUndSchwanz(kop, schwanz int)
}
