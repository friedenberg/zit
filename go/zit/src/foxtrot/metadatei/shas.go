package metadatei

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/charlie/sha"
)

const (
	ShaKeySelbstMetadatei              = "SelbstMetadatei"
	ShaKeySelbstMetadateiMutter        = "SelbstMetadateiMutter"
	ShaKeySelbstMetadateiKennungMutter = "SelbstMetadateiKennungMutter"
	ShaKeyMutterMetadatei              = "MutterMetadatei"
	ShaKeyMutterMetadateiMutter        = "MutterMetadateiMutter"
	ShaKeyMutterMetadateiKennungMutter = "MutterMetadateiKennungMutter"
)

type Shas struct {
	SelbstMetadatei              sha.Sha
	SelbstMetadateiMutter        sha.Sha
	SelbstMetadateiKennungMutter sha.Sha
	MutterMetadatei              sha.Sha
	MutterMetadateiMutter        sha.Sha
	MutterMetadateiKennungMutter sha.Sha
}

func (s *Shas) Reset() {
	s.SelbstMetadatei.Reset()
	s.SelbstMetadateiMutter.Reset()
	s.SelbstMetadateiKennungMutter.Reset()
	s.MutterMetadatei.Reset()
	s.MutterMetadateiMutter.Reset()
	s.MutterMetadateiKennungMutter.Reset()
}

func (a *Shas) ResetWith(b *Shas) {
	a.SelbstMetadatei.ResetWith(&b.SelbstMetadatei)
	a.SelbstMetadateiMutter.ResetWith(&b.SelbstMetadateiMutter)
	a.SelbstMetadateiKennungMutter.ResetWith(&b.SelbstMetadateiKennungMutter)
	a.MutterMetadatei.ResetWith(&b.MutterMetadatei)
	a.MutterMetadateiMutter.ResetWith(&b.MutterMetadateiMutter)
	a.MutterMetadateiKennungMutter.ResetWith(&b.MutterMetadateiKennungMutter)
}

func (s *Shas) Add(k, v string) (err error) {
	switch k {
	case ShaKeySelbstMetadatei:
		if err = s.SelbstMetadatei.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

	case ShaKeySelbstMetadateiMutter:
		if err = s.SelbstMetadateiMutter.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

	case ShaKeySelbstMetadateiKennungMutter:
		if err = s.SelbstMetadateiKennungMutter.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

	case ShaKeyMutterMetadatei:
		if err = s.MutterMetadatei.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

	case ShaKeyMutterMetadateiMutter:
		if err = s.MutterMetadateiMutter.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

	case ShaKeyMutterMetadateiKennungMutter:
		if err = s.MutterMetadateiKennungMutter.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

	default:
		err = errors.Errorf("unrecognized sha kind: %q", k)
		return
	}

	return
}
