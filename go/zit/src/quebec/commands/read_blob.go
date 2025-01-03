package commands

import (
	"encoding/json"
	"flag"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/november/repo_local"
)

type ReadBlob struct{}

func init() {
	registerCommand(
		"read-blob",
		func(f *flag.FlagSet) CommandWithRepo {
			c := &ReadBlob{}

			return c
		},
	)
}

type readBlobEntry struct {
	Blob string `json:"blob"`
}

func (c ReadBlob) RunWithRepo(u *repo_local.Repo, args ...string) {
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

func (ReadBlob) readOneBlob(u *repo_local.Repo, entry readBlobEntry) (sh *sha.Sha, err error) {
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
