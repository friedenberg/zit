package metadatei

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
)

const (
	ShaKeySelbstMetadatei              = "SelbstMetadatei"
	ShaKeySelbstMetadateiSansTai       = "SelbstMetadateiMutterSansTai"
	ShaKeySelbstMetadateiKennungMutter = "SelbstMetadateiKennungMutter"
	ShaKeyMutterMetadateiKennungMutter = "MutterMetadateiKennungMutter"
	ShaKeySelbst                       = "MetadateiTai"
	ShaKeyMutter                       = "MutterMetadateiMutterTai"
)

type Sha struct {
	*sha.Sha
	string
}

// TODO make this a map
type Shas struct {
	Akte                         sha.Sha
	SelbstMetadatei              sha.Sha
	SelbstMetadateiSansTai       sha.Sha
	SelbstMetadateiKennungMutter sha.Sha
	MutterMetadateiKennungMutter sha.Sha
}

func (s *Shas) Reset() {
	s.Akte.Reset()
	s.SelbstMetadatei.Reset()
	s.SelbstMetadateiSansTai.Reset()
	s.SelbstMetadateiKennungMutter.Reset()
	s.MutterMetadateiKennungMutter.Reset()
}

func (a *Shas) ResetWith(b *Shas) {
	a.Akte.ResetWith(&b.Akte)
	a.SelbstMetadatei.ResetWith(&b.SelbstMetadatei)
	a.SelbstMetadateiSansTai.ResetWith(&b.SelbstMetadateiSansTai)
	a.SelbstMetadateiKennungMutter.ResetWith(&b.SelbstMetadateiKennungMutter)
	a.MutterMetadateiKennungMutter.ResetWith(&b.MutterMetadateiKennungMutter)
}

func (s *Shas) String() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "%s: %s\n", "Akte", &s.Akte)
	fmt.Fprintf(&sb, "%s: %s\n", ShaKeySelbstMetadatei, &s.SelbstMetadatei)
	fmt.Fprintf(&sb, "%s: %s\n", ShaKeySelbstMetadateiSansTai, &s.SelbstMetadateiSansTai)

	return sb.String()
}

func (s *Shas) Add(k, v string) (err error) {
	switch k {
	case ShaKeySelbstMetadatei:
		if err = s.SelbstMetadatei.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

	case ShaKeySelbstMetadateiSansTai:
		if err = s.SelbstMetadateiSansTai.Set(v); err != nil {
			err = errors.Wrap(err)
			return
		}

	case ShaKeySelbstMetadateiKennungMutter:
		if err = s.SelbstMetadateiKennungMutter.Set(v); err != nil {
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
