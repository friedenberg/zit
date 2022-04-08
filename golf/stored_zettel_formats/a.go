package stored_zettel_formats

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/alfa/node_type"
	"github.com/friedenberg/zit/alfa/bezeichnung"
	"github.com/friedenberg/zit/bravo/akte_ext"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

const (
	_TypeMutter      = node_type.TypeMutter
	_TypeKinder      = node_type.TypeKinder
	_TypeAkte        = node_type.TypeAkte
	_TypeAkteExt     = node_type.TypeAkteExt
	_TypeBezeichnung = node_type.TypeBezeichnung
	_TypeEtikett     = node_type.TypeEtikett
)

type (
	_StoredZettel = stored_zettel.Stored
	_Type         = node_type.Type
	_Etikett      = etikett.Etikett
	_Zettel       = zettel.Zettel
	_EtikettSet   = etikett.Set
	_AkteExt      = akte_ext.AkteExt
	_Sha          = sha.Sha
	_Bezeichnung  = bezeichnung.Bezeichnung
)

var (
	_Errorf        = errors.Errorf
	_Error         = errors.Error
	_EtikettNewSet = etikett.NewSet
)
