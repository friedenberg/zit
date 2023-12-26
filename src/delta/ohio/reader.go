package ohio

import (
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
)

func ReadAllOrDieTrying(r io.Reader, b []byte) (n int, err error) {
	var acc int

	for n < len(b) {
    acc, err = r.Read(b[n:])
		n += acc
		if err != nil {
			return
		}
	}

	return
}

func MakeLineReaderIterateStrict(
	rffs ...schnittstellen.FuncSetString,
) schnittstellen.FuncSetString {
	si, _ := errors.MakeStackInfo(1)
	var i int64

	return func(v string) (err error) {
		if int64(len(rffs))-1 < i {
			err = errors.Wrap(&ErrExhaustedFuncSetStringersLine{
				error:  err,
				string: v,
			})

			return
		}

		if err = rffs[i](v); err != nil {
			err = si.Wrapf(err, "Value: %s", v)
			return
		}

		i++

		return
	}
}

func MakeLineReaderIterate(
	rffs ...schnittstellen.FuncSetString,
) schnittstellen.FuncSetString {
	si, _ := errors.MakeStackInfo(1)
	var i int64

	return func(v string) (err error) {
		for {
			if int64(len(rffs))-1 < i {
				err = errors.Wrap(&ErrExhaustedFuncSetStringersLine{
					error:  err,
					string: v,
				})

				return
			}

			if err = rffs[i](v); err != nil {
				i++
				err = si.Wrapf(err, "Value: %s", v)
				continue
			}

			return
		}
	}
}

func MakeLineReaderKeyValues(
	dict map[string]schnittstellen.FuncSetString,
) schnittstellen.FuncSetString {
	si, _ := errors.MakeStackInfo(1)

	return func(line string) (err error) {
		loc := strings.Index(line, " ")

		if loc == -1 {
			err = si.Errorf("expected at least one space, but found none: %q", line)
			return
		}

		key := line[:loc]
		value := line[loc+1:]

		var reader schnittstellen.FuncSetString
		ok := false

		if reader, ok = dict[key]; !ok {
			err = si.Errorf("key not supported: %q", key)
			return
		}

		if err = reader(value); err != nil {
			err = si.Errorf("%s: %q", err, value)
			return
		}

		return
	}
}

func MakeLineReaderKeyValue(
	key string,
	valueReader schnittstellen.FuncSetString,
) schnittstellen.FuncSetString {
	return func(line string) (err error) {
		loc := strings.Index(line, " ")

		if loc == -1 {
			err = errors.Errorf("expected at least one space, but found none: %q", line)
			return
		}

		keyActual := line[:loc]
		value := line[loc+1:]

		if keyActual != key {
			err = errors.Wrap(&ErrExhaustedFuncSetStringersLine{
				string: value,
			})

			return
		}

		if err = valueReader(value); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}
}

func MakeLineReaderRepeat(
	in schnittstellen.FuncSetString,
) schnittstellen.FuncSetString {
	return func(line string) (err error) {
		if err = in(line); err != nil {
			err = errors.Wrap(&ErrExhaustedFuncSetStringersLine{
				error:  err,
				string: line,
			})

			return
		}

		return
	}
}

func MakeLineReaderIgnoreErrors(
	in schnittstellen.FuncSetString,
) schnittstellen.FuncSetString {
	return func(line string) (err error) {
		in(line)

		return
	}
}

func MakeLineReaderNop() schnittstellen.FuncSetString {
	return func(line string) (err error) {
		return
	}
}
