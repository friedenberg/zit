package commands

import (
	"flag"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type CatBlob struct{}

func init() {
	registerCommand(
		"cat-objekte",
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
	akteWriter := iter.MakeSyncSerializer(
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

		me := errors.MakeMulti()

		if err = c.blob(u, &sh, akteWriter); err == nil {
			continue
		}

		me.Add(err)

		// if sk, err = u.StoreUtil().GetVerzeichnisse().ReadOneShas(&sh); err == nil {
		// 	log.Out().Printf("%s", sk)
		// 	continue
		// }

		// me.Add(err)

		var rc sha.ReadCloser

		if rc, err = u.GetStore().ReaderFor(&sh); err == nil {
			if err = akteWriter(rc); err != nil {
				err = errors.Wrap(err)
				return
			}

			continue
		}

		me.Add(err)
		ui.Err().Print(me)
	}

	return
}

func (c CatBlob) blob(
	u *env.Env,
	sh *sha.Sha,
	akteWriter interfaces.FuncIter[io.ReadCloser],
) (err error) {
	var r io.ReadCloser

	if r, err = u.GetFSHome().BlobReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = akteWriter(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
