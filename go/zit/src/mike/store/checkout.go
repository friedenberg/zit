package store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/store_fs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/query"
)

func (s *Store) DeleteCheckout(col sku.CheckedOutLike) (err error) {
	switch cot := col.(type) {
	default:
		err = errors.Errorf("unsupported checkout: %T, %s", cot, cot)
		return

	case *store_fs.CheckedOut:
		if err = s.GetCwdFiles().Delete(&cot.External); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (s *Store) CheckoutQueryFS(
	options checkout_options.Options,
	qg *query.Group,
	f schnittstellen.FuncIter[*store_fs.CheckedOut],
) (err error) {
	if err = s.QueryWithCwd(
		qg,
		func(t *sku.Transacted) (err error) {
			var cop *store_fs.CheckedOut

			if cop, err = s.CheckoutOneFS(options, t); err != nil {
				err = errors.Wrap(err)
				return
			}

			cop.DetermineState(true)

			if err = s.checkedOutLogPrinter(cop); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = f(cop); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

// TODO [radi/kof !task "add support for kasten in checkouts and external" project-2021-zit-features today zz-inbox]
func (s *Store) CheckoutOneFS(
	options checkout_options.Options,
	sz *sku.Transacted,
) (cz *store_fs.CheckedOut, err error) {
	if cz, err = s.cwdFiles.CheckoutOneFS(options, sz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) CheckoutOne(
	options checkout_options.Options,
	sz *sku.Transacted,
) (cz sku.CheckedOutLike, err error) {
	if cz, err = s.cwdFiles.CheckoutOneFS(options, sz); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
