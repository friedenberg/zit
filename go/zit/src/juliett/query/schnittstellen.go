package query

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type Reducer interface {
	Reduce(*Builder) error
}

type Kasten interface {
	sku.Query
	GetCwdFDs() fd.Set
	GetKennungForFD(*fd.FD) (*kennung.Kennung2, error)
}
