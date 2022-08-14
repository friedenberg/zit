package commands

import (
	"bufio"
	"flag"
	"io"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/delta/age"
	"github.com/friedenberg/zit/delta/konfig"
	"github.com/friedenberg/zit/echo/umwelt"
	"github.com/friedenberg/zit/golf/etiketten"
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

func (c Test) Run(u *umwelt.Umwelt, args ...string) (err error) {
	var a age.Age

	if a, err = u.Age(); err != nil {
		err = errors.Error(err)
		return
	}

	var e etiketten.Etiketten

	if e, err = etiketten.New(konfig.Konfig{}, a, u.DirZit("Etiketten")); err != nil {
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
