package errors

import (
	"fmt"
	"io"
	"os"
)

type ProdPrinter interface {
	Print(v ...interface{}) error
	Printf(format string, v ...interface{}) error
}

type DevPrinter interface {
	ProdPrinter
	Caller(i int, f string, vs ...interface{})
	FunctionName(skip int)
}

var (
	printerOut, printerErr   prodPrinter
	printerLog, printerDebug devPrinter
)

func init() {
	printerOut = prodPrinter{
		f:  os.Stdout,
		on: true,
	}

	printerErr = prodPrinter{
		f:  os.Stderr,
		on: true,
	}

	printerLog = devPrinter{
		prodPrinter: prodPrinter{
			f: os.Stderr,
		},
		includesStack: true,
	}

	printerDebug = devPrinter{
		prodPrinter: prodPrinter{
			f: os.Stderr,
			// TODO-P2 determine thru compilation
			on: true,
		},
		includesStack: true,
	}
}

type prodPrinter struct {
	f  io.Writer
	on bool
}

type devPrinter struct {
	prodPrinter
	includesStack bool
}

func Out() ProdPrinter {
	return printerOut
}

func Err() ProdPrinter {
	return printerErr
}

func Log() DevPrinter {
	return printerLog
}

func Debug() DevPrinter {
	return printerDebug
}

// func MakePrinter(o io.Writer) *printer {
// 	return &printer{
// 		f:  o,
// 		on: true,
// 	}
// }

// func (p printer) PrintDebug(vs ...interface{}) (err error) {
// 	if !p.on {
// 		return
// 	}

// 	si, _ := MakeStackInfo(1)

// 	for _, v := range vs {
// 		if _, err = fmt.Fprintf(
// 			p.f,
// 			"%s%#v\n",
// 			si,
// 			v,
// 		); err != nil {
// 			err = Wrap(err)
// 			return
// 		}
// 	}

// 	return
// }

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

func (p devPrinter) Print(a ...interface{}) (err error) {
	if !p.on {
		return
	}

	if p.includesStack {
		si, _ := MakeStackInfo(1)
		a = append([]interface{}{si}, a...)
	}

	return p.prodPrinter.Print(a...)
}

func (p prodPrinter) printfStack(depth int, f string, a ...interface{}) (err error) {
	if !p.on {
		return
	}

	si, _ := MakeStackInfo(1 + depth)
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

func (p devPrinter) Printf(f string, a ...interface{}) (err error) {
	if !p.on {
		return
	}

	if p.includesStack {
		si, _ := MakeStackInfo(1)
		f = "%s" + f
		a = append([]interface{}{si}, a...)
	}

	return p.prodPrinter.Printf(f, a...)
}

func (p devPrinter) Caller(i int, f string, vs ...interface{}) {
	if !p.on {
		return
	}

	st, _ := MakeStackInfo(i + 1)

	vs = append([]interface{}{st}, vs...)
	// TODO-P4 strip trailing newline and add back
	io.WriteString(p.f, fmt.Sprintf("%s"+f+"\n", vs...))
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

	st, _ := MakeStackInfo(skip + 1)
	io.WriteString(p.f, fmt.Sprintf("%s%s\n", st, st.function))
}
