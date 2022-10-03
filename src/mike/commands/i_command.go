package commands

import "github.com/friedenberg/zit/src/kilo/umwelt"

type Command interface {
	Run(*umwelt.Umwelt, ...string) error
}
