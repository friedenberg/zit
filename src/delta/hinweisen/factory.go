package hinweisen

import (
	"io"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/kennung"
	"github.com/friedenberg/zit/src/bravo/open_file_guard"
	"github.com/friedenberg/zit/src/charlie/hinweis"
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
	counter    kennung.Int
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
		err = errors.Wrap(err)
		return
	}

	if f.yang, err = newProvider(providerPathYang); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = f.Refresh(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (hf *factory) Left() provider {
	return hf.yin
}

func (hf *factory) Right() provider {
	return hf.yang
}

func (hf *factory) Refresh() (err error) {
	hf.Lock()
	defer hf.Unlock()

	err = hf.refresh()

	return
}

func (hf *factory) refresh() (err error) {
	var old string

	if old, err = open_file_guard.ReadAllString(hf.pathLastId); err != nil {
		return
	}

	if hf.counter, err = strconv.ParseUint(old, 10, 64); err != nil {
		return
	}

	return
}

func (hf *factory) Make() (h hinweis.Hinweis, err error) {
	errors.Print("making")
	hf.Lock()
	defer hf.Unlock()
	defer func() {
		if err == nil {
			err = hf.flush()
		}
	}()

	// if err = hf.refresh(); err != nil {
	// 	err = errors.Error(err)
	// 	return
	// }

	newInt := hf.counter + 1

	errors.Printf("next kennung: %d", newInt)
	if h, err = hinweis.New(newInt, hf.yin, hf.yang); err != nil {
		err = errors.Wrap(err)
		return
	}

	hf.counter = newInt

	return
}

func (hf factory) Flush() (err error) {
	hf.Lock()
	defer hf.Unlock()

	return hf.flush()
}

func (hf factory) flush() (err error) {
	var f *os.File

	if f, err = open_file_guard.TempFile(); err != nil {
		return
	}

	defer open_file_guard.Close(f)

	io.WriteString(f, strconv.FormatInt(int64(hf.counter), 10))

	if err = os.Rename(f.Name(), hf.pathLastId); err != nil {
		return
	}

	return
}
