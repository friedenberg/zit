package commands

import (
	"encoding/json"
	"flag"
	"log"
)

type Log struct {
}

func init() {
	registerCommand(
		"log",
		func(f *flag.FlagSet) Command {
			c := &Log{}

			return commandWithZettels{c}
		},
	)
}

func (c Log) RunWithZettels(u _Umwelt, zs _Zettels, args ...string) (err error) {
	var rawId string

	switch len(args) {

	case 0:
		err = _Errorf("hinweis or zettel sha required")
		return

	default:
		_Errf("ignoring extra arguments: %q\n", args[1:])

		fallthrough

	case 1:
		rawId = args[0]

	}

	var id _Id

	if id, err = c.getIdFromArg(rawId); err != nil {
		err = _Error(err)
		return
	}

	var chain _ZettelsChain

	if chain, err = zs.AllInChain(id); err != nil {
		err = _Error(err)
		return
	}

	b, err := json.Marshal(chain)

	if err != nil {
		log.Print(err)
	} else {
		_Out(string(b))
	}

	return
}

func (c Log) getIdFromArg(arg string) (id _Id, err error) {
	var sha _Sha

	if err = sha.Set(arg); err == nil {
		id = sha
		return
	}

	hinweis := _HinweisNewEmpty()

	if err = hinweis.Set(arg); err == nil {
		id = hinweis
		return
	}

	err = _Errorf("incorrect format for id: '%s'", arg)

	return
}
