package objekte_format

import (
	"fmt"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/log"
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
	keyContentLength = catgut.MakeFromString("ContentLength")
	keyAkte          = catgut.MakeFromString("Akte")
	keyBezeichnung   = catgut.MakeFromString("Bezeichnung")
	keyEtikett       = catgut.MakeFromString("Etikett")
	keyGattung       = catgut.MakeFromString("Gattung")
	keyKennung       = catgut.MakeFromString("Kennung")
	keyKomment       = catgut.MakeFromString("Komment")
	keySigil         = catgut.MakeFromString("Sigil")
	keyTai           = catgut.MakeFromString("Tai")
	keyTyp           = catgut.MakeFromString("Typ")

	keyMutter = catgut.MakeFromString("zzMutter")
	keySha    = catgut.MakeFromString("zzSha")

	keyShasMutterMetadatei              = catgut.MakeFromString("Shas" + metadatei.ShaKeyMutterMetadatei)
	keyShasMutterMetadateiMutter        = catgut.MakeFromString("Shas" + metadatei.ShaKeyMutterMetadateiMutter)
	keyShasMutterMetadateiKennungMutter = catgut.MakeFromString("Shas" + metadatei.ShaKeyMutterMetadateiKennungMutter)

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
	"Metadatei":              {keyAkte, keyBezeichnung, keyEtikett, keyTyp, keyTai},
	"MetadateiMutter":        {keyAkte, keyBezeichnung, keyEtikett, keyTyp, keyTai, keyShasMutterMetadateiMutter},
	"MetadateiKennungMutter": {keyAkte, keyBezeichnung, keyEtikett, keyKennung, keyTyp, keyTai, keyShasMutterMetadateiKennungMutter},
}

type formats struct {
	metadatei              FormatGeneric
	metadateiMutter        FormatGeneric
	metadateiKennungMutter FormatGeneric
}

func (fs formats) Metadatei() FormatGeneric {
	return fs.metadatei
}

func (fs formats) MetadateiMutter() FormatGeneric {
	return fs.metadateiMutter
}

func (fs formats) MetadateiKennungMutter() FormatGeneric {
	return fs.metadateiKennungMutter
}

var Formats formats

func init() {
	Formats.metadatei.key = "Metadatei"
	Formats.metadatei.keys = FormatsGeneric["Metadatei"]

	Formats.metadateiMutter.key = "MetadateiMutter"
	Formats.metadateiMutter.keys = FormatsGeneric["MetadateiMutter"]

	Formats.metadateiKennungMutter.key = "MetadateiKennungMutter"
	Formats.metadateiKennungMutter.keys = FormatsGeneric["MetadateiKennungMutter"]
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

func (f FormatGeneric) WriteMetadateiTo(
	w io.Writer,
	c FormatterContext,
) (n int64, err error) {
	var n1 int64

	for _, k := range f.keys {
		n1, err = WriteMetadateiKeyTo(w, c, k)
		n += n1

		if err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func WriteMetadateiKeyTo(
	w io.Writer,
	c FormatterContext,
	key *catgut.String,
) (n int64, err error) {
	m := c.GetMetadatei()

	var n1 int

	switch key {
	case keyAkte:
		if !m.Akte.IsNull() {
			n1, err = ohio.WriteKeySpaceValueNewlineString(
				w,
				keyAkte.String(),
				m.Akte.String(),
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

	case keyKennung:
		n1, err = ohio.WriteKeySpaceValueNewlineString(
			w,
			keyGattung.String(),
			c.GetKennung().GetGattung().GetGattungString(),
		)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n1, err = ohio.WriteKeySpaceValueNewlineString(
			w,
			keyKennung.String(),
			c.GetKennung().String(),
		)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
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

	case keyShasMutterMetadatei:
		if !m.MutterMetadatei.IsNull() {
			n1, err = ohio.WriteKeySpaceValueNewlineString(
				w,
				keyShasMutterMetadatei.String(),
				m.MutterMetadatei.String(),
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

	case keyShasMutterMetadateiMutter:
		if !m.MutterMetadatei.IsNull() {
			n1, err = ohio.WriteKeySpaceValueNewlineString(
				w,
				keyShasMutterMetadateiMutter.String(),
				m.MutterMetadateiMutter.String(),
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

	case keyShasMutterMetadateiKennungMutter:
		if !m.MutterMetadatei.IsNull() {
			n1, err = ohio.WriteKeySpaceValueNewlineString(
				w,
				keyShasMutterMetadateiKennungMutter.String(),
				m.MutterMetadateiKennungMutter.String(),
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

func GetShaForContext(f FormatGeneric, c FormatterContext) (sh *Sha, err error) {
	m := c.GetMetadatei()

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

	if m.GetTai().IsEmpty() {
		err = ErrEmptyTai
		return
	}

	return getShaForContext(f, c)
}

func GetShaForMetadatei(f FormatGeneric, m *Metadatei) (sh *Sha, err error) {
	return GetShaForContext(f, nopFormatterContext{m})
}

func WriteMetadatei(w io.Writer, f FormatGeneric, c FormatterContext) (sh *Sha, err error) {
	sw := sha.MakeWriter(w)

	_, err = f.WriteMetadateiTo(sw, c)

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = sha.GetPool().Get()

	if err = sh.SetShaLike(sw); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func getShaForContext(f FormatGeneric, c FormatterContext) (sh *Sha, err error) {
	if sh, err = WriteMetadatei(nil, f, c); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func getShaForMetadateiDebug(f FormatGeneric, c *Metadatei) (sh *Sha, err error) {
	var sb strings.Builder
	sw := sha.MakeWriter(&sb)

	_, err = f.WriteMetadateiTo(sw, nopFormatterContext{c})

	if err != nil {
		err = errors.Wrap(err)
		return
	}

	sh = &Sha{}

	if err = sh.SetShaLike(sw); err != nil {
		err = errors.Wrap(err)
		return
	}

	log.DebugAllowCommit().Caller(2, "%s:%s -> %s", f.key, sh, &sb)

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
