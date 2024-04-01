package hinweisen

import (
	"bufio"
	"io"
	"os"

	"code.linenisgreat.com/zit/src/alfa/coordinates"
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/charlie/files"
)

type provider []string

func newProvider(path string) (p provider, err error) {
	var f *os.File

	if f, err = files.Open(path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.Deferred(&err, f.Close)

	r := bufio.NewReader(f)

	for {
		var line string
		line, err = r.ReadString('\n')

		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		p = append(p, Clean(line))
	}

	return
}

func (p provider) Hinweis(i coordinates.Int) (s string, err error) {
	if len(p)-1 < int(i) {
		err = errors.Errorf(
			"insuffient ids. requested %d, have %d, last %s",
			int(i),
			len(p)-1,
			p[len(p)-1],
		)

		return
	}

	s = p[i]

	return
}

func (p provider) Len() int {
	return len(p)
}

func (p provider) Kennung(v string) (i int, err error) {
	v = Clean(v)

	var s string

	for i, s = range p {
		if s == v {
			return
		}
	}

	err = ErrDoesNotExist{
		Value: v,
	}

	return
}
