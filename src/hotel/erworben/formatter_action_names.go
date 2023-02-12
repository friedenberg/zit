package erworben

import (
	"fmt"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/script_config"
)

type formatterActionNames struct{}

func MakeFormatterActionNames() *formatterActionNames {
	return &formatterActionNames{}
}

func (f formatterActionNames) Format(
	w io.Writer,
	ct map[string]script_config.ScriptConfig,
) (n int64, err error) {
	for v, v1 := range ct {
		var n1 int

		if n1, err = io.WriteString(
			w,
			fmt.Sprintf("%s\t%s\n", v, v1.Description),
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += int64(n1)
	}

	return
}
