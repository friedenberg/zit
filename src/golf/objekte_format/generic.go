package objekte_format

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/ohio"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type (
	Metadatei = metadatei.Metadatei
	Sha       = sha.Sha
)

var (
	keyAkte        = catgut.MakeFromString("Akte")
	keyBezeichnung = catgut.MakeFromString("Bezeichnung")
	keyEtikett     = catgut.MakeFromString("Etikett")
	keyGattung     = catgut.MakeFromString("Gattung")
	keyKennung     = catgut.MakeFromString("Kennung")
	keyKomment     = catgut.MakeFromString("Komment")
	keyTai         = catgut.MakeFromString("Tai")
	keyTyp         = catgut.MakeFromString("Typ")

	keyMutter = catgut.MakeFromString("zzMutter")
	keySha    = catgut.MakeFromString("zzSha")

	keyVerzeichnisseArchiviert      = catgut.MakeFromString("Verzeichnisse-Archiviert")
	keyVerzeichnisseEtikettImplicit = catgut.MakeFromString("Verzeichnisse-Etikett-Implicit")
	keyVerzeichnisseEtikettExpanded = catgut.MakeFromString("Verzeichnisse-Etikett-Expanded")
)

type FormatGeneric struct {
	key  string
	keys []*catgut.String
}

var FormatsGeneric = map[string][]*catgut.String{
	// "Akte":                {keyAkte},
	// "AkteBez":             {keyAkte, keyBezeichnung},
	// "AkteTyp":             {keyAkte, keyTyp},
	// "MetadateiSansTai":    {keyAkte, keyBezeichnung, keyEtikett, keyTyp},
	// "Metadatei":           {keyAkte, keyBezeichnung, keyEtikett, keyTyp, keyTai},
	"MetadateiPlusMutter": {keyAkte, keyBezeichnung, keyEtikett, keyTyp, keyTai, keyMutter},
}

func FormatForKeyError(k string) (fo FormatGeneric, err error) {
	f, ok := FormatsGeneric[k]

	if !ok {
		err = errInvalidGenericFormat(k)
		return
	}

	fo = FormatGeneric{
		key:  k,
		keys: f,
	}

	return
}

func FormatForKey(k string) FormatGeneric {
	f, err := FormatForKeyError(k)
	errors.PanicIfError(err)
	return f
}

func (f FormatGeneric) printKeys(
	w io.Writer,
	m *Metadatei,
) (n int64, err error) {
	var n1 int64

	for _, k := range f.keys {
		n1, err = f.printKey(w, m, k)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (f FormatGeneric) printKey(
	w io.Writer,
	m *Metadatei,
	key *catgut.String,
) (n int64, err error) {
	var n1 int

	switch key {
	case keyAkte:
		if !m.Akte.IsNull() {
			n1, err = ohio.WriteKeySpaceValueNewlineString(
				w,
				keyAkte.String(),
				m.Sha.String(),
			)
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	case keyBezeichnung:
		lines := strings.Split(m.Bezeichnung.String(), "\n")

		for _, line := range lines {
			if line == "" {
				continue
			}

			n1, err = ohio.WriteKeySpaceValueNewlineString(
				w,
				keyBezeichnung.String(),
				line,
			)
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	case keyEtikett:
		es := m.GetEtiketten()

		for _, e := range iter.SortedValues[kennung.Etikett](es) {
			n1, err = ohio.WriteKeySpaceValueNewlineString(
				w,
				keyEtikett.String(),
				e.String(),
			)
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	case keyMutter:

		if !m.Mutter.IsNull() {
			n1, err = ohio.WriteKeySpaceValueNewlineString(
				w,
				keyMutter.String(),
				m.Mutter.String(),
			)
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}

			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	case keyTai:
		n1, err = ohio.WriteKeySpaceValueNewlineString(
			w,
			keyTai.String(),
			m.Tai.String(),
		)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

	case keyTyp:
		if !m.Typ.IsEmpty() {
			n1, err = ohio.WriteKeySpaceValueNewlineString(
				w,
				keyTyp.String(),
				m.GetTyp().String(),
			)
			n += int64(n1)

			if err != nil {
				err = errors.Wrap(err)
				return
			}
		}

	default:
		panic(fmt.Sprintf("unsupported key: %s", key))
	}

	return
}

func GetShaForMetadatei(f FormatGeneric, m *Metadatei) (sh *Sha, err error) {
	switch f.key {
	case "Akte", "AkteTyp":
		if m.Akte.IsNull() {
			return
		}

	case "AkteBez":
		if m.Akte.IsNull() && m.Bezeichnung.IsEmpty() {
			return
		}
	}

	var sb strings.Builder
	sw := sha.MakeWriter(&sb)

	_, err = f.printKeys(sw, m)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = &Sha{}

	if err = sh.SetShaLike(sw); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func GetShasForMetadatei(m *Metadatei) (shas map[string]*sha.Sha, err error) {
	shas = make(map[string]*sha.Sha, len(FormatsGeneric))

	for k := range FormatsGeneric {
		f := FormatForKey(k)

		var sh *Sha

		if sh, err = GetShaForMetadatei(f, m); err != nil {
			err = errors.Wrap(err)
			return
		}

		if sh == nil {
			continue
		}

		shas[k] = sh
	}

	return
}
