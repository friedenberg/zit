package remote_http

import (
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/tridex"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
)

type serverBlobCache struct {
	ui             fd.Std
	localBlobStore interfaces.LocalBlobStore
	shas           interfaces.MutableTridex
	init           sync.Once
}

func (serverBlobCache *serverBlobCache) populate() (err error) {
	serverBlobCache.shas = tridex.Make()

	{
		count := 0

		for sh, errIter := range serverBlobCache.localBlobStore.AllBlobs() {
			if errIter != nil {
				err = errors.Wrap(errIter)
				return
			}

			serverBlobCache.shas.Add(sh.String())
			count++
		}

		ui.Log().Printf("have blobs: %d", count)
	}

	return
}

func (serverBlobCache *serverBlobCache) HasBlob(
	blobSha interfaces.Sha,
) (ok bool, err error) {
	serverBlobCache.init.Do(
		func() {
			if err = serverBlobCache.populate(); err != nil {
				err = errors.Wrap(err)
			}
		},
	)

	if err != nil {
		return
	}

	if serverBlobCache.shas.ContainsExpansion(blobSha.String()) {
		ok = true
		return
	}

	return
}
