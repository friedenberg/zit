package todo

import "github.com/friedenberg/zit/src/alfa/errors"

var errNotImplemented = errors.New("not implemented")

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

func Optimize() {
	errors.TodoP0("start logging this")
}

func Implement() (err error) {
	errors.TodoP0("start logging this")
	return errors.WrapN(1, errNotImplemented)
}

func Remove() {
	errors.TodoP0("start logging this")
}