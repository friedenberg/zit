package zettel

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type FormatValue struct {
	Format
}

func (f FormatValue) String() string {
	switch f1 := f.Format.(type) {
	case *Text:
		return "text"

	case *Objekte:
		return "objekte"

	case *Akte:
		return "akte"

	default:
		return fmt.Sprintf("%T", f1)
	}
}

func (f *FormatValue) Set(v string) (err error) {
	switch strings.TrimSpace(strings.ToLower(v)) {
	case "akte":
		f.Format = &Akte{}

	case "objekte":
		f.Format = &Objekte{}

	case "text":
		f.Format = &Text{}

	default:
		err = errors.Errorf("unsupported format type: %s", v)
		return
	}

	return
}
