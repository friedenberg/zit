package verzeichnisse

import (
	"bufio"
	"encoding/gob"
	"io"
	"path"
	"sync"

	"github.com/friedenberg/zit/alfa/logz"
	"github.com/friedenberg/zit/bravo/errors"
	"github.com/friedenberg/zit/bravo/stdprinter"
	"github.com/friedenberg/zit/charlie/open_file_guard"
	age_io "github.com/friedenberg/zit/delta/age_io"
)

type RowMaker func() ([]Row, error)

type Index struct {
	path          string
	idTransformer IdTransformer
	rwLock        *sync.RWMutex
	pages         map[string]*page
	age_io.ReadCloserFactory
	age_io.WriteCloserFactory
}

type page struct {
	rows       []Row
	hasChanges bool
}

func NewIndex(
	path string,
	r age_io.ReadCloserFactory,
	w age_io.WriteCloserFactory,
	idTransformer IdTransformer,
) (i *Index, err error) {
	logz.Print("initing verzeichnisse")
	i = &Index{
		path:               path,
		pages:              make(map[string]*page),
		idTransformer:      idTransformer,
		rwLock:             &sync.RWMutex{},
		ReadCloserFactory:  r,
		WriteCloserFactory: w,
	}

	return
}

func (i *Index) Flush() (err error) {
	i.rwLock.RLock()
	defer i.rwLock.RUnlock()

	//TODO make atomic
	for id, p := range i.pages {
		if err = i.writePage(id, p); err != nil {
			err = errors.Wrapped(err, "failed to flush page: %s: %s", id, err)
			return
		}
	}

	return
}

func (i *Index) ReadPages(r Reader, ids ...string) (err error) {
	logz.Printf("reading pages: %s", ids)
	wg := &sync.WaitGroup{}
	wg.Add(len(ids))

	if err = r.Begin(); err != nil {
		err = errors.Wrapped(err, "closing index reader failed")
		return
	}

	for _, id := range ids {
		logz.PrintDebug(id)
		go func(id string) {
			defer wg.Done()

			var p *page

			if p, err = i.readPage(id); err != nil {
				err = errors.Error(err)
				return
			}

			if err = i.readPageRows(id, p, r); err != nil {
				err = errors.Error(err)
				return
			}
		}(id)
	}

	wg.Wait()

	if err = r.End(); err != nil {
		err = errors.Wrapped(err, "closing index reader failed")
		return
	}

	return
}

func (i *Index) GetAllPageIds() (ids []string, err error) {
	if ids, err = open_file_guard.ReadDirNames(i.path); err != nil {
		err = errors.Wrapped(err, "failed to read all page ids: %s", i.path)
		return
	}

	return
}

func (i *Index) ReadAll(r Reader) (err error) {
	var ids []string

	if ids, err = i.GetAllPageIds(); err != nil {
		err = errors.Error(err)
		return
	}

	return i.ReadPages(r, ids...)
}

func (i *Index) WriteRows(rowMaker RowMaker) (err error) {
	var rs []Row

	if rs, err = rowMaker(); err != nil {
		err = errors.Error(err)
		return
	}

	return i.Write(rs...)
}

func (i *Index) Write(rs ...Row) (err error) {
	for _, r := range rs {
		id := i.idTransformer(r.Sha)

		if err = i.writeRow(id, r); err != nil {
			err = errors.Error(err)
			return
		}
	}

	return
}

func (i *Index) writeRow(id string, row Row) (err error) {
	var p *page

	if p, err = i.readPage(id); err != nil {
		err = errors.Error(err)
		return
	}

	if p == nil {
		err = errors.Errorf("read page returned nil for page: %s", id)
		return
	}

	logz.Print("writing row")
	//TODO should this be deduped?
	p.rows = append(p.rows, row)
	p.hasChanges = true

	return
}

func (i *Index) readPage(id string) (p *page, err error) {
	ok := false

	i.rwLock.RLock()

	if p, ok = i.pages[id]; ok {
		logz.Printf("page was cached: %s", id)
		i.rwLock.RUnlock()
		return
	}

	logz.Printf("page needs to be read: %s", id)

	i.rwLock.RUnlock()

	p = &page{
		rows: make([]Row, 0),
	}

	var r1 io.ReadCloser

	r1, err = i.ReadCloser(path.Join(i.path, id))

	if err == nil {
		defer r1.Close()

		r := bufio.NewReader(r1)

		dec := gob.NewDecoder(r)

		var rows []Row

		err = dec.Decode(&rows)

		if err == nil {
			p.rows = rows
		} else {
			if errors.IsEOF(err) {
				err = nil
			} else {
				err = errors.Wrapped(err, "failed to decode page: %s", id)
				return
			}
		}
	} else {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrapped(err, "failed to create reader for page: %s", id)
			return
		}
	}

	i.rwLock.Lock()
	defer i.rwLock.Unlock()

	i.pages[id] = p

	return
}

//   _               _      ____                  _              _
//  | |    ___   ___| | __ |  _ \ ___  __ _ _   _(_)_ __ ___  __| |
//  | |   / _ \ / __| |/ / | |_) / _ \/ _` | | | | | '__/ _ \/ _` |
//  | |__| (_) | (__|   <  |  _ <  __/ (_| | |_| | | | |  __/ (_| |
//  |_____\___/ \___|_|\_\ |_| \_\___|\__, |\__,_|_|_|  \___|\__,_|
//                                       |_|

func (i *Index) readPageRows(id string, p *page, rr Reader) (err error) {
	if rr == nil {
		return
	}

	for _, row := range p.rows {
		if err = rr.ReadRow(id, row); err != nil {
			err = errors.Wrapped(err, "row reader failed to read row")
			return
		}
	}

	return
}

func (i *Index) writePage(id string, p *page) (err error) {
	if !p.hasChanges {
		return
	}

	var w1 io.WriteCloser

	if w1, err = i.WriteCloser(path.Join(i.path, id)); err != nil {
		err = errors.Wrapped(err, "failed to make write closer for page: %s", id)
		return
	}

	defer stdprinter.PanicIfError(w1.Close)

	w := bufio.NewWriter(w1)

	defer stdprinter.PanicIfError(w.Flush)

	enc := gob.NewEncoder(w)

	if err = enc.Encode(p.rows); err != nil {
		err = errors.Wrapped(err, "failed to write encoded page: %s", id)
		return
	}

	return
}
