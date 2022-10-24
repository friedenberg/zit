package id_set

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/id"
)

type ReaderLine struct {
	ProtoIdSet
	Set
}

func (rl *ReaderLine) ReadFrom(r1 io.Reader) (n int64, err error) {
	rl.Set = Make(0)
	r := bufio.NewReader(r1)

	for {
		var line string

		line, err = r.ReadString('\n')
		n += int64(len(line))

		switch {
		case err == nil:
			break

		case errors.IsEOF(err):
			err = nil
			return

		default:
			err = errors.Wrap(err)
			return
		}

		if line == "" {
			continue
		}

		var i id.Id

		if i, err = rl.ProtoIdSet.MakeOne(line); err != nil {
			err = errors.Wrap(err)
			return
		}

		rl.Set.Add(i)
	}
}
