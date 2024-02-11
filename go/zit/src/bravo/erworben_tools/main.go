package erworben_tools

import (
	"flag"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"github.com/google/shlex"
)

type Tools struct {
	Merge []string `toml:"merge"`
}

func (c *Tools) AddToFlags(f *flag.FlagSet) {
	f.Func(
		"merge-tool",
		"utility to launch for merge conflict resolution",
		func(value string) (err error) {
			if c.Merge, err = shlex.Split(value); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		},
	)
}
