package errors

import (
	"fmt"
	"io"
	"os"
)

var (
	out, err, log printer
)

func init() {
	out = printer{
		f:  os.Stdout,
		on: true,
	}

	err = printer{
		f:  os.Stderr,
		on: true,
	}

	log = printer{
		includesStack: true,
		f:             os.Stderr,
	}
}

type printer struct {
	f             *os.File
	includesStack bool
	on            bool
}

func Out() printer {
	return out
}

func Err() printer {
	return err
}

func Log() printer {
	return log
}

func (p printer) PrintDebug(vs ...interface{}) (err error) {
	if !p.on {
		return
	}

	si, _ := MakeStackInfo(1)

	for _, v := range vs {
		if _, err = fmt.Fprintf(
			p.f,
			"%s%#v\n",
			si,
			v,
		); err != nil {
			err = Wrap(err)
			return
		}
	}

	return
}

func (p printer) Print(a ...interface{}) (err error) {
	if !p.on {
		return
	}

	//TODO add support for includesStack
	_, err = fmt.Fprintln(
		p.f,
		a...,
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

func (p printer) Caller(i int, f string, vs ...interface{}) {
	if !p.on {
		return
	}

	st, _ := MakeStackInfo(i + 1)

	vs = append([]interface{}{st}, vs...)
	//TODO strip trailing newline and add back
	io.WriteString(p.f, fmt.Sprintf("%s"+f+"\n", vs...))
}

func (p printer) CallerNonEmpty(i int, v interface{}) {
	if v != nil {
		p.Caller(i+1, "%s", v)
	}
}

func (p printer) FunctionName(skip int) {
	if !p.on {
		return
	}

	st, _ := MakeStackInfo(skip + 1)
	io.WriteString(p.f, fmt.Sprintf("%s%s\n", st, st.function))
}
