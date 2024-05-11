package commands

import (
	"encoding/json"
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/november/umwelt"
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

func (c CheckinJson) DefaultGattungen() kennung.Gattung {
	return kennung.MakeGattung()
}

type TomlBookmark struct {
	Kennung   string
	Etiketten []string
	Url       string
}

func (c CheckinJson) Run(
	u *umwelt.Umwelt,
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
