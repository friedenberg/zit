package akten

import (
	"code.linenisgreat.com/zit-go/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit-go/src/delta/etikett_akte"
	"code.linenisgreat.com/zit-go/src/delta/kasten_akte"
	"code.linenisgreat.com/zit-go/src/delta/standort"
	"code.linenisgreat.com/zit-go/src/delta/typ_akte"
	"code.linenisgreat.com/zit-go/src/india/erworben"
	"code.linenisgreat.com/zit-go/src/juliett/objekte"
	"code.linenisgreat.com/zit-go/src/kilo/objekte_store"
)

type Store[
	A schnittstellen.Akte[A],
	APtr schnittstellen.AktePtr[A],
] interface {
	SaveAkteText(APtr) (schnittstellen.ShaLike, int64, error)
	objekte_store.StoredParseSaver[A, APtr]
	objekte.AkteFormat[A, APtr]
	schnittstellen.AkteGetterPutter[APtr]
}

type Akten struct {
	etikett_v0 Store[etikett_akte.V0, *etikett_akte.V0]
	kasten_v0  Store[kasten_akte.V0, *kasten_akte.V0]
	konfig_v0  Store[erworben.Akte, *erworben.Akte]
	typ_v0     Store[typ_akte.V0, *typ_akte.V0]
}

func Make(
	st standort.Standort,
) *Akten {
	return &Akten{
		etikett_v0: objekte_store.MakeAkteStore[etikett_akte.V0, *etikett_akte.V0](
			st,
			objekte_store.MakeAkteFormat[etikett_akte.V0, *etikett_akte.V0](
				objekte.MakeTextParserIgnoreTomlErrors[etikett_akte.V0](
					st,
				),
				objekte.ParsedAkteTomlFormatter[etikett_akte.V0, *etikett_akte.V0]{},
				st,
			),
			func(a *etikett_akte.V0) {
				a.Reset()
			},
		),
		kasten_v0: objekte_store.MakeAkteStore[kasten_akte.V0, *kasten_akte.V0](
			st,
			objekte_store.MakeAkteFormat[kasten_akte.V0, *kasten_akte.V0](
				objekte.MakeTextParserIgnoreTomlErrors[kasten_akte.V0](
					st,
				),
				objekte.ParsedAkteTomlFormatter[kasten_akte.V0, *kasten_akte.V0]{},
				st,
			),
			func(a *kasten_akte.V0) {
				a.Reset()
			},
		),
		konfig_v0: objekte_store.MakeAkteStore[erworben.Akte, *erworben.Akte](
			st,
			objekte_store.MakeAkteFormat[erworben.Akte, *erworben.Akte](
				objekte.MakeTextParserIgnoreTomlErrors[erworben.Akte](
					st,
				),
				objekte.ParsedAkteTomlFormatter[erworben.Akte, *erworben.Akte]{},
				st,
			),
			func(a *erworben.Akte) {
				a.Reset()
			},
		),
		typ_v0: objekte_store.MakeAkteStore[typ_akte.V0, *typ_akte.V0](
			st,
			objekte_store.MakeAkteFormat[typ_akte.V0, *typ_akte.V0](
				objekte.MakeTextParserIgnoreTomlErrors[typ_akte.V0](
					st,
				),
				objekte.ParsedAkteTomlFormatter[typ_akte.V0, *typ_akte.V0]{},
				st,
			),
			func(a *typ_akte.V0) {
				a.Reset()
			},
		),
	}
}

func (a *Akten) GetEtikettV0() Store[etikett_akte.V0, *etikett_akte.V0] {
	return a.etikett_v0
}

func (a *Akten) GetKastenV0() Store[kasten_akte.V0, *kasten_akte.V0] {
	return a.kasten_v0
}

func (a *Akten) GetKonfigV0() Store[erworben.Akte, *erworben.Akte] {
	return a.konfig_v0
}

func (a *Akten) GetTypV0() Store[typ_akte.V0, *typ_akte.V0] {
	return a.typ_v0
}
