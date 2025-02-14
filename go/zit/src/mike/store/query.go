package store

import (
	"slices"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

func (s *Store) QueryPrimitive(
	qg sku.PrimitiveQueryGroup,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	e := query.MakeExecutorPrimitive(
		qg,
		s.GetStreamIndex().ReadPrimitiveQuery,
		s.ReadOneInto,
	)

	if err = e.ExecuteTransacted(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryTransacted(
	qg *query.Query,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	var e query.Executor

	if e, err = s.makeQueryExecutor(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	var sk *sku.Transacted

	switch {
	case true:
		// TODO why does this not work with trying to read internal
		if sk, err = e.ExecuteExactlyOneExternalObject(false); err != nil {
			err = nil
			break
		}

		if err = f(sk); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = e.ExecuteTransacted(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryTransactedAsSkuType(
	qg *query.Query,
	f interfaces.FuncIter[sku.SkuType],
) (err error) {
	var e query.Executor

	if e, err = s.makeQueryExecutor(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.ExecuteTransactedAsSkuType(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (store *Store) QuerySkuType(
	queryGroup *query.Query,
	output interfaces.FuncIter[sku.SkuType],
) (err error) {
	var e query.Executor

	if e, err = store.makeQueryExecutor(queryGroup); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.ExecuteSkuType(output); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (store *Store) QueryExactlyOneExternal(
	queryGroup *query.Query,
) (sk *sku.Transacted, err error) {
	var executor query.Executor

	if executor, err = store.makeQueryExecutor(queryGroup); err != nil {
		err = errors.Wrap(err)
		return
	}

	if sk, err = executor.ExecuteExactlyOneExternalObject(true); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (store *Store) QueryExactlyOne(
	queryGroup *query.Query,
) (sk *sku.Transacted, err error) {
	var executor query.Executor

	if executor, err = store.makeQueryExecutor(queryGroup); err != nil {
		err = errors.Wrap(err)
		return
	}

	if sk, err = executor.ExecuteExactlyOne(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) MakeBlobShaBytesMap() (blobShaBytes map[sha.Bytes][]string, err error) {
	blobShaBytes = make(map[sha.Bytes][]string)
	var l sync.Mutex

	if err = s.QueryPrimitive(
		sku.MakePrimitiveQueryGroup(),
		func(sk *sku.Transacted) (err error) {
			l.Lock()
			defer l.Unlock()

			k := sk.Metadata.Blob.GetBytes()
			oids := blobShaBytes[k]
			oid := sk.ObjectId.String()
			loc, found := slices.BinarySearch(oids, oid)

			if found {
				return
			}

			oids = slices.Insert(oids, loc, oid)

			blobShaBytes[sk.Metadata.Blob.GetBytes()] = oids

			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
