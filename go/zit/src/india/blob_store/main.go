package blob_store

import (
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/etikett_akte"
	"code.linenisgreat.com/zit/go/zit/src/delta/typ_akte"
	"code.linenisgreat.com/zit/go/zit/src/echo/kasten_akte"
	"code.linenisgreat.com/zit/go/zit/src/echo/standort"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/erworben"
)

type Store[
	A interfaces.Blob[A],
	APtr interfaces.BlobPtr[A],
] interface {
	SaveBlobText(APtr) (interfaces.ShaLike, int64, error)
	Format[A, APtr]
	interfaces.BlobGetterPutter[APtr]
}

type VersionedStores struct {
	etikett_v0 Store[etikett_akte.V0, *etikett_akte.V0]
	etikett_v1 Store[etikett_akte.V1, *etikett_akte.V1]
	kasten_v0  Store[kasten_akte.V0, *kasten_akte.V0]
	konfig_v0  Store[erworben.Akte, *erworben.Akte]
	typ_v0     Store[typ_akte.V0, *typ_akte.V0]
}

func Make(
	st standort.Standort,
) *VersionedStores {
	return &VersionedStores{
		etikett_v0: MakeBlobStore(
			st,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[etikett_akte.V0](
					st,
				),
				ParsedBlobTomlFormatter[etikett_akte.V0, *etikett_akte.V0]{},
				st,
			),
			func(a *etikett_akte.V0) {
				a.Reset()
			},
		),
		etikett_v1: MakeBlobStore(
			st,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[etikett_akte.V1](
					st,
				),
				ParsedBlobTomlFormatter[etikett_akte.V1, *etikett_akte.V1]{},
				st,
			),
			func(a *etikett_akte.V1) {
				a.Reset()
			},
		),
		kasten_v0: MakeBlobStore(
			st,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[kasten_akte.V0](
					st,
				),
				ParsedBlobTomlFormatter[kasten_akte.V0, *kasten_akte.V0]{},
				st,
			),
			func(a *kasten_akte.V0) {
				a.Reset()
			},
		),
		konfig_v0: MakeBlobStore(
			st,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[erworben.Akte](
					st,
				),
				ParsedBlobTomlFormatter[erworben.Akte, *erworben.Akte]{},
				st,
			),
			func(a *erworben.Akte) {
				a.Reset()
			},
		),
		typ_v0: MakeBlobStore(
			st,
			MakeBlobFormat(
				MakeTextParserIgnoreTomlErrors[typ_akte.V0](
					st,
				),
				ParsedBlobTomlFormatter[typ_akte.V0, *typ_akte.V0]{},
				st,
			),
			func(a *typ_akte.V0) {
				a.Reset()
			},
		),
	}
}

func (a *VersionedStores) GetEtikettV0() Store[etikett_akte.V0, *etikett_akte.V0] {
	return a.etikett_v0
}

func (a *VersionedStores) GetEtikettV1() Store[etikett_akte.V1, *etikett_akte.V1] {
	return a.etikett_v1
}

func (a *VersionedStores) GetKastenV0() Store[kasten_akte.V0, *kasten_akte.V0] {
	return a.kasten_v0
}

func (a *VersionedStores) GetKonfigV0() Store[erworben.Akte, *erworben.Akte] {
	return a.konfig_v0
}

func (a *VersionedStores) GetTypV0() Store[typ_akte.V0, *typ_akte.V0] {
	return a.typ_v0
}
