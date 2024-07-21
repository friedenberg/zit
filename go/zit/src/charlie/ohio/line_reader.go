package ohio

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
)

func MakeLineReaderIterateStrict(
	rffs ...interfaces.FuncSetString,
) interfaces.FuncSetString {
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
	rffs ...interfaces.FuncSetString,
) interfaces.FuncSetString {
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
	dict map[string]interfaces.FuncSetString,
) interfaces.FuncSetString {
	si, _ := errors.MakeStackInfo(1)

	return func(line string) (err error) {
		loc := strings.Index(line, " ")

		if loc == -1 {
			err = si.Errorf(
				"expected at least one space, but found none: %q",
				line,
			)
			return
		}

		key := line[:loc]
		value := line[loc+1:]

		var reader interfaces.FuncSetString
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

func MakeLineReaderRepeat(
	in interfaces.FuncSetString,
) interfaces.FuncSetString {
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
	in interfaces.FuncSetString,
) interfaces.FuncSetString {
	return func(line string) (err error) {
		in(line)

		return
	}
}
