package hinweisen

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/kennung"
	"github.com/friedenberg/zit/src/charlie/open_file_guard"
)

type provider []string

func newProvider(path string) (p provider, err error) {
	var f *os.File

	if f, err = open_file_guard.Open(path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer open_file_guard.Close(f)

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

		p = append(p, strings.TrimSuffix(line, "\n"))
	}

	return
}

func (p provider) Hinweis(i kennung.Int) (s string, err error) {
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
	v = strings.ToLower(v)
	v = strings.Map(
		func(r rune) rune {
			if r > 'z' {
				return -1
			}

			return r
		},
		v,
	)

	var s string

	for i, s = range p {
		//TODO move to init and make common
		s = strings.ToLower(s)
		s = strings.Map(
			func(r rune) rune {
				if r > 'z' {
					return -1
				}

				return r
			},
			s,
		)

		if s == v {
			return
		}
	}

	err = ErrDoesNotExist{
		Value: v,
	}

	return
}
