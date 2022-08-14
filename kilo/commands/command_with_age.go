package commands

import (
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/delta/age"
	"github.com/friedenberg/zit/delta/umwelt"
)

type CommandWithAge interface {
	RunWithAge(*umwelt.Umwelt, age.Age, ...string) error
}

type commandWithAge struct {
	CommandWithAge
}

func (c commandWithAge) Run(u *umwelt.Umwelt, args ...string) (err error) {
	var age age.Age

	if age, err = u.Age(); err != nil {
		err = errors.Error(err)
		return
	}

	if err = c.RunWithAge(u, age, args...); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
