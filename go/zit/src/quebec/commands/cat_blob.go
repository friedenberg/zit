package commands

import (
	"flag"
	"io"
	"os/exec"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/delim_io"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type CatBlob struct {
	Utility   script_value.Utility
	PrefixSha bool
}

func init() {
	registerCommand(
		"cat-blob",
		func(f *flag.FlagSet) Command {
			c := &CatBlob{}

			f.Var(&c.Utility, "utility", "")
			f.BoolVar(&c.PrefixSha, "prefix-sha", false, "")

			return c
		},
	)
}

type shaWithReadCloser struct {
	Sha        *sha.Sha
	ReadCloser io.ReadCloser
}

func (c CatBlob) makeBlobWriter(u *env.Local) interfaces.FuncIter[shaWithReadCloser] {
	if c.Utility.IsEmpty() {
		return quiter.MakeSyncSerializer(
			func(rc shaWithReadCloser) (err error) {
				if err = c.copy(u, rc); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		)
	} else {
		return quiter.MakeSyncSerializer(
			func(rc shaWithReadCloser) (err error) {
				defer errors.DeferredCloser(&err, rc.ReadCloser)

				cmd := exec.Command(c.Utility.Head(), c.Utility.Tail()...)
				cmd.Stdin = rc.ReadCloser

				var out io.ReadCloser

				if out, err = cmd.StdoutPipe(); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = cmd.Start(); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = c.copy(
					u,
					shaWithReadCloser{
						Sha:        rc.Sha,
						ReadCloser: out,
					},
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				if err = cmd.Wait(); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			},
		)
	}
}

func (c CatBlob) Run(
	u *env.Local,
	args ...string,
) (err error) {
	blobWriter := c.makeBlobWriter(u)

	for _, v := range args {
		var sh sha.Sha

		if err = sh.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = c.blob(u, &sh, blobWriter); err != nil {
			ui.Err().Print(err)
			err = nil
		}
	}

	return
}

func (c CatBlob) copy(
	u *env.Local,
	rc shaWithReadCloser,
) (err error) {
	defer errors.DeferredCloser(&err, rc.ReadCloser)

	if c.PrefixSha {
		if _, err = delim_io.CopyWithPrefixOnDelim(
			'\n',
			rc.Sha.GetShaLike().GetShaString(),
			u.Out(),
			rc.ReadCloser,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if _, err = io.Copy(u.Out(), rc.ReadCloser); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c CatBlob) blob(
	u *env.Local,
	sh *sha.Sha,
	blobWriter interfaces.FuncIter[shaWithReadCloser],
) (err error) {
	var r sha.ReadCloser

	if r, err = u.GetDirectoryLayout().BlobReader(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = blobWriter(shaWithReadCloser{Sha: sh, ReadCloser: r}); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
