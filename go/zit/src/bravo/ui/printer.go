package ui

import (
	"fmt"
	"os"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/primordial"
)

func MakePrinter(f *os.File) printer {
	return MakePrinterOn(f, true)
}

func MakePrinterOn(f *os.File, on bool) printer {
	return printer{
		f:     f,
		isTty: primordial.IsTty(f),
		on:    on,
	}
}

type printer struct {
	f     *os.File
	isTty bool
	on    bool
}

func (p printer) withOn(on bool) printer {
	p.on = on
	return p
}

func (p printer) GetPrinter() Printer {
	return p
}

func (p printer) Write(b []byte) (n int, err error) {
	if !p.on {
		n = len(b)
		return
	}

	return p.f.Write(b)
}

func (p printer) GetFile() *os.File {
	return p.f
}

func (p printer) IsTty() bool {
	return p.isTty
}

func (p printer) Print(a ...interface{}) (err error) {
	if !p.on {
		return
	}

	_, err = fmt.Fprintln(
		p.f,
		a...,
	)

	return
}

func (p printer) printfStack(depth int, f string, a ...interface{}) (err error) {
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

func (p printer) Printf(f string, a ...interface{}) (err error) {
	if !p.on {
		return
	}

	_, err = fmt.Fprintln(
		p.f,
		fmt.Sprintf(f, a...),
	)

	return
}
