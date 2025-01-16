package commands

import (
	"encoding/json"
	"flag"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type ReadBlob struct{}

func init() {
	registerCommandOld(
		"read-blob",
		func(f *flag.FlagSet) WithLocalWorkingCopy {
			c := &ReadBlob{}

			return c
		},
	)
}

type readBlobEntry struct {
	Blob string `json:"blob"`
}

func (c ReadBlob) Run(u *local_working_copy.Repo, args ...string) {
	dec := json.NewDecoder(u.GetInFile())

	for {
		var entry readBlobEntry

		if err := dec.Decode(&entry); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				u.CancelWithError(err)
			}

			return
		}

		{
			var err error

			if _, err = c.readOneBlob(u, entry); err != nil {
				u.CancelWithError(err)
			}
		}
	}
}

func (ReadBlob) readOneBlob(u *local_working_copy.Repo, entry readBlobEntry) (sh *sha.Sha, err error) {
	var aw sha.WriteCloser

	if aw, err = u.GetRepoLayout().BlobWriter(); err != nil {
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
