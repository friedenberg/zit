package store_checkout

import "github.com/friedenberg/zit/src/delta/konfig"

type Konfig struct {
	konfig.Konfig
	CacheEnabled bool
}
