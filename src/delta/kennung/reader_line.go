package kennung

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type ReaderLine struct {
	Set
}

func (rl *ReaderLine) ReadFrom(r1 io.Reader) (n int64, err error) {
	errors.TodoP4("add expanders")
	rl.Set = MakeSet(Expanders{})
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

		if err = rl.Set.Set(line); err != nil {
			err = errors.Wrap(err)
			return
		}
	}
}
