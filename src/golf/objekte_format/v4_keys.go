package objekte_format

import "github.com/friedenberg/zit/src/charlie/catgut"

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
