package checkout_store

import (
	"strconv"
	"strings"
	"time"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/hinweis"
)

// type EntryType int

// const (
// 	EntryTypeUnknown = EntryType(iota)
// 	EntryTypeZettel  = node_type.TypeZettel
// 	EntryTypeAkte    = node_type.TypeAkte
// )

type Entry struct {
	Time    time.Time
	Hinweis hinweis.Hinweis
	Sha     sha.Sha
}

func (e Entry) String() string {
	sb := &strings.Builder{}
	sb.WriteString(e.Sha.String())
	sb.WriteString(" ")
	sb.WriteString(e.Hinweis.String())
	sb.WriteString(" ")
	sb.WriteString(strconv.FormatInt(e.Time.Unix(), 10))
	return sb.String()
}

func (e *Entry) Set(s string) (err error) {
	elements := strings.Split(s, " ")

	if len(elements) != 2 {
		err = errors.Errorf("expected 2 elements, but got %d: %q", len(elements), elements)
		return
	}

	if err = e.Sha.Set(elements[0]); err != nil {
		err = errors.Error(err)
		return
	}

	if err = e.Hinweis.Set(elements[1]); err != nil {
		err = errors.Error(err)
		return
	}

	var i int64

	if i, err = strconv.ParseInt(elements[2], 10, 64); err != nil {
		err = errors.Error(err)
		return
	}

	e.Time = time.Unix(i, 0)

	return
}

// func (e Entry) Type() EntryType {
//   switch strings.ToLower(path.Ext(e.Path)) {
//     case ".md":
//     default:
//   }
// }
