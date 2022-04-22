package organize_text

import (
	"github.com/friedenberg/zit/alfa/bezeichnung"
	"github.com/friedenberg/zit/alfa/errors"
	"github.com/friedenberg/zit/bravo/line_format"
	"github.com/friedenberg/zit/charlie/etikett"
	"github.com/friedenberg/zit/charlie/hinweis"
	zettel1 "github.com/friedenberg/zit/echo/zettel"
	"github.com/friedenberg/zit/foxtrot/stored_zettel"
	"github.com/friedenberg/zit/foxtrot/zettel_formats"
)

const (
	_MetadateiBoundary = zettel_formats.MetadateiBoundary
)

type (
	_NamedZettel          = stored_zettel.Named
	_Zettel               = zettel1.Zettel
	_Etikett              = etikett.Etikett
	_EtikettSet           = etikett.Set
	_Hinweis              = hinweis.Hinweis
	_Bezeichnung          = bezeichnung.Bezeichnung
	_EtikettExpanderRight = etikett.ExpanderRight
)

var (
	_Error                  = errors.Error
	_HinweisNewEmpty        = hinweis.NewEmpty
	_Errorf                 = errors.Errorf
	_LineFormatNewWriter    = line_format.NewWriter
	_EtikettNewSetFromSlice = etikett.NewSetFromSlice
	_EtikettNewSet          = etikett.NewSet
)
