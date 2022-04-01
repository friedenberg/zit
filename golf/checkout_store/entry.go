package checkout_store

import (
	"time"

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

// func (e Entry) Type() EntryType {
//   switch strings.ToLower(path.Ext(e.Path)) {
//     case ".md":
//     default:
//   }
// }
