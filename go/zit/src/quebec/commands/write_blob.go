package commands

import (
	"flag"
	"io"
	"sync/atomic"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/delta/script_value"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type WriteBlob struct {
	Check           bool
	AgeIdentity     age.Identity
	CompressionType immutable_config.CompressionType
	UtilityBefore   script_value.Utility
	UtilityAfter    script_value.Utility
}

func init() {
	registerCommand(
		"write-blob",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &WriteBlob{}

			f.BoolVar(&c.Check, "check", false, "only check if the object already exists")
			f.Var(&c.AgeIdentity, "age-identity", "")
			c.CompressionType.SetFlagSet(f)

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

func (c WriteBlob) RunWithRepo(u *repo_local.Repo, args ...string) {
	var failCount atomic.Uint32

	sawStdin := false

	var ag age.Age

	if err := ag.AddIdentity(c.AgeIdentity); err != nil {
		u.CancelWithErrorAndFormat(err, "age-identity: %q", &c.AgeIdentity)
	}

	for _, p := range args {
		switch {
		case sawStdin:
			ui.Err().Print("'-' passed in more than once. Ignoring")
			continue

		case p == "-":
			sawStdin = true
		}

		a := answer{Path: p}

		a.Sha, a.error = c.doOne(&ag, u.GetRepoLayout(), p)

		if a.error != nil {
			ui.Err().Printf("%s: (error: %q)", a.Path, a.error)
			failCount.Add(1)
			continue
		}

		hasBlob := u.GetRepoLayout().HasBlob(a.Sha)

		if hasBlob {
			if c.Check {
				u.GetUI().Printf("%s %s (already checked in)", a.GetShaLike(), a.Path)
			} else {
				u.GetUI().Printf("%s %s (checked in)", a.GetShaLike(), a.Path)
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
		u.CancelWithBadRequestf("untracked objects: %d", fc)
		return
	}
}

func (c WriteBlob) doOne(
	ag *age.Age,
	arf interfaces.BlobWriterFactory,
	p string,
) (sh interfaces.Sha, err error) {
	var rc io.ReadCloser

	o := repo_layout.FileReadOptions{
		Age:             ag,
		Path:            p,
		CompressionType: c.CompressionType,
	}

	if rc, err = repo_layout.NewFileReader(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, rc)

	var wc sha.WriteCloser

	if c.Check {
		wc = sha.MakeWriter(nil)
	} else {
		if wc, err = arf.BlobWriter(); err != nil {
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
