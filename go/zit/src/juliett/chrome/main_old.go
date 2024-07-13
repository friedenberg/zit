package chrome

import (
	"net/url"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

func (c *Store) CommitTransacted(kinder, mutter *sku.Transacted) (err error) {
	if c.konfig.DryRun {
		return
	}

	var dt diff

	if dt, err = c.getDiff(kinder, mutter); err != nil {
		err = errors.Wrap(err)
		return
	}

	if dt.diffType == diffTypeIgnore {
		return
	}

	var u *url.URL

	if u, err = c.getUrl(kinder); err != nil {
		err = errors.Wrap(err)
		return
	}

	if _, ok := c.urls[*u]; !ok {
		// TODO fetch previous URL
		return
	}

	switch dt.diffType {
	case diffTypeDelete:
		ui.Debug().Print("deleted", "TODO add to dedicated printer", kinder)
		c.removed[*u] = struct{}{}

	default:
		ui.Debug().Print("TODO not implemented", dt, kinder, mutter)
	}

	return
}

func (c *Store) Query(
	qg *query.Group,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	return
}

func (c *Store) ContainsSku(sk *sku.Transacted) bool {
	return false
}
