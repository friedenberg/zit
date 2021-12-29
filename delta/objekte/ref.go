package objekte

import (
	"fmt"
	"strings"
)

type Ref struct {
	Type  _Type
	Value string
}

func (o Ref) String() string {
	return fmt.Sprintf("%s %s", o.Type, o.Value)
}

func (o *Ref) Set(v string) (err error) {
	idxSpace := strings.Index(v, " ")

	if idxSpace == -1 {
		err = _Errorf("expected at least one space character, but got none")

		return
	}

	if err = o.Type.Set(strings.TrimSpace(v[:idxSpace])); err != nil {
		err = _Error(err)
		return
	}

	o.Value = strings.TrimSpace(v[idxSpace:])

	return
}

func (o Ref) Sha() (s _Sha, err error) {
	if err = s.Set(o.Value); err != nil {
		err = _Error(err)
		return
	}

	return
}
