package commands

import (
	"encoding/json"
	"flag"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type ReadBlob struct{}

func init() {
	registerCommand(
		"read-akte",
		func(f *flag.FlagSet) Command {
			c := &ReadBlob{}

			return c
		},
	)
}

type readBlobEntry struct {
	Blob string `json:"akte"`
}

func (c ReadBlob) Run(u *env.Env, args ...string) (err error) {
	dec := json.NewDecoder(u.In())

	for {
		var entry readBlobEntry

		if err = dec.Decode(&entry); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		var sh *sha.Sha

		if sh, err = c.readOneAkte(u, entry); err != nil {
			err = errors.Wrap(err)
			return
		}

		ui.Debug().Print(sh)
	}

	return
}

func (ReadBlob) readOneAkte(u *env.Env, entry readBlobEntry) (sh *sha.Sha, err error) {
	var aw sha.WriteCloser

	if aw, err = u.GetFSHome().BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, aw)

	if _, err = io.Copy(aw, strings.NewReader(entry.Blob)); err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = sha.GetPool().Get()

	if err = sh.SetShaLike(aw.GetShaLike()); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
