package object_probe_index

import (
	"bufio"
	"io"
	"os"
	"sync"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/files"
	"code.linenisgreat.com/zit/go/zit/src/delta/heap"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
)

type page struct {
	sync.Mutex // for the buffered reader
	f          *os.File
	br         bufio.Reader
	added      *heap.Heap[row, *row]
	dirLayout  repo_layout.Layout
	searchFunc func(*sha.Sha) (mid int64, err error)
	sha.PageId
}

func (p *page) initialize(
	equaler interfaces.Equaler1[*row],
	s repo_layout.Layout,
	pid sha.PageId,
) (err error) {
	p.added = heap.Make(
		equaler,
		rowLessor{},
		rowResetter{},
	)

	p.dirLayout = s
	p.PageId = pid

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

func (e *page) GetObjectProbeIndexPage() pageInterface {
	return e
}

func (e *page) AddSha(sh *Sha, loc Loc) (err error) {
	if sh.IsNull() {
		return
	}

	e.Lock()
	defer e.Unlock()

	return e.addSha(sh, loc)
}

func (e *page) addSha(sh *Sha, loc Loc) (err error) {
	if sh.IsNull() {
		return
	}

	r := &row{
		Loc: loc,
	}

	if err = r.sha.SetShaLike(sh); err != nil {
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

	n = fi.Size()/RowSize - 1

	return
}

func (e *page) ReadOne(sh *Sha) (loc Loc, err error) {
	e.Lock()
	defer e.Unlock()

	var start int64

	if start, err = e.searchFunc(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.seekAndResetTo(start); err != nil {
		err = errors.Wrap(err)
		return
	}

	if loc, _, err = e.readCurrentLoc(sh, &e.br); err != nil {
		err = errors.Wrapf(err, "Start: %d", start)
		return
	}

	return
}

func (e *page) ReadMany(sh *Sha, locs *[]Loc) (err error) {
	e.Lock()
	defer e.Unlock()

	var start int64

	if start, err = e.searchFunc(sh); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = e.seekAndResetTo(start); err != nil {
		err = errors.Wrap(err)
		return
	}

	isEOF := false

	for !isEOF {
		var loc Loc
		var found bool

		loc, found, err = e.readCurrentLoc(sh, &e.br)

		if err == io.EOF {
			err = nil
			isEOF = true
		} else if err != nil {
			err = errors.Wrap(err)
			return
		}

		if found {
			*locs = append(*locs, loc)
		}
	}

	return
}

func (e *page) readCurrentLoc(
	in *sha.Sha,
	r io.Reader,
) (out Loc, found bool, err error) {
	if in.IsNull() {
		err = errors.Errorf("empty sha")
		return
	}

	sh := sha.GetPool().Get()
	defer sha.GetPool().Put(sh)

	if _, err = sh.ReadFrom(r); err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	if !in.Equals(sh) {
		err = io.EOF
		return
	}

	var n int64
	n, err = out.ReadFrom(r)

	if n > 0 {
		found = true
	}

	if err != nil {
		err = errors.WrapExcept(err, io.EOF)
		return
	}

	return
}

func (e *page) seekAndResetTo(loc int64) (err error) {
	if _, err = e.f.Seek(loc*RowSize, io.SeekStart); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.br.Reset(e.f)

	return
}

func (e *page) PrintAll(env *env.Env) (err error) {
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

		env.GetUI().Printf("%s", &current)
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

	if err = heap.MergeStream(
		e.added,
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
