package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/bravo/todo"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

// TODO [radi/kof !task "add support for kasten in checkouts and external" project-2021-zit-features today zz-inbox]
func (s *Store) ReadExternal(
	qg query.GroupWithKasten,
	f schnittstellen.FuncIter[sku.CheckedOutLike],
) (err error) {
	switch qg.Kasten.String() {
	case "chrome":
		err = todo.Implement()

	default:
		if err = s.cwdFiles.ReadQuery(
			qg.Group,
			func(cofs *store_fs.CheckedOut) (err error) {
				return f(cofs)
			},
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
