package errors

import (
	"fmt"
	"os"
)

var (
	todo todoPrinter
)

func init() {
	todo = todoPrinter{
		f:             os.Stderr,
		includesStack: true,
	}
}

func SetTodoOn() {
	todo.on = true
}

type todoPrinter printer

//go:generate stringer -type=Priority
type Priority int

const (
	P0 = Priority(iota)
	P1
	P2
	P3
	P4
	P5
)

func Todo(pr Priority, f string, a ...interface{}) (err error) {
	return todo.Printf(pr, f, a...)
}

func (p todoPrinter) Printf(pr Priority, f string, a ...interface{}) (err error) {
	if !p.on {
		return
	}

	if p.includesStack {
		si, _ := MakeStackInfo(1)
		f = "%s %s" + f
		a = append([]interface{}{pr, si}, a...)
	}

	_, err = fmt.Fprintln(
		p.f,
		fmt.Sprintf(f, a...),
	)

	return
}
