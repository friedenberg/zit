package format

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/ohio"
)

type lineReader struct {
	delim         byte
	passthruEmpty bool
	reader        schnittstellen.FuncSetString
}

func MakeLineReaderConsumeEmpty(
	rffs schnittstellen.FuncSetString,
) io.ReaderFrom {
	return lineReader{
		delim:  '\n',
		reader: rffs,
	}
}

func MakeLineReaderPassThruEmpty(
	rffs schnittstellen.FuncSetString,
) io.ReaderFrom {
	return lineReader{
		delim:         '\n',
		passthruEmpty: true,
		reader:        rffs,
	}
}

func MakeDelimReaderConsumeEmpty(
	delim byte,
	rffs schnittstellen.FuncSetString,
) io.ReaderFrom {
	return lineReader{
		delim:  delim,
		reader: rffs,
	}
}

func ReadLines(
	r1 io.Reader,
	rffs schnittstellen.FuncSetString,
) (n int64, err error) {
	lr := MakeLineReaderConsumeEmpty(rffs)
	return lr.ReadFrom(r1)
}

func ReadSep(
	delim byte,
	r1 io.Reader,
	rffs schnittstellen.FuncSetString,
) (n int64, err error) {
	lr := MakeDelimReaderConsumeEmpty(delim, rffs)
	return lr.ReadFrom(r1)
}

func (lr lineReader) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	var (
		lineNo int64
		isEOF  bool
	)

	for {
		if isEOF {
			break
		}

		var rawLine, line string

		rawLine, err = r.ReadString(lr.delim)
		n1 := len(rawLine)
		n += int64(n1)

		if err != nil && !errors.IsEOF(err) {
			err = errors.Wrap(err)
			return
		}

		if errors.IsEOF(err) {
			isEOF = true
			err = nil

			if n1 == 0 {
				break
			}
		}

		line = strings.TrimSuffix(rawLine, string([]byte{lr.delim}))

		if line == "" && !lr.passthruEmpty {
			continue
		}

		if err = lr.reader(line); err != nil {
			err = errors.Wrap(ohio.ErrExhaustedFuncSetStringersWithDelim(err, lr.delim))
			return
		}

		lineNo++
	}

	return
}
