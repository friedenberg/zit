package remote_transfers

import (
	"io"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

type boundReader struct {
	lContinue sync.Locker
	lr        *io.LimitedReader
	io.Reader
}

func makeBoundReader(
	in io.Reader,
	lContinue sync.Locker,
	n int64,
) sha.ReadCloser {
	lr := &io.LimitedReader{
		R: in,
		N: n,
	}

	r := &boundReader{
		lContinue: lContinue,
		lr:        lr,
		Reader:    lr,
	}

	return sha.MakeReadCloser(r)
}

func (r *boundReader) Close() (err error) {
	ui.Log().Printf("closing bound reader")
	defer ui.Log().Printf("did close bound")

	defer r.lContinue.Unlock()

	if _, err = io.Copy(io.Discard, r.Reader); err != nil {
		if errors.IsEOF(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
