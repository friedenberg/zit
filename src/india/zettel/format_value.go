package zettel

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/hotel/konfig_compiled"
)

type FormatValue struct {
	out    io.Writer
	konfig konfig_compiled.Compiled
	Format
}

func MakeFormatValue(out io.Writer, k konfig_compiled.Compiled) *FormatValue {
	return &FormatValue{
		out:    out,
		konfig: k,
	}
}

func (f FormatValue) String() string {
	switch f1 := f.Format.(type) {
	// case *collections_coding.EncoderJson[Zettel]:
	// 	return "json"

	case *FormatToFormat2:
		switch f2 := f1.Format.(type) {
		case *Text2:
			return "text"

		default:
			return fmt.Sprintf("%T", f2)
		}

	case *FormatObjekte:
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
		f.Format = &FormatObjekte{}

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
