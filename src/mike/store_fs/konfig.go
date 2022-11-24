package store_fs

import "github.com/friedenberg/zit/src/echo/konfig"

type Konfig struct {
	konfig.Konfig
	CacheEnabled bool
}
