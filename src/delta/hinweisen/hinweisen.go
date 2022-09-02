package hinweisen

import (
	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/sha"
	"github.com/friedenberg/zit/src/charlie/age"
	"github.com/friedenberg/zit/src/charlie/hinweis"
)

type Hinweisen interface {
	StoreNew(sha sha.Sha) (h hinweis.Hinweis, err error)
	Flush() error
	Factory() *factory
}

type hinweisen struct {
	basePath string
	factory  *factory
}

func New(age age.Age, basePath string) (s *hinweisen, err error) {
	s = &hinweisen{
		basePath: basePath,
	}

	if s.factory, err = newFactory(basePath); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (hn hinweisen) Factory() *factory {
	return hn.factory
}

func (zs *hinweisen) Flush() (err error) {
	errors.Print()

	if err = zs.factory.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (hn *hinweisen) StoreNew(sha sha.Sha) (h hinweis.Hinweis, err error) {
	if h, err = hn.factory.Make(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
