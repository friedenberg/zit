package ui

import (
	"fmt"
	"io"
	"log"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
)

var (
	SetOutput = log.SetOutput
	verbose   bool
	isTest    bool
)

func SetVerbose() {
	printerLog.on = true
	printerDebug.on = true
	verbose = true
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
	log.Print("verbose")
}

func SetTesting() {
	isTest = true
	errors.SetTesting()
	SetVerbose()
}

func IsVerbose() bool {
	return verbose
}

type ProdPrinter interface {
	Print(v ...interface{}) error
	Printf(format string, v ...interface{}) error
}

type DevPrinter interface {
	ProdPrinter
	Caller(i int, vs ...interface{})
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

func DebugAllowCommit() DevPrinter {
	return printerDebug
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

func (p devPrinter) Print(a ...interface{}) (err error) {
	if !p.on {
		return
	}

	if p.includesStack {
		si, _ := errors.MakeStackInfo(1)
		a = append([]interface{}{si}, a...)
	}

	return p.prodPrinter.Print(a...)
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

func (p devPrinter) Printf(f string, a ...interface{}) (err error) {
	if !p.on {
		return
	}

	if p.includesStack {
		si, _ := errors.MakeStackInfo(1)
		f = "%s" + f
		a = append([]interface{}{si}, a...)
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
