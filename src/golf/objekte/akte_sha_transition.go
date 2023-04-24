package objekte

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

type InNeedOfAkteShaCorrection interface {
	metadatei.Getter
	GetAkteSha() schnittstellen.Sha
	GetSkuAkteSha() schnittstellen.Sha
}

type InNeedOfAkteShaCorrectionPtr interface {
	InNeedOfAkteShaCorrection
	metadatei.Setter
	SetAkteSha(schnittstellen.Sha)
}

func CorrectAkteSha(
	inoasc InNeedOfAkteShaCorrectionPtr,
	corrector metadatei.Getter,
) {
	if inoasc == nil {
		panic("InNeedOfAkteShaCorrection was nil")
	}

	if corrector == nil {
		panic("metadatei.Getter was nil")
	}

	inoasc.SetAkteSha(corrector.GetMetadatei().AkteSha)
}

func AssertAkteShasMatch(inoa InNeedOfAkteShaCorrection) {
	shSku := inoa.GetSkuAkteSha()
	shMetadatei := inoa.GetMetadatei().AkteSha

	if !shSku.EqualsSha(shMetadatei) {
		panic(errors.Errorf(
			"akte sha in sku was %s while akte sha in metadatei was %s",
			shSku,
			shMetadatei,
		))
	}
}
