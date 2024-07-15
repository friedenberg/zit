package object_inventory_format

import (
	"fmt"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/catgut"
	"code.linenisgreat.com/zit/go/zit/src/delta/german_keys"
	"code.linenisgreat.com/zit/go/zit/src/delta/keys"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/object_metadata"
)

type (
	Metadatei = object_metadata.Metadata
	Sha       = sha.Sha
)

const (
	KeyFormatMetadata               = "Metadatei"
	KeyFormatMetadataWithoutTai     = "MetadateiSansTai"
	KeyFormatMetadataObjectIdParent = "MetadateiKennungMutter"
)

var (
	keyAkte                             = german_keys.KeyAkte
	keyBezeichnung                      = german_keys.KeyBezeichnung
	keyEtikett                          = german_keys.KeyEtikett
	keyGattung                          = german_keys.KeyGattung
	keyKennung                          = german_keys.KeyKennung
	keyKomment                          = german_keys.KeyKomment
	keyTyp                              = german_keys.KeyTyp
	keyShasMutterMetadateiKennungMutter = german_keys.KeyShasMutterMetadateiKennungMutter
	keyVerzeichnisseArchiviert          = german_keys.KeyVerzeichnisseArchiviert
	keyVerzeichnisseEtikettImplicit     = german_keys.KeyVerzeichnisseEtikettImplicit
	keyVerzeichnisseEtikettExpanded     = german_keys.KeyVerzeichnisseEtikettExpanded

	keySigil = keys.KeySigil
	keyTai   = keys.KeyTai
	keySha   = keys.KeySha
)

type FormatGeneric struct {
	key  string
	keys []*catgut.String
}

var FormatsGeneric = map[string][]*catgut.String{
	KeyFormatMetadataWithoutTai:     {keyAkte, keyBezeichnung, keyEtikett, keyTyp},
	KeyFormatMetadata:               {keyAkte, keyBezeichnung, keyEtikett, keyTyp, keyTai},
	KeyFormatMetadataObjectIdParent: {keyAkte, keyBezeichnung, keyEtikett, keyKennung, keyTyp, keyTai, keyShasMutterMetadateiKennungMutter},
}

type formats struct {
	metadateiSansTai       FormatGeneric
	metadatei        FormatGeneric
	metadateiKennungMutter FormatGeneric
}

func (fs formats) MetadateiSansTai() FormatGeneric {
	return fs.metadateiSansTai
}

func (fs formats) Metadatei() FormatGeneric {
	return fs.metadatei
}

func (fs formats) MetadateiKennungMutter() FormatGeneric {
	return fs.metadateiKennungMutter
}

var Formats formats

func init() {
	Formats.metadatei.key = KeyFormatMetadata
	Formats.metadatei.keys = FormatsGeneric[KeyFormatMetadata]

	Formats.metadateiSansTai.key = KeyFormatMetadataWithoutTai
	Formats.metadateiSansTai.keys = FormatsGeneric[KeyFormatMetadataWithoutTai]

	Formats.metadateiKennungMutter.key = KeyFormatMetadataObjectIdParent
	Formats.metadateiKennungMutter.keys = FormatsGeneric[KeyFormatMetadataObjectIdParent]
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
	m := c.GetMetadata()

	var n1 int

	switch key {
	case keyAkte:
		n1, err = writeShaKeyIfNotNull(
			w,
			keyAkte,
			&m.Blob,
		)

		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

	case keyBezeichnung:
		lines := strings.Split(m.Description.String(), "\n")

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
		es := m.GetTags()

		if es == nil {
			break
		}

		// TODO fix issue with es being nil sometimes
		for _, e := range iter.SortedValues(es) {
			if e.IsVirtual() {
				continue
			}

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
			c.GetObjectId().GetGenre().GetGenreString(),
		)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		n1, err = ohio.WriteKeySpaceValueNewlineString(
			w,
			keyKennung.String(),
			c.GetObjectId().String(),
		)
		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

	case keyShasMutterMetadateiKennungMutter:
		n1, err = writeShaKeyIfNotNull(
			w,
			keyShasMutterMetadateiKennungMutter,
			m.Mutter(),
		)

		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

	case keyShasMutterMetadateiKennungMutter:
		n1, err = writeShaKeyIfNotNull(
			w,
			keyShasMutterMetadateiKennungMutter,
			&m.ParentMetadataObjectIdParent,
		)

		n += int64(n1)

		if err != nil {
			err = errors.Wrap(err)
			return
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
		if !m.Type.IsEmpty() {
			n1, err = ohio.WriteKeySpaceValueNewlineString(
				w,
				keyTyp.String(),
				m.GetType().String(),
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

func writeShaKeyIfNotNull(
	w io.Writer,
	key *catgut.String,
	sh *sha.Sha,
) (n int, err error) {
	if sh.IsNull() {
		return
	}

	n, err = ohio.WriteKeySpaceValueNewlineString(
		w,
		key.String(),
		sh.String(),
	)
	if err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func GetShaForContext(f FormatGeneric, c FormatterContext) (sh *Sha, err error) {
	m := c.GetMetadata()

	switch f.key {
	case "Akte", "AkteTyp":
		if m.Blob.IsNull() {
			return
		}

	case "AkteBez":
		if m.Blob.IsNull() && m.Description.IsEmpty() {
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

func GetShaForContextDebug(
	f FormatGeneric,
	c FormatterContext,
) (sh *Sha, err error) {
	var sb strings.Builder
	sw := sha.MakeWriter(&sb)

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

	ui.DebugAllowCommit().Caller(2, "%s:%s -> %s", f.key, sh, &sb)

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
