package akten

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/go/zit/src/delta/etikett_akte"
	"code.linenisgreat.com/zit/go/zit/src/delta/typ_akte"
	"code.linenisgreat.com/zit/go/zit/src/echo/kasten_akte"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/erworben"
)

type Store[
	A schnittstellen.Akte[A],
	APtr schnittstellen.AktePtr[A],
] interface {
	SaveAkteText(APtr) (schnittstellen.ShaLike, int64, error)
	StoredParseSaver[A, APtr]
	Format[A, APtr]
	schnittstellen.AkteGetterPutter[APtr]
}

type Akten struct {
	etikett_v0 Store[etikett_akte.V0, *etikett_akte.V0]
	etikett_v1 Store[etikett_akte.V1, *etikett_akte.V1]
	kasten_v0  Store[kasten_akte.V0, *kasten_akte.V0]
	konfig_v0  Store[erworben.Akte, *erworben.Akte]
	typ_v0     Store[typ_akte.V0, *typ_akte.V0]
}

func Make(
	st standort.Standort,
) *Akten {
	return &Akten{
		etikett_v0: MakeAkteStore(
			st,
			MakeAkteFormat(
				MakeTextParserIgnoreTomlErrors[etikett_akte.V0](
					st,
				),
				ParsedAkteTomlFormatter[etikett_akte.V0, *etikett_akte.V0]{},
				st,
			),
			func(a *etikett_akte.V0) {
				a.Reset()
			},
		),
		etikett_v1: MakeAkteStore(
			st,
			MakeAkteFormat(
				MakeTextParserIgnoreTomlErrors[etikett_akte.V1](
					st,
				),
				ParsedAkteTomlFormatter[etikett_akte.V1, *etikett_akte.V1]{},
				st,
			),
			func(a *etikett_akte.V1) {
				a.Reset()
			},
		),
		kasten_v0: MakeAkteStore(
			st,
			MakeAkteFormat(
				MakeTextParserIgnoreTomlErrors[kasten_akte.V0](
					st,
				),
				ParsedAkteTomlFormatter[kasten_akte.V0, *kasten_akte.V0]{},
				st,
			),
			func(a *kasten_akte.V0) {
				a.Reset()
			},
		),
		konfig_v0: MakeAkteStore(
			st,
			MakeAkteFormat(
				MakeTextParserIgnoreTomlErrors[erworben.Akte](
					st,
				),
				ParsedAkteTomlFormatter[erworben.Akte, *erworben.Akte]{},
				st,
			),
			func(a *erworben.Akte) {
				a.Reset()
			},
		),
		typ_v0: MakeAkteStore(
			st,
			MakeAkteFormat(
				MakeTextParserIgnoreTomlErrors[typ_akte.V0](
					st,
				),
				ParsedAkteTomlFormatter[typ_akte.V0, *typ_akte.V0]{},
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

func (a *Akten) GetEtikettV1() Store[etikett_akte.V1, *etikett_akte.V1] {
	return a.etikett_v1
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
