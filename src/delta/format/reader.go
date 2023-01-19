package format

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/charlie/collections"
)

type readerFrom[T any] struct {
	rf FuncReaderFormat[T]
	e  *T
}

func (rf readerFrom[T]) ReadFrom(r io.Reader) (n int64, err error) {
	return rf.rf(r, rf.e)
}

func MakeReaderFrom[T any](
	rf FuncReaderFormat[T],
	e *T,
) io.ReaderFrom {
	return readerFrom[T]{
		rf: rf,
		e:  e,
	}
}

func ReadSep(
	delim byte,
	r1 io.Reader,
	rffs ...schnittstellen.FuncSetString,
) (n int64, err error) {
	r := bufio.NewReader(r1)
	i := 0

	var last error

	for {
		var rawLine, line string

		rawLine, err = r.ReadString(delim)
		n += int64(len(rawLine))

		if err != nil && !errors.IsEOF(err) {
			err = errors.Wrap(err)
			return
		}

		isEof := errors.IsEOF(err)
		err = nil

		line = strings.TrimSuffix(rawLine, string([]byte{delim}))

		if len(rffs) == i {
			//TODO add line
			err = errors.Errorf("ran out of read line funcs before fully consuming reader")

			if last != nil {
				err = errors.MakeMulti(err, last)
			}

			return
		}

		frl := rffs[i]

		if err = frl(line); err != nil {
			last = err
			err = nil
		}

		i++

		if isEof {
			break
		}
	}

	if last != nil && !errors.Is(last, collections.ErrStopIteration) {
		err = last
	}

	return
}

func ReadLines(
	r1 io.Reader,
	rffs ...schnittstellen.FuncSetString,
) (n int64, err error) {
	r := bufio.NewReader(r1)
	i := 0

	var last error

	isEOF := false

	for {
		if isEOF {
			break
		}

		var rawLine, line string

		rawLine, err = r.ReadString('\n')
		n += int64(len(rawLine))

		if err != nil && !errors.IsEOF(err) {
			err = errors.Wrap(err)
			return
		}

		if errors.IsEOF(err) {
			isEOF = true
			err = nil
		}

		line = strings.TrimSuffix(rawLine, "\n")

		if line == "" {
			continue
		}

		if len(rffs) == i {
			//TODO add line
			err = errors.Errorf("ran out of read line funcs before fully consuming reader")

			if last != nil {
				err = errors.MakeMulti(err, last)
			}

			return
		}

		frl := rffs[i]

		if err = frl(line); err != nil {
			last = err
			err = nil
			i++
		}
	}

	if last != nil && !errors.Is(last, collections.ErrStopIteration) {
		err = last
	}

	return
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

		err = collections.ErrStopIteration

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
			err = errors.Errorf("expected key %q but got %q", key, keyActual)
			return
		}

		if err = valueReader(value); err != nil {
			err = errors.Errorf("%s: %q", err, value)
			return
		}

		err = collections.ErrStopIteration

		return
	}
}

func MakeLineReaderRepeat(
	in schnittstellen.FuncSetString,
) schnittstellen.FuncSetString {
	return func(line string) (err error) {
		if err = in(line); err != nil {
			if errors.Is(err, collections.ErrStopIteration) {
				err = nil
			} else {
				err = errors.Wrap(err)
				return
			}
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
