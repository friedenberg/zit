package typ

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/collections"
	"github.com/friedenberg/zit/src/delta/konfig"
)

type EncoderValue struct {
	collections.EncoderLike[Typ]
	out    io.Writer
	konfig konfig.Konfig
}

func MakeEncoderValue(
	k konfig.Konfig,
	out io.Writer,
) EncoderValue {
	return EncoderValue{
		out:    out,
		konfig: k,
	}
}

func (f EncoderValue) String() string {
	switch f1 := f.EncoderLike.(type) {
	// case *Text:
	// 	return "text"

	// case *Objekte:
	// 	return "objekte"

	case *EncoderActionNames:
		return "action-names"

	default:
		return fmt.Sprintf("%T", f1)
	}
}

func (f *EncoderValue) Set(v string) (err error) {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case "action-names":
		f.EncoderLike = MakeEncoderActionNames(f.out, f.konfig)

	default:
		err = errors.Errorf("unsupported format type: %s", v)
		return
	}

	return
}
