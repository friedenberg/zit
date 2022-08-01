package hinweisen

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/friedenberg/zit/alfa/errors"
)

type provider []string

func newProvider(path string) (p provider, err error) {
	var f *os.File

	if f, err = _Open(path); err != nil {
		err = errors.Error(err)
		return
	}

	defer _Close(f)

	r := bufio.NewReader(f)

	for {
		var line string
		line, err = r.ReadString('\n')

		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			err = errors.Error(err)
			return
		}

		p = append(p, strings.TrimSuffix(line, "\n"))
	}

	return
}

func (p provider) Hinweis(i _Int) (s string, err error) {
	if len(p)-1 < int(i) {
		err = errors.Errorf("insuffient ids")
		return
	}

	s = p[i]

	return
}
