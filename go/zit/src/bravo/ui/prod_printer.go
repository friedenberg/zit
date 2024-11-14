package ui

import (
	"fmt"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

type prodPrinter struct {
	f  io.Writer
	on bool
}

func (p prodPrinter) Print(a ...interface{}) (err error) {
	if !p.on {
		return
	}

	_, err = fmt.Fprintln(
		p.f,
		a...,
	)

	return
}

func (p prodPrinter) printfStack(depth int, f string, a ...interface{}) (err error) {
	if !p.on {
		return
	}

	si, _ := errors.MakeStackInfo(1 + depth)
	f = "%s" + f
	a = append([]interface{}{si}, a...)

	_, err = fmt.Fprintln(
		p.f,
		fmt.Sprintf(f, a...),
	)

	return
}

func (p prodPrinter) Printf(f string, a ...interface{}) (err error) {
	if !p.on {
		return
	}

	_, err = fmt.Fprintln(
		p.f,
		fmt.Sprintf(f, a...),
	)

	return
}
