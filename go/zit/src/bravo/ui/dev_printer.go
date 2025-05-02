package ui

import (
	"fmt"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

type devPrinter struct {
	printer
	includesStack bool
}

func (p devPrinter) Print(a ...interface{}) (err error) {
	if !p.on {
		return
	}

	if p.includesStack {
		si, _ := errors.MakeStackFrame(1)
		a = append([]interface{}{si.StringNoFunctionName()}, a...)
	}

	return p.printer.Print(a...)
}

func (p devPrinter) Printf(f string, a ...interface{}) (err error) {
	if !p.on {
		return
	}

	if p.includesStack {
		si, _ := errors.MakeStackFrame(1)
		f = "%s " + f
		a = append([]interface{}{si.StringNoFunctionName()}, a...)
	}

	return p.printer.Printf(f, a...)
}

func (p devPrinter) Caller(i int, vs ...interface{}) {
	if !p.on {
		return
	}

	st, _ := errors.MakeStackFrame(i + 1)

	vs = append([]interface{}{st}, vs...)
	// TODO-P4 strip trailing newline and add back
	p.printer.Print(vs...)
}

func (p devPrinter) CallerNonEmpty(i int, v interface{}) {
	if v != nil {
		p.Caller(i+1, "%s", v)
	}
}

func (p devPrinter) FunctionName(skip int) {
	if !p.on {
		return
	}

	st, _ := errors.MakeStackFrame(skip + 1)
	io.WriteString(p.f, fmt.Sprintf("%s%s\n", st, st.Function))
}

//go:noinline
func (p devPrinter) Stack(skip, count int) {
	if !p.on {
		return
	}

	frames := errors.MakeStackFrames(skip+1, count)

	io.WriteString(
		p.f,
		fmt.Sprintf(
			"Printing Stack (skip: %d, count requested: %d, count actual: %d):\n\n",
			skip,
			count,
			len(frames),
		),
	)

	for i, frame := range frames {
		io.WriteString(
			p.f,
			fmt.Sprintf("%s (%d)\n", frame.StringLogLine(), i),
		)
	}
}
