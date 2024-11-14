package store

import (
	"slices"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

func (s *Store) QueryPrimitive(
	qg sku.PrimitiveQueryGroup,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	var e query.ExecutorPrimitive

	if e, err = s.MakeQueryExecutorPrimitive(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.ExecuteTransacted(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryTransacted(
	qg *query.Group,
	f interfaces.FuncIter[*sku.Transacted],
) (err error) {
	var e query.Executor

	if e, err = s.MakeQueryExecutor(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.ExecuteTransacted(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QueryTransactedAsSkuType(
	qg *query.Group,
	f interfaces.FuncIter[sku.SkuType],
) (err error) {
	var e query.Executor

	if e, err = s.MakeQueryExecutor(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.ExecuteTransactedAsSkuType(f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (s *Store) QuerySkuType(
	qg *query.Group,
	f interfaces.FuncIter[sku.SkuType],
) (err error) {
	var e query.Executor

	if e, err = s.MakeQueryExecutor(qg); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.ExecuteSkuType(f); err != nil {
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
