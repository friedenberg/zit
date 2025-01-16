package commands

import (
	"encoding/json"
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type CheckinJson struct{}

func init() {
	registerCommandOld(
		"checkin-json",
		func(f *flag.FlagSet) WithLocalWorkingCopy {
			c := &CheckinJson{}

			return c
		},
	)
}

func (c CheckinJson) DefaultGenres() ids.Genre {
	return ids.MakeGenre()
}

type TomlBookmark struct {
	ObjectId string
	Tags     []string
	Url      string
}

func (c CheckinJson) Run(
	u *local_working_copy.Repo,
	args ...string,
) {
	dec := json.NewDecoder(u.GetInFile())

	for {
		var entry TomlBookmark

		if err := dec.Decode(&entry); err != nil {
			if errors.IsEOF(err) {
				err = nil
				break
			} else {
				u.CancelWithError(err)
			}
		}
	}
}
