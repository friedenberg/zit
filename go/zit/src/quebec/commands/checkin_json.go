package commands

import (
	"encoding/json"
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/november/env"
)

type CheckinJson struct{}

func init() {
	registerCommand(
		"checkin-json",
		func(f *flag.FlagSet) Command {
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
	u *env.Local,
	args ...string,
) (err error) {
	dec := json.NewDecoder(u.In())

	for {
		var entry TomlBookmark

		if err = dec.Decode(&entry); err != nil {
			if errors.IsEOF(err) {
				err = nil
				break
			} else {
				err = errors.Wrap(err)
				return
			}
		}

		ui.Debug().Print(entry)
	}

	return
}
