package commands

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/delta/kennung"
	"github.com/friedenberg/zit/src/november/umwelt"
)

type CommandWithHinweis interface {
	RunWithHinweis(store *umwelt.Umwelt, h kennung.Hinweis) error
}

type commandWithHinweis struct {
	CommandWithHinweis
}

func (c commandWithHinweis) Complete(u *umwelt.Umwelt, args ...string) (err error) {
	err = errors.Implement()
	return
}

func (c commandWithHinweis) Run(u *umwelt.Umwelt, args ...string) (err error) {
	errors.TodoP1("add metaid type to kennung and support for sigils")
	if len(args) != 1 {
		err = errors.Normalf("only one hinweis is accepted")
		return
	}

	var h kennung.Hinweis
	v := args[0]

	if h, err = u.StoreObjekten().GetAbbrStore().ExpandHinweisString(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = c.RunWithHinweis(u, h); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
