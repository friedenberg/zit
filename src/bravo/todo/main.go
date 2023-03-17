package todo

import "github.com/friedenberg/zit/src/alfa/errors"

func Change(_ string) {
	errors.TodoP0("start logging this")
}

func Decide(_ string) {
	errors.TodoP0("start logging this")
}

func Refactor(_ string) {
	errors.TodoP0("start logging this")
}

func Parallelize() {
	errors.TodoP0("start logging this")
}

func Implement() (err error) {
	errors.TodoP0("start logging this")
	return errors.Implement()
}

func Remove() {
	errors.TodoP0("start logging this")
}
