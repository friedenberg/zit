package stored_zettel

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/echo/zettel"
)

type (
	_Zettel  = zettel.Zettel
	_Sha     = sha.Sha
	_Hinweis = hinweis.Hinweis

	_EtikettSet = etikett.Set
)

var (
	_Error = errors.Error
)
