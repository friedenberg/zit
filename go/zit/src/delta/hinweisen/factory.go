package hinweisen

import (
	"path"
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
)

const (
	FilePathKennungYin  = "Yin"
	FilePathKennungYang = "Yang"
)

type Hinweisen struct {
	sync.Locker
	yin  provider
	yang provider
}

func New(ps schnittstellen.Standort) (f *Hinweisen, err error) {
	providerPathYin := path.Join(ps.DirKennung(), FilePathKennungYin)
	providerPathYang := path.Join(ps.DirKennung(), FilePathKennungYang)

	f = &Hinweisen{
		Locker: &sync.Mutex{},
	}

	if f.yin, err = newProvider(providerPathYin); err != nil {
		err = errors.Wrap(err)
		return
	}

	if f.yang, err = newProvider(providerPathYang); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (hf *Hinweisen) Left() provider {
	return hf.yin
}

func (hf *Hinweisen) Right() provider {
	return hf.yang
}
