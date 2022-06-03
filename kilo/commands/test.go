package commands

import (
	"bufio"
	"flag"
	"io"

	"github.com/friedenberg/zit/alfa/errors"
)

type Test struct {
}

func init() {
	registerCommand(
		"test",
		func(f *flag.FlagSet) Command {
			c := &Test{}

			return c
		},
	)
}

func (c Test) Run(u _Umwelt, args ...string) (err error) {
	var a _Age

	if a, err = u.Age(); err != nil {
		err = errors.Error(err)
		return
	}

	var e _Etiketten

	if e, err = _NewEtiketten(_Konfig{}, a, u.DirZit("Etiketten")); err != nil {
		err = errors.Error(err)
		return
	}

	r := bufio.NewReader(u.In)

	for {
		var l string

		l, err = r.ReadString('\n')

		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			err = errors.Error(err)
			return
		}

		e.AddString(l)
	}

	if err = e.Flush(); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
