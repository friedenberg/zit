package zettels

import (
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

type Chain struct {
	Hinweis hinweis.Hinweis
	//stored in reverse (latest is at 0)
	Zettels []stored_zettel.Stored
}
