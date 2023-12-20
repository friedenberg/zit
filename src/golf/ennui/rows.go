package ennui

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/catgut"
	"github.com/friedenberg/zit/src/charlie/sha"
)

var formats = map[string][]*catgut.String{
	"Akte":                {keyAkte},
	"AkteBez":             {keyAkte, keyBezeichnung},
	"AkteTyp":             {keyAkte, keyTyp},
	"MetadateiSansTai":    {keyAkte, keyBezeichnung, keyEtikett, keyTyp},
	"Metadatei":           {keyAkte, keyBezeichnung, keyEtikett, keyTyp, keyTai},
	"MetadateiPlusMutter": {keyAkte, keyBezeichnung, keyEtikett, keyTyp, keyTai, keyMutter},
}

func (e *ennui) getShasForMetadatei(m *Metadatei) (shas map[string]*sha.Sha, err error) {
	shas = make(map[string]*sha.Sha, len(formats))

	for k, f1 := range formats {
		f := format{
			key:  k,
			keys: f1,
		}

		switch k {
		case "Akte", "AkteTyp":
			if m.Akte.IsNull() {
				continue
			}

		case "AkteBez":
			if m.Akte.IsNull() && m.Bezeichnung.IsEmpty() {
				continue
			}
		}

		sw := sha.MakeWriter(nil)

		_, err = f.printKeys(sw, m)

		if err != nil {
			err = errors.Wrap(err)
			return
		}

		var sh sha.Sha

		if err = sh.SetShaLike(sw); err != nil {
			err = errors.Wrap(err)
			return
		}

    // log.Debug().Printf("%s: %s", k, &sh)
		shas[k] = &sh
	}

	return
}
