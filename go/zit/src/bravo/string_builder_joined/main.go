package string_builder_joined

import (
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
)

type builder struct {
	didWrite   bool
	joinString string
	builder    strings.Builder
}

func Make(joinString string) (b *builder) {
	b = &builder{
		joinString: joinString,
	}

	return
}

func (b *builder) setDidWrite() {
	b.didWrite = true
}

func (b *builder) WriteString(v string) (n int, err error) {
	defer b.setDidWrite()
	var n1 int

	if b.didWrite {
		if n1, err = b.builder.WriteString(b.joinString); err != nil {
			err = errors.Wrap(err)
			return
		}

		n += n1
	}

	if n1, err = b.builder.WriteString(v); err != nil {
		err = errors.Wrap(err)
		return
	}

	n += n1

	return
}

func (b *builder) String() string {
	return b.builder.String()
}
