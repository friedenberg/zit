package node_type

import (
	"strings"

	"github.com/friedenberg/zit/alfa/errors"
)

type Type int

const (
	TypeUnknown = Type(iota)
	TypeAkte
	TypeAkteExt
	TypeEtikett
	TypeZettel
	TypeMutter
	TypeKinder
	TypeBezeichnung
	TypeHinweis
)

func (t Type) String() string {
	switch t {
	case TypeAkte:
		return "Akte"

	case TypeAkteExt:
		return "AkteExt"

	case TypeEtikett:
		return "Etikett"

	case TypeZettel:
		return "Zettel"

	case TypeMutter:
		return "Mutter"

	case TypeKinder:
		return "Kinder"

	case TypeBezeichnung:
		return "Bezeichnung"

	case TypeHinweis:
		return "Hinweis"

	default:
		return "Unknown"
	}
}

func (t *Type) Set(v string) (err error) {
	v1 := strings.ToLower(v)
	switch v1 {
	case "akte":
		*t = TypeAkte

	case "akteext":
		*t = TypeAkteExt

	case "etikett":
		*t = TypeEtikett

	case "zettel":
		*t = TypeZettel

	case "mutter":
		*t = TypeMutter

	case "kinder":
		*t = TypeKinder

	case "bezeichnung":
		*t = TypeBezeichnung

	case "hinweis":
		*t = TypeHinweis

	default:
		err = errors.Errorf("unknown object type: %s", v)
	}

	return
}
