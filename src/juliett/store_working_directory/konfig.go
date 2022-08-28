package store_working_directory

import "github.com/friedenberg/zit/src/delta/konfig"

type Konfig struct {
	konfig.Konfig
	CacheEnabled bool
}
