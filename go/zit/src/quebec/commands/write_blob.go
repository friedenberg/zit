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
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

type WriteBlob struct {
	Check         bool
	UtilityBefore script_value.Utility
	UtilityAfter  script_value.Utility
}

func init() {
	registerCommand(
		"write-blob",
		func(f *flag.FlagSet) WithBlobStore {
			c := &WriteBlob{}

			f.BoolVar(&c.Check, "check", false, "only check if the object already exists")

			f.Var(&c.UtilityBefore, "utility-before", "")
			f.Var(&c.UtilityAfter, "utility-after", "")

			return c
		},
	)
}

type answer struct {
	error
	interfaces.Sha
	Path string
}

func (c WriteBlob) RunWithBlobStore(
	blobStore command_components.BlobStoreWithEnv,
	args ...string,
) {
	var failCount atomic.Uint32

	sawStdin := false

	for _, p := range args {
		switch {
		case sawStdin:
			ui.Err().Print("'-' passed in more than once. Ignoring")
			continue

		case p == "-":
			sawStdin = true
		}

		a := answer{Path: p}

		a.Sha, a.error = c.doOne(blobStore, p)

		if a.error != nil {
			blobStore.GetErr().Printf("%s: (error: %q)", a.Path, a.error)
			failCount.Add(1)
			continue
		}

		hasBlob := blobStore.HasBlob(a.Sha)

		if hasBlob {
			if c.Check {
				blobStore.GetUI().Printf("%s %s (already checked in)", a.GetShaLike(), a.Path)
			} else {
				blobStore.GetUI().Printf("%s %s (checked in)", a.GetShaLike(), a.Path)
			}
		} else {
			ui.Err().Printf("%s %s (untracked)", a.GetShaLike(), a.Path)

			if c.Check {
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
