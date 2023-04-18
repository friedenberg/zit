package metadatei

import (
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/delta/kennung"
)

type Getter interface {
	GetMetadatei() Metadatei
}

type Setter interface {
	SetMetadatei(Metadatei)
}

type PersistentFormatterContext interface {
	Getter
	GetAkteSha() schnittstellen.Sha
}

type PersistentParserContext interface {
	Getter
	Setter
}

type TextFormatterContext interface {
	PersistentFormatterContext
	GetAktePath() string
}

type TextParserContext interface {
	PersistentParserContext
	SetAkteFD(kennung.FD) error
}
