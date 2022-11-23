package hinweisen

import (
	"path"
	"sync"

	"github.com/friedenberg/zit/src/alfa/coordinates"
	"github.com/friedenberg/zit/src/alfa/errors"
)

const (
	FilePathKennungYin  = "Kennung/Yin"
	FilePathKennungYang = "Kennung/Yang"
)

type Hinweisen struct {
	sync.Locker
	yin     provider
	yang    provider
	counter coordinates.Int
}

func New(basePath string) (f *Hinweisen, err error) {
	providerPathYin := path.Join(basePath, FilePathKennungYin)
	providerPathYang := path.Join(basePath, FilePathKennungYang)

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
