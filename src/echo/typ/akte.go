package typ

import (
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/konfig"
)

type Akte struct {
	sha.Sha
	konfig.KonfigTyp
}
