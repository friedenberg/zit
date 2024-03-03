package organize_text2

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
)

type line struct {
	isComment bool
	prefix    string
	value     string
}

func (l line) String() string {
	return fmt.Sprintf("%s %s", l.prefix, l.value)
}

func (l *line) Set(v string) (err error) {
	v = strings.TrimSpace(v)

	if len(v) == 0 {
		err = errors.Errorf("line not long enough")
		return
	}

	firstSpace := strings.Index(v, " ")

	if firstSpace == -1 {
		l.prefix = v
		return
	}

	l.prefix = strings.TrimSpace(v[:firstSpace])
	l.value = strings.TrimSpace(v[firstSpace:])

	return
}

func (l line) PrefixRune() rune {
	if len(l.prefix) == 0 {
		panic(errors.Errorf("cannot find prefix in line: %q", l.value))
	}

	return rune(l.prefix[0])
}

func (l line) Depth(r rune) (depth int, err error) {
	for i, c := range l.prefix {
		if c != r {
			err = errors.Errorf("rune at index %d is %c and not %c", i, c, r)
			return
		}

		depth++
	}

	return
}
