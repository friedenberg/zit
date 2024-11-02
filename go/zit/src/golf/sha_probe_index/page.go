package sha_probe_index

import (
	"bufio"
	"io"
	"os"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/heap"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
)

type page struct {
	sync.Mutex
	sha.PageId

	f          *os.File
	br         bufio.Reader
	equaler    interfaces.Equaler1[*row]
	added      addedMap
	dirLayout  dir_layout.DirLayout
	searchFunc func(*sha.Sha) (mid int64, err error)
	rowSize    int
}

func (p *page) initialize(
	equaler interfaces.Equaler1[*row],
	s dir_layout.DirLayout,
	pid sha.PageId,
	rowSize int,
) (err error) {
	p.equaler = equaler

	// p.added = make([]*row, 0)
	p.added = make(addedMap)
	// p.added = heap.Make(
	// 	p.equaler,
	// 	rowLessor{},
	// 	rowResetter{},
	// )

	p.dirLayout = s
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

func (e *page) GetIndexPage() pageInterface {
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

	// e.added = append(e.added, r)
	e.added[r.left.GetBytes()] = r

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

		ui.Out().Printf("%s", &current)
	}
}

func (e *page) Flush() (err error) {
	e.Lock()
	defer e.Unlock()

	if len(e.added) == 0 {
		return
	}

	if e.f != nil {
		if err = e.seekAndResetTo(0); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	var ft *os.File

	if ft, err = e.dirLayout.TempLocal.FileTemp(); err != nil {
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

	// e.added.SortStableAndRemoveDuplicates()
	s := e.added.ToSlice()

	h := heap.MakeHeapFromSlice(
		e.equaler,
		rowLessor{},
		rowResetter{},
		s,
	)

	if err = heap.MergeStream(
		h,
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

	clear(e.added)
	// e.added = e.added[:0]

	if err = e.open(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
