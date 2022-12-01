package format

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

func ReadLines(
	r1 io.Reader,
	rffs ...FuncReadLine,
) (n int64, err error) {
	r := bufio.NewReader(r1)
	i := 0

	for {
		var rawLine, line string

		rawLine, err = r.ReadString('\n')
		n += int64(len(rawLine))

		if err != nil && !errors.IsEOF(err) {
			err = errors.Wrap(err)
			return
		}

		if errors.IsEOF(err) {
			err = nil
			break
		}

		line = strings.TrimSuffix(rawLine, "\n")

		if len(rffs) == i {
			//TODO add line
			err = errors.Errorf("ran out of read line funcs before fully consuming reader")
			return
		}

		frl := rffs[i]

		if err = frl(line); err != nil {
			if errors.IsEOF(err) {
				err = nil
				i++
				continue
			} else {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func MakeLineReaderKeyValue(
	key string,
	value FuncReadLine,
) FuncReadLine {
	return func(line string) (err error) {
		loc := strings.Index(line, " ")

		if loc == -1 {
			err = errors.Errorf("expected at least one space, but found none: %q", line)
			return
		}

		if line[:loc] != key {
			err = errors.Errorf("expected key %q but got %q", key, line[:loc])
			return
		}

		if err = value(line[loc+1:]); err != nil {
			err = errors.Errorf("%s: %q", err, line[loc+1:])
			return
		}

		err = io.EOF

		return
	}
}

func MakeLineReaderIgnoreErrors(
	in FuncReadLine,
) FuncReadLine {
	return func(line string) (err error) {
		in(line)

		return
	}
}
