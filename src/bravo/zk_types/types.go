package zk_types

import (
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Type int

const (
	TypeUnknown = Type(iota)
	TypeAkte
	TypeAkteTyp
	TypeBezeichnung
	TypeEtikett
	TypeHinweis
	TypeTransaktion
	TypeZettel

	TypeTyp = TypeAkteTyp
)

func (t Type) String() string {
	switch t {
	case TypeAkte:
		return "Akte"

	case TypeAkteTyp:
		return "AkteTyp"

	case TypeEtikett:
		return "Etikett"

	case TypeZettel:
		return "Zettel"

	case TypeBezeichnung:
		return "Bezeichnung"

	case TypeHinweis:
		return "Hinweis"

	case TypeTransaktion:
		return "Transaktion"

	default:
		return "Unknown"
	}
}

func (t *Type) Set(v string) (err error) {
	v1 := strings.ToLower(v)
	switch v1 {
	case "akte":
		*t = TypeAkte

	case "typ":
		*t = TypeTyp

	case "aktetyp":
		*t = TypeAkteTyp

	case "etikett":
		*t = TypeEtikett

	case "zettel":
		*t = TypeZettel

	case "bezeichnung":
		*t = TypeBezeichnung

	case "hinweis":
		*t = TypeHinweis

	case "transaktion":
		*t = TypeTransaktion

	default:
		err = errors.Errorf("unknown object type: %s", v)
	}

	return
}
