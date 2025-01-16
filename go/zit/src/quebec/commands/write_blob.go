package commands

import (
	"flag"
	"io"
	"sync/atomic"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	command.Register("write-blob", &WriteBlob{})
}

type WriteBlob struct {
	command_components.BlobStoreLocal

	Check         bool
	UtilityBefore script_value.Utility
	UtilityAfter  script_value.Utility
}

func (cmd *WriteBlob) SetFlagSet(f *flag.FlagSet) {
	f.BoolVar(&cmd.Check, "check", false, "only check if the object already exists")

	f.Var(&cmd.UtilityBefore, "utility-before", "")
	f.Var(&cmd.UtilityAfter, "utility-after", "")
}

type answer struct {
	error
	interfaces.Sha
	Path string
}

func (cmd WriteBlob) Run(
	dep command.Dep,
) {
	blobStore := cmd.MakeBlobStoreLocal(
		dep.Context,
		dep.Config,
		env.Options{},
		local_working_copy.OptionsEmpty,
	)

	var failCount atomic.Uint32

	sawStdin := false

	for _, p := range dep.Args() {
		switch {
		case sawStdin:
			ui.Err().Print("'-' passed in more than once. Ignoring")
			continue

		case p == "-":
			sawStdin = true
		}

		a := answer{Path: p}

		a.Sha, a.error = cmd.doOne(blobStore, p)

		if a.error != nil {
			blobStore.GetErr().Printf("%s: (error: %q)", a.Path, a.error)
			failCount.Add(1)
			continue
		}

		hasBlob := blobStore.HasBlob(a.Sha)

		if hasBlob {
			if cmd.Check {
				blobStore.GetUI().Printf("%s %s (already checked in)", a.GetShaLike(), a.Path)
			} else {
				blobStore.GetUI().Printf("%s %s (checked in)", a.GetShaLike(), a.Path)
			}
		} else {
			ui.Err().Printf("%s %s (untracked)", a.GetShaLike(), a.Path)

			if cmd.Check {
				failCount.Add(1)
			}
		}
	}

	fc := failCount.Load()

	if fc > 0 {
		blobStore.CancelWithBadRequestf("untracked objects: %d", fc)
		return
	}
}

func (c WriteBlob) doOne(
	blobStore command_components.BlobStoreWithEnv,
	p string,
) (sh interfaces.Sha, err error) {
	var rc io.ReadCloser

	o := dir_layout.FileReadOptions{
		Path: p,
	}

	if rc, err = dir_layout.NewFileReader(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, rc)

	var wc sha.WriteCloser

	if c.Check {
		wc = sha.MakeWriter(nil)
	} else {
		if wc, err = blobStore.BlobWriter(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	defer errors.DeferredCloser(&err, wc)

	if _, err = io.Copy(wc, rc); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = wc.GetShaLike()

	return
}
