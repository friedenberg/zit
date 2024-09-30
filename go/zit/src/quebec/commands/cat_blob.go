package commands

import (
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type CatBlob struct{}

func init() {
	registerCommand(
		"cat-blob",
		func(_ *flag.FlagSet) Command {
			c := &CatBlob{}

			return c
		},
	)
}

func (c CatBlob) Run(
	u *env.Env,
	args ...string,
) (err error) {
	blobWriter := quiter.MakeSyncSerializer(
		func(rc io.ReadCloser) (err error) {
			defer errors.DeferredCloser(&err, rc)

			if _, err = io.Copy(u.Out(), rc); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)

	for _, v := range args {
		var sh sha.Sha

		if err = sh.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = c.blob(u, &sh, blobWriter); err == nil {
			continue
		}

		var rc sha.ReadCloser

		if rc, err = u.GetStore().ReaderFor(&sh); err == nil {
			if err = blobWriter(rc); err != nil {
				ui.Err().Print(err)
				err = nil
				continue
			}

			continue
		}

		ui.Err().Print(err)
	}

	return
}

func (c CatBlob) blob(
	u *env.Env,
	sh *sha.Sha,
	blobWriter interfaces.FuncIter[io.ReadCloser],
) (err error) {
	var r io.ReadCloser

	if r, err = u.GetFSHome().BlobReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = blobWriter(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
