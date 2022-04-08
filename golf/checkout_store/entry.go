package checkout_store

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

// type EntryType int

// const (
// 	EntryTypeUnknown = EntryType(iota)
// 	EntryTypeZettel  = node_type.TypeZettel
// 	EntryTypeAkte    = node_type.TypeAkte
// )

type Entry struct {
	ZettelTime time.Time
	AkteTime   time.Time
	External   stored_zettel.External
}

func (e Entry) String() string {
	sb := &strings.Builder{}
	var b []byte
	var err error

	if b, err = json.Marshal(e); err != nil {
		panic(errors.Error(err))
	}

	sb.Write(b)

	return sb.String()
}

func (e *Entry) Set(s string) (err error) {
	if err = json.Unmarshal([]byte(s), &e); err != nil {
		err = errors.Error(err)
	}

	return
}

// func (e Entry) Type() EntryType {
//   switch strings.ToLower(path.Ext(e.Path)) {
//     case ".md":
//     default:
//   }
// }
