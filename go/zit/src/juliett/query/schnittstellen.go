package query

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
)

type VirtualStore interface {
	// sku.ExternalStore
	Initialize() error
	Flush() error
	// TODO add objekte mode
	// CommitTransacted(kinder, mutter *sku.Transacted) error
	// ModifySku(*sku.Transacted) error
	// Query(*Group, schnittstellen.FuncIter[*sku.Transacted]) error
	// sku.Queryable
}

type Reducer interface {
	Reduce(*Builder) error
}

type Cwd interface {
	sku.Query
	GetCwdFDs() fd.Set
	GetKennungForFD(*fd.FD) (*kennung.Kennung2, error)
}
