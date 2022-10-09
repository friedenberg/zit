package transaktion

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Reader struct {
	Transaktion
	readState
}

type readState struct {
	lineNo int
}

func (r *Reader) ReadFrom(r1 io.Reader) (n int64, err error) {
	br := bufio.NewReader(r1)

	r.Transaktion.Objekten = make(map[string]Objekte)

	for {
		var line string

		line, err = br.ReadString('\n')
		n += int64(len(line))

		if err != nil && err != io.EOF {
			err = errors.Wrap(err)
			return
		}

		if err == io.EOF {
			err = nil
			break
		}

		line = strings.TrimSuffix(line, "\n")

		switch r.readState.lineNo {
		case 0:
			if err = r.Transaktion.Time.Set(line); err != nil {
				err = errors.Wrapf(err, "failed to read time: %s", line)
				return
			}

		default:
			var o Objekte

			if err = o.Set(line); err != nil {
				err = errors.Wrapf(err, "failed to read line: %s", line)
				return
			}

			k := o.GetKey()
			o1, ok := r.Transaktion.Objekten[k]

			if ok {
				errors.Err().Printf(
					"Transation %s has duplicate entries:  (%s %s %s) & (%s %s %s)",
					r.Transaktion,
					o1.Gattung,
					o1.Id,
					o1.Sha,
					o.Gattung,
					o.Id,
					o.Sha,
				)
			}

			r.Transaktion.Objekten[k] = o
		}

		r.readState.lineNo += 1
	}

	return
}
