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
	return todo.printf(1, pr, f, a...)
}

func TodoP0(f string, a ...interface{}) (err error) {
	return todo.printf(1, P0, f, a...)
}

func TodoP1(f string, a ...interface{}) (err error) {
	return todo.printf(1, P1, f, a...)
}

func TodoP2(f string, a ...interface{}) (err error) {
	return todo.printf(1, P2, f, a...)
}

func TodoP3(f string, a ...interface{}) (err error) {
	return todo.printf(1, P3, f, a...)
}

func TodoP4(f string, a ...interface{}) (err error) {
	return todo.printf(1, P4, f, a...)
}

func TodoP5(f string, a ...interface{}) (err error) {
	return todo.printf(1, P5, f, a...)
}

func (p todoPrinter) Printf(pr Priority, f string, a ...interface{}) (err error) {
	return p.printf(1, pr, f, a...)
}

func (p todoPrinter) printf(depth int, pr Priority, f string, a ...interface{}) (err error) {
	if !p.on {
		return
	}

	if p.includesStack {
		si, _ := MakeStackInfo(1 + depth)
		f = "%s %s" + f
		a = append([]interface{}{pr, si}, a...)
	}

	_, err = fmt.Fprintln(
		p.f,
		fmt.Sprintf(f, a...),
	)

	return
}
