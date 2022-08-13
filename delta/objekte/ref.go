package objekte

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/node_type"
	"github.com/friedenberg/zit/bravo/sha"
)

type Ref struct {
	Type  node_type.Type
	Value string
}

func (o Ref) String() string {
	return fmt.Sprintf("%s %s", o.Type, o.Value)
}

func (o *Ref) Set(v string) (err error) {
	idxSpace := strings.Index(v, " ")

	if idxSpace == -1 {
		err = errors.Errorf("expected at least one space character, but got none")

		return
	}

	if err = o.Type.Set(strings.TrimSpace(v[:idxSpace])); err != nil {
		err = errors.Error(err)
		return
	}

	o.Value = strings.TrimSpace(v[idxSpace:])

	return
}

func (o Ref) Sha() (s sha.Sha, err error) {
	if err = s.Set(o.Value); err != nil {
		err = errors.Error(err)
		return
	}

	return
}
