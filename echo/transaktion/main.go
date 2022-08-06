package transaktion

import (
	"time"

	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/hinweis"
)

type Transaktion struct {
	time.Time
	LineItems []TransaktionLineItem
}

type TransaktionLineItem struct {
	hinweis.Hinweis
	sha.Sha
}
