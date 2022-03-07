package alfred

import (
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/alfred"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/charlie/hinweis"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
)

type (
	_NamedZettel  = stored_zettel.Named
	_AlfredItem   = alfred.Item
	_AlfredWriter = alfred.Writer
	_Etikett      = etikett.Etikett
	_EtikettSet   = etikett.Set
	_Hinweis      = hinweis.Hinweis
)

var (
	_AlfredNewMatchBuilder = alfred.NewMatchBuilder
	_AlfredNewWriter       = alfred.NewWriter
	_Error                 = errors.Error
)
