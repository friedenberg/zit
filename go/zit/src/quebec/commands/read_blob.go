package commands

import (
	"encoding/json"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/papa/command_components"
)

func init() {
	registerCommand("read-blob", &ReadBlob{})
}

type ReadBlob struct {
	command_components.RepoLayout
}

type readBlobEntry struct {
	Blob string `json:"blob"`
}

func (c ReadBlob) Run(dep command.Dep) {
	repoLayout := c.MakeRepoLayout(dep, false)

	dec := json.NewDecoder(repoLayout.GetInFile())

	for {
		var entry readBlobEntry

		if err := dec.Decode(&entry); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				repoLayout.CancelWithError(err)
			}

			return
		}

		{
			var err error

			if _, err = c.readOneBlob(repoLayout, entry); err != nil {
				repoLayout.CancelWithError(err)
			}
		}
	}
}

func (ReadBlob) readOneBlob(
	repoLayout repo_layout.Layout,
	entry readBlobEntry,
) (sh *sha.Sha, err error) {
	var aw sha.WriteCloser

	if aw, err = repoLayout.BlobWriter(); err != nil {
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
