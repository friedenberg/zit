package zettel

import (
	"github.com/friedenberg/zit/alfa/bezeichnung"
	"github.com/friedenberg/zit/bravo/akte_ext"
	"github.com/friedenberg/zit/bravo/sha"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/delta/objekte"
)

type (
	_Bezeichnung   = bezeichnung.Bezeichnung
	_EtikettSet    = etikett.Set
	_AkteExt       = akte_ext.AkteExt
	_Sha           = sha.Sha
	_ObjekteWriter = objekte.Writer
)

var (
	_EtikettNewSet = etikett.NewSet
)
