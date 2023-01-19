package transaktion

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/foxtrot/sku"
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

	r.Transaktion.Skus = sku.MakeMutableSet()

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
			var o sku.SkuLike

			if o, err = sku.MakeSkuTransacted(r.Transaktion.Time, line); err != nil {
				err = errors.Wrapf(err, "failed to read line: %s", line)
				return
			}

			r.Transaktion.Skus.Add(o)
		}

		r.readState.lineNo += 1
	}

	return
}
