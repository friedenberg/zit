package ui

import (
	"fmt"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

type devPrinter struct {
	prodPrinter
	includesStack bool
}

func (p devPrinter) Print(a ...interface{}) (err error) {
	if !p.on {
		return
	}

	if p.includesStack {
		si, _ := errors.MakeStackInfo(1)
		a = append([]interface{}{si.StringNoFunctionName()}, a...)
	}

	return p.prodPrinter.Print(a...)
}

func (p devPrinter) Printf(f string, a ...interface{}) (err error) {
	if !p.on {
		return
	}

	if p.includesStack {
		si, _ := errors.MakeStackInfo(1)
		f = "%s " + f
		a = append([]interface{}{si.StringNoFunctionName()}, a...)
	}

	return p.prodPrinter.Printf(f, a...)
}

func (p devPrinter) Caller(i int, vs ...interface{}) {
	if !p.on {
		return
	}

	st, _ := errors.MakeStackInfo(i + 1)

	vs = append([]interface{}{st}, vs...)
	// TODO-P4 strip trailing newline and add back
	p.prodPrinter.Print(vs...)
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

	st, _ := errors.MakeStackInfo(skip + 1)
	io.WriteString(p.f, fmt.Sprintf("%s%s\n", st, st.Function))
}
