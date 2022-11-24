package zettel

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/konfig"
)

type FormatValue struct {
	out    io.Writer
	konfig konfig.Konfig
	Format
}

func MakeFormatValue(out io.Writer, k konfig.Konfig) *FormatValue {
	return &FormatValue{
		out:    out,
		konfig: k,
	}
}

func (f FormatValue) String() string {
	switch f1 := f.Format.(type) {
	// case *collections_coding.EncoderJson[Zettel]:
	// 	return "json"

	case *Text:
		return "text"

	case *Objekte:
		return "objekte"

	case *Akte:
		return "akte"

	case *EncoderTypActionNames:
		return "typ-action-names"

	default:
		return fmt.Sprintf("%T", f1)
	}
}

func (f *FormatValue) Set(v string) (err error) {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case "akte":
		f.Format = &Akte{}

	// case "json":
	// 	f.EncoderLike = collections_coding.MakeEncoderJson[Zettel](f.out)

	case "objekte":
		f.Format = &Objekte{}

	case "text":
		f.Format = &Text{}

	case "typ-action-names":
		f.Format = &EncoderTypActionNames{
			konfig: f.konfig,
		}

	default:
		err = errors.Errorf("unsupported format type: %s", v)
		return
	}

	return
}
