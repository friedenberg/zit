package hinweisen

import (
	"io"
	"os"
	"path"
	"strconv"
	"sync"
)

const (
	FilePathKennungYin     = "Kennung/Yin"
	FilePathKennungYang    = "Kennung/Yang"
	FilePathKennungCounter = "Kennung/Counter"
)

type factory struct {
	sync.Locker
	pathLastId string
	yin        provider
	yang       provider
	counter    _Int
}

func newFactory(basePath string) (f *factory, err error) {
	providerPathYin := path.Join(basePath, FilePathKennungYin)
	providerPathYang := path.Join(basePath, FilePathKennungYang)
	idLockPath := path.Join(basePath, FilePathKennungCounter)

	f = &factory{
		Locker:     &sync.Mutex{},
		pathLastId: idLockPath,
	}

	if f.yin, err = newProvider(providerPathYin); err != nil {
		err = _Error(err)
		return
	}

	if f.yang, err = newProvider(providerPathYang); err != nil {
		err = _Error(err)
		return
	}

	var old string

	if old, err = _ReadAllString(f.pathLastId); err != nil {
		return
	}

	if f.counter, err = strconv.ParseUint(old, 10, 64); err != nil {
		return
	}

	return
}

func (hf *factory) Make() (h _Hinweis, err error) {
	hf.Lock()
	defer hf.Unlock()

	newInt := hf.counter + 1

	if h, err = _NewHinweis(newInt, hf.yin, hf.yang); err != nil {
		err = _Error(err)
		return
	}

	hf.counter = newInt

	return
}

func (hf factory) Flush() (err error) {
	hf.Lock()
	defer hf.Unlock()

	var f *os.File

	if f, err = _TempFile(); err != nil {
		return
	}

	defer _Close(f)

	io.WriteString(f, strconv.FormatInt(int64(hf.counter), 10))

	if err = os.Rename(f.Name(), hf.pathLastId); err != nil {
		return
	}

	return
}
