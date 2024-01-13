package commands

import (
	"flag"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/oscar/umwelt"
)

type CatObjekte struct{}

func init() {
	registerCommand(
		"cat-objekte",
		func(_ *flag.FlagSet) Command {
			c := &CatObjekte{}

			return c
		},
	)
}

func (c CatObjekte) Run(
	u *umwelt.Umwelt,
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

		if err = c.akte(u, &sh, akteWriter); err == nil {
			continue
		}

		me.Add(err)

		// if sk, err = u.StoreUtil().GetVerzeichnisse().ReadOneShas(&sh); err == nil {
		// 	log.Out().Printf("%s", sk)
		// 	continue
		// }

		// me.Add(err)

		var rc sha.ReadCloser

		if rc, err = u.StoreUtil().ReaderFor(&sh); err == nil {
			if err = akteWriter(rc); err != nil {
				err = errors.Wrap(err)
				return
			}

			continue
		}

		me.Add(err)
		log.Err().Print(me)
	}

	return
}

func (c CatObjekte) akte(
	u *umwelt.Umwelt,
	sh *sha.Sha,
	akteWriter schnittstellen.FuncIter[io.ReadCloser],
) (err error) {
	var r io.ReadCloser

	if r, err = u.Standort().AkteReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = akteWriter(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
