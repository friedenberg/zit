package store_fs

import "github.com/friedenberg/zit/src/delta/konfig"

type Konfig struct {
	konfig.Konfig
	CacheEnabled bool
}
