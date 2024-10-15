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
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type WriteBlob struct {
	Check           bool
	AgeIdentity     age.Identity
	CompressionType immutable_config.CompressionType
}

func init() {
	registerCommand(
		"write-blob",
		func(f *flag.FlagSet) Command {
			c := &WriteBlob{}

			f.BoolVar(&c.Check, "check", false, "only check if the object already exists")
			f.Var(&c.AgeIdentity, "age-identity", "")
			c.CompressionType.AddToFlagSet(f)

			return c
		},
	)
}

type answer struct {
	error
	interfaces.Sha
	Path string
}

func (c WriteBlob) Run(u *env.Env, args ...string) (err error) {
	// wg := &sync.WaitGroup{}
	// wg.Add(len(args))

	var failCount atomic.Uint32

	// chShas := make(chan answer)

	sawStdin := false

	var ag age.Age

	if err = ag.AddIdentity(c.AgeIdentity); err != nil {
		err = errors.Wrapf(err, "age-identity: %q", &c.AgeIdentity)
		return
	}

	// go func() {
	for _, p := range args {
		switch {
		case sawStdin:
			ui.Err().Print("'-' passed in more than once. Ignoring")
			continue

		case p == "-":
			sawStdin = true
		}

		a := answer{Path: p}

		a.Sha, a.error = c.doOne(&ag, u.GetFSHome(), p)

		if a.error != nil {
			ui.Err().Printf("%s: %s", a.Path, a.error)
			failCount.Add(1)
			continue
		}

		hasBlob := u.GetFSHome().HasBlob(u.GetConfig().GetStoreVersion(), a.Sha)

		if hasBlob {
			ui.Out().Printf("%s %s (checked in)", a.GetShaLike(), a.Path)
		} else {
			ui.Err().Printf("%s %s (untracked)", a.GetShaLike(), a.Path)

			if c.Check {
				failCount.Add(1)
			}
		}
	}
	// }()

	// for _, p := range args {
	// 	switch {
	// 	case sawStdin:
	// 		ui.Err().Print("'-' passed in more than once. Ignoring")
	// 		continue

	// 	case p == "-":
	// 		sawStdin = true
	// 	}

	// 	go func() {
	// 		defer wg.Done()

	// 		a := answer{Path: p}

	// 		a.Sha, a.error = c.doOne(&ag, u.GetFSHome(), p)

	// 		chShas <- a
	// 	}()
	// }

	// wg.Wait()

	fc := failCount.Load()

	if fc > 0 {
		err = errors.BadRequestf("untracked objects: %d", fc)
		return
	}

	return
}

func (c WriteBlob) doOne(
	ag *age.Age,
	arf interfaces.BlobWriterFactory,
	p string,
) (sh interfaces.Sha, err error) {
	var rc io.ReadCloser

	o := fs_home.FileReadOptions{
		Age:             ag,
		Path:            p,
		CompressionType: c.CompressionType,
	}

	if rc, err = fs_home.NewFileReader(o); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, rc)

	var wc sha.WriteCloser

	if c.Check {
		wc = sha.MakeWriter(nil)
	} else {
		if wc, err = arf.BlobWriter(); err != nil {
			ui.Debug().Print(err)
			err = errors.Wrap(err)
			return
		}
	}

	defer errors.DeferredCloser(&err, wc)

	if _, err = io.Copy(wc, rc); err != nil {
		ui.Debug().Print(err)
		err = errors.Wrap(err)
		return
	}

	sh = wc.GetShaLike()

	return
}
