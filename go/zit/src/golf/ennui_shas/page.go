package ennui_shas

import (
	"bufio"
	"io"
	"os"
	"sync"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/bravo/files"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/src/charlie/sha"
	"code.linenisgreat.com/zit/src/delta/heap"
	"code.linenisgreat.com/zit/src/delta/standort"
)

type page struct {
	sync.Mutex
	f          *os.File
	br         bufio.Reader
	added      *heap.Heap[row, *row]
	standort   standort.Standort
	searchFunc func(*sha.Sha) (mid int64, err error)
	rowSize    int
	sha.PageId
}

func (p *page) initialize(
	equaler schnittstellen.Equaler1[*row],
	s standort.Standort,
	pid sha.PageId,
	rowSize int,
) (err error) {
	p.added = heap.Make(
		equaler,
		rowLessor{},
		rowResetter{},
	)

	p.standort = s
	p.PageId = pid

	p.rowSize = rowSize

	p.searchFunc = p.seekToFirstBinarySearch

	if err = p.open(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e *page) open() (err error) {
	if e.f != nil {
		if err = e.f.Close(); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if e.f, err = files.OpenFile(
		e.Path(),
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

func (e *page) GetEnnuiPage() pageInterface {
	return e
}

func (e *page) AddSha(left, right *Sha) (err error) {
	if left.IsNull() {
		return
	}

	e.Lock()
	defer e.Unlock()

	return e.addSha(left, right)
}

func (e *page) addSha(left, right *Sha) (err error) {
	if left.IsNull() {
		return
	}

	r := &row{}

	if err = r.left.SetShaLike(left); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = r.right.SetShaLike(right); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.added.Push(r)

	return
}

func (e *page) GetRowCount() (n int64, err error) {
	var fi os.FileInfo

	if fi, err = e.f.Stat(); err != nil {
		err = errors.Wrap(err)
		return
	}

	n = fi.Size()/int64(e.rowSize) - 1

	return
}

func (e *page) ReadOne(left *Sha) (right *Sha, err error) {
	e.Lock()
	defer e.Unlock()

	var start int64

	if start, err = e.searchFunc(left); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.seekAndResetTo(start); err != nil {
		err = errors.Wrap(err)
		return
	}

	if right, err = e.readCurrentLoc(left, e.f); err != nil {
		err = errors.Wrapf(err, "Start: %d", start)
		return
	}

	return
}

func (e *page) ReadMany(sh *Sha, locs *[]*sha.Sha) (err error) {
	e.Lock()
	defer e.Unlock()

	var start int64

	if start, err = e.searchFunc(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.collectLocs(sh, locs, start); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e *page) readCurrentLoc(
	in *sha.Sha,
	r io.Reader,
) (out *sha.Sha, err error) {
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

	out = sha.GetPool().Get()

	if _, err = out.ReadFrom(e.f); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (e *page) seekAndResetTo(loc int64) (err error) {
	if _, err = e.f.Seek(loc*int64(e.rowSize), io.SeekStart); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.br.Reset(e.f)

	return
}

func (e *page) collectLocs(
	shMet *sha.Sha,
	h *[]*sha.Sha,
	start int64,
) (err error) {
	if err = e.seekAndResetTo(start); err != nil {
		err = errors.Wrap(err)
		return
	}

	for {
		var sh *sha.Sha

		sh, err = e.readCurrentLoc(shMet, &e.br)

		if err != nil {
			if err == io.EOF {
				err = nil
			}

			return
		}

		*h = append(*h, sh)
	}
}

func (e *page) PrintAll() (err error) {
	e.Lock()
	defer e.Unlock()

	if e.f == nil {
		return
	}

	if err = e.seekAndResetTo(0); err != nil {
		err = errors.Wrap(err)
		return
	}

	for {
		var current row

		if _, err = current.ReadFrom(&e.br); err != nil {
			err = errors.WrapExceptAsNil(err, io.EOF)
			return
		}

		log.Out().Printf("%s", &current)
	}
}

func (e *page) Flush() (err error) {
	e.Lock()
	defer e.Unlock()

	if e.added.Len() == 0 {
		return
	}

	if e.f != nil {
		if err = e.seekAndResetTo(0); err != nil {
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

	getOne := func() (r *row, err error) {
		if e.f == nil {
			err = io.EOF
			return
		}

		var n int64
		n, err = current.ReadFrom(&e.br)

		if err != nil {
			if errors.Is(err, io.ErrUnexpectedEOF) && n == 0 {
				err = io.EOF
			}

			err = errors.WrapExcept(err, io.EOF)
			return
		}

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
		e.Path(),
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
