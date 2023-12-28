package ennui

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/log"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/heap"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/golf/objekte_format"
)

type (
	Sha = sha.Sha

	Ennui interface {
		GetEnnui() Ennui
		AddMetadatei(*metadatei.Metadatei, uint8, uint64) error
		AddSha(*sha.Sha, uint8, uint64) error
		ReadOneSha(sh *Sha) (loc Loc, err error)
		ReadOneKey(string, *metadatei.Metadatei) (Loc, error)
		ReadManyKeys(string, *metadatei.Metadatei, *[]Loc) error
		ReadAll(*metadatei.Metadatei, *[]Loc) error
		PrintAll() error
		errors.Flusher
	}
)

type Metadatei = metadatei.Metadatei

type ennui struct {
	sync.Mutex
	f        *os.File
	br       bufio.Reader
	added    *heap.Heap[row, *row]
	standort standort.Standort
	dir      string
}

func Make(s standort.Standort, dir string) (e *ennui, err error) {
	e = &ennui{
		added: heap.Make(
			rowEqualer{},
			rowLessor{},
			rowResetter{},
		),
		standort: s,
		dir:      dir,
	}

	if err = e.open(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e *ennui) open() (err error) {
	if e.f != nil {
		if err = e.f.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if e.f, err = files.OpenFile(
		path.Join(e.dir, "Ennui"),
		os.O_RDONLY,
		0o666,
	); err != nil {
		if errors.IsNotExist(err) {
			err = nil
		} else {
			err = errors.Wrap(err)
		}

		return
	}

	return
}

func (e *ennui) GetEnnui() Ennui {
	return e
}

func (e *ennui) AddMetadatei(m *Metadatei, page uint8, offset uint64) (err error) {
	e.Lock()
	defer e.Unlock()

	var shas map[string]*sha.Sha

	if shas, err = objekte_format.GetShasForMetadatei(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, s := range shas {
		if err = e.addSha(s, page, offset); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (e *ennui) AddSha(sh *Sha, page uint8, offset uint64) (err error) {
	if sh.IsNull() {
		return
	}

	e.Lock()
	defer e.Unlock()

	return e.addSha(sh, page, offset)
}

func (e *ennui) addSha(sh *Sha, page uint8, offset uint64) (err error) {
	if sh.IsNull() {
		return
	}

	r := &row{}

	if err = r.sha.SetShaLike(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	r.page[0] = page
	r.setOffset(offset)
	e.added.Push(r)

	return
}

func (e *ennui) GetRowCount() (n int64, err error) {
	var fi os.FileInfo

	if fi, err = e.f.Stat(); err != nil {
		err = errors.Wrap(err)
		return
	}

	n = fi.Size()/RowSize - 1

	return
}

func (e *ennui) ReadOneSha(sh *Sha) (loc Loc, err error) {
	e.Lock()
	defer e.Unlock()

	if err = e.seekToFirstBinarySearch(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	if loc, err = e.readCurrentLoc(sh, e.f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e *ennui) ReadOneKey(kf string, m *metadatei.Metadatei) (loc Loc, err error) {
	var f objekte_format.FormatGeneric

	if f, err = objekte_format.FormatForKeyError(kf); err != nil {
		err = errors.Wrap(err)
		return
	}

	var sh *Sha

	if sh, err = objekte_format.GetShaForMetadatei(f, m); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer sha.GetPool().Put(sh)

	if loc, err = e.ReadOneSha(sh); err != nil {
		err = errors.Wrapf(err, "Key: %s", kf)
		return
	}

	return
}

func (e *ennui) ReadManyKeys(
	kf string,
	m *metadatei.Metadatei,
	h *[]Loc,
) (err error) {
	var f objekte_format.FormatGeneric

	if f, err = objekte_format.FormatForKeyError(kf); err != nil {
		err = errors.Wrap(err)
		return
	}

	var sh *Sha

	if sh, err = objekte_format.GetShaForMetadatei(f, m); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.Lock()
	defer e.Unlock()

	if err = e.seekToFirstBinarySearch(sh); err != nil {
		err = errors.Wrapf(err, "Key: %s", kf)
		return
	}

	if err = e.collectLocs(sh, h); err != nil {
		err = errors.Wrapf(err, "Key: %s", kf)
		return
	}

	return
}

func (e *ennui) ReadAll(m *metadatei.Metadatei, h *[]Loc) (err error) {
	var shas map[string]*sha.Sha

	if shas, err = objekte_format.GetShasForMetadatei(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.Lock()
	defer e.Unlock()

	me := errors.MakeMulti()

	for k, s := range shas {
		if err = e.seekToFirstBinarySearch(s); err != nil {
			me.Add(errors.Wrapf(err, "Key: %s", k))
			err = nil
			continue
		}

		if err = e.collectLocs(s, h); err != nil {
			me.Add(errors.Wrapf(err, "Key: %s", k))
			err = nil
			continue
		}
	}

	if me.Len() > 0 {
		err = me
	}

	return
}

func (e *ennui) seekToFirstLinearSearch(shMet *sha.Sha) (err error) {
	if e.f == nil {
		err = collections.ErrNotFound("fd nil: " + shMet.String())
		return
	}

	var rowCount int64
	shMid := &sha.Sha{}

	if rowCount, err = e.GetRowCount(); err != nil {
		err = errors.Wrap(err)
		return
	}

	for mid := int64(0); mid < rowCount; mid++ {
		// var loc int64

		if _, err = e.f.Seek(mid*RowSize, io.SeekStart); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = shMid.ReadFrom(e.f); err != nil {
			err = errors.Wrap(err)
			return
		}

		if bytes.Equal(shMet.GetShaBytes(), shMid.GetShaBytes()) {
			// log.Debug().Printf("%d", loc)

			if _, err = e.f.Seek(mid*RowSize, io.SeekStart); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}
	}

	err = collections.ErrNotFound(shMet.String())

	return
}

func (e *ennui) seekToFirstBinarySearch(shMet *sha.Sha) (err error) {
	if e.f == nil {
		err = collections.ErrNotFound("fd nil: " + shMet.String())
		return
	}

	var low, mid, hi int64
	shMid := &sha.Sha{}

	var rowCount int64

	if rowCount, err = e.GetRowCount(); err != nil {
		err = errors.Wrap(err)
		return
	}

	hi = rowCount - 1
	loops := 0

	for low <= hi {
		loops++
		mid = (hi + low) / 2

		// var loc int64

		if _, err = e.f.Seek(mid*RowSize, io.SeekStart); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = shMid.ReadFrom(e.f); err != nil {
			err = errors.Wrap(err)
			return
		}

		cmp := bytes.Compare(shMet.GetShaBytes(), shMid.GetShaBytes())
		// log.Debug().Printf("%s", shMid)
		// log.Debug().Printf(
		// 	"Lo: %d, Mid: %d, Hi: %d, Loc: %d, Max: %d, cmp: %d",
		// 	low,
		// 	mid,
		// 	hi,
		// 	loc,
		// 	rowCount,
		// 	cmp,
		// )

		switch cmp {
		case -1:
			if low == hi-1 {
				low = hi
			} else {
				hi = mid - 1
			}

		case 0:
			// found
			if _, err = e.f.Seek(mid*RowSize, io.SeekStart); err != nil {
				err = errors.Wrap(err)
				return
			}

			return

		case 1:
			low = mid + 1

		default:
			panic("not possible")
		}
	}

	err = collections.ErrNotFound(fmt.Sprintf("%d: %s", loops, shMet.String()))

	return
}

func (e *ennui) readCurrentLoc(
	in *sha.Sha,
	r io.Reader,
) (out Loc, err error) {
	if in.IsNull() {
		err = errors.Errorf("empty sha")
		return
	}

	sh := sha.GetPool().Get()
	defer sha.GetPool().Put(sh)

	if _, err = sh.ReadFrom(r); err != nil {
		if err != io.EOF {
			err = errors.Wrap(err)
		}

		return
	}

	if !in.Equals(sh) {
		err = io.EOF
		return
	}

	if _, err = out.ReadFrom(e.f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e *ennui) collectLocs(
	shMet *sha.Sha,
	h *[]Loc,
) (err error) {
	e.br.Reset(e.f)

	for {
		var loc Loc

		loc, err = e.readCurrentLoc(shMet, &e.br)

		if err != nil {
			if err == io.EOF {
				err = nil
			}

			return
		}

		*h = append(*h, loc)
	}
}

func (e *ennui) PrintAll() (err error) {
	e.Lock()
	defer e.Unlock()

	if e.f == nil {
		return
	}

	if _, err = e.f.Seek(0, io.SeekStart); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.br.Reset(e.f)

	for {
		var current row

		if _, err = current.ReadFrom(&e.br); err != nil {
			err = errors.WrapExceptAsNil(err, io.EOF)
			return
		}

		log.Out().Printf("%s", &current)
	}
}

func (e *ennui) Flush() (err error) {
	e.Lock()
	defer e.Unlock()

	if e.added.Len() == 0 {
		return
	}

	if e.f != nil {
		if _, err = e.f.Seek(0, io.SeekStart); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	var ft *os.File

	if ft, err = e.standort.FileTempLocal(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ft)

	w := bufio.NewWriter(ft)

	defer errors.DeferredFlusher(&err, w)

	var current row

	e.br.Reset(e.f)

	getOne := func() (r *row, err error) {
		if e.f == nil {
			err = io.EOF
			return
		}

		_, err = current.ReadFrom(&e.br)
		r = &current

		return
	}

	if err = e.added.MergeStream(
		func() (tz *row, err error) {
			tz, err = getOne()

			if errors.IsEOF(err) || tz == nil {
				err = collections.MakeErrStopIteration()
			}

			return
		},
		func(r *row) (err error) {
			_, err = r.WriteTo(w)
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = os.Rename(
		ft.Name(),
		path.Join(e.dir, "Ennui"),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.added.Reset()

	if err = e.open(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
