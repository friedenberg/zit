package ennui

import (
	"bytes"
	"io"
	"os"
	"sync"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/collections"
	"github.com/friedenberg/zit/src/charlie/sha"
	"github.com/friedenberg/zit/src/delta/heap"
	"github.com/friedenberg/zit/src/delta/standort"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
)

// Sha Sha
const RowSize = sha.ByteSize + sha.ByteSize

type (
	// Formatter interface {
	// 	Format(
	// 		io.Writer,
	// 		metadatei.PersistentFormatterContext,
	// 		Options,
	// 	) (int64, error)
	// }

	// Parser interface {
	// 	ParsePersistentMetadatei(
	// 		*catgut.RingBuffer,
	// 		metadatei.PersistentParserContext,
	// 		Options,
	// 	) (int64, error)
	// }

	// Format interface {
	// 	Formatter
	// 	Parser
	// }
	Sha = sha.Sha

	Ennui interface {
		GetEnnui() Ennui
		Add(*metadatei.Metadatei, *sha.Sha) error
		Read(*metadatei.Metadatei, *heap.Heap[Sha, *Sha]) error
		errors.Flusher
	}
)

type Metadatei = metadatei.Metadatei

type ennui struct {
	sync.Mutex
	f        *os.File
	added    *heap.Heap[row, *row]
	standort standort.Standort
}

func Make(s standort.Standort) (e *ennui, err error) {
	e = &ennui{
		added: heap.Make(
			rowEqualer{},
			rowLessor{},
			rowResetter{},
		),
		standort: s,
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
		e.standort.FileVerzeichnisseEnnui(),
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

func (e *ennui) Add(m *Metadatei, sh *sha.Sha) (err error) {
	e.Lock()
	defer e.Unlock()

	var shas []*sha.Sha

	if shas, err = e.getShasForMetadatei(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, s := range shas {
		r := &row{}
		r[0].SetShaLike(s)
		r[1].SetShaLike(sh)
		e.added.Push(r)
	}

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

func (e *ennui) Read(m *Metadatei, h *heap.Heap[Sha, *Sha]) (err error) {
	var shas []*sha.Sha

	if shas, err = e.getShasForMetadatei(m); err != nil {
		err = errors.Wrap(err)
		return
	}

	e.Lock()
	defer e.Unlock()

	for _, s := range shas {
		if err = e.seekToFirst(s); err != nil {
			err = errors.Wrap(err)
			return
		}

		if err = e.collectShas(s, h); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (e *ennui) seekToFirst(shMet *sha.Sha) (err error) {
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

	for low <= hi {
		mid = (hi + low) / 2

		if _, err = e.f.Seek(mid*RowSize, io.SeekStart); err != nil {
			err = errors.Wrap(err)
			return
		}

		if _, err = shMid.ReadFrom(e.f); err != nil {
			err = errors.Wrap(err)
			return
		}

		cmp := bytes.Compare(shMet.GetShaBytes(), shMid.GetShaBytes())

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

	err = collections.ErrNotFound(shMet.String())

	return
}

func (e *ennui) collectShas(
	shMet *sha.Sha,
	h *heap.Heap[Sha, *Sha],
) (err error) {
	for {
		shMid := &sha.Sha{}

		if _, err = shMid.ReadFrom(e.f); err != nil {
			if err == io.EOF {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		if !bytes.Equal(shMet.GetShaBytes(), shMid.GetShaBytes()) {
			return
		}

		if _, err = shMid.ReadFrom(e.f); err != nil {
			err = errors.Wrap(err)
			return
		}

		var sh sha.Sha
		errors.PanicIfError(sh.SetShaLike(shMid))
		h.Push(&sh)
	}
}

func (e *ennui) Flush() (err error) {
	e.Lock()
	defer e.Unlock()

	if e.added.Len() == 0 {
		return
	}

	var ft *os.File

	if ft, err = e.standort.FileTempLocal(); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer errors.DeferredCloser(&err, ft)

	var current row

	getOne := func() (r *row, err error) {
		if e.f == nil {
			err = io.EOF
			return
		}

		_, err = current.ReadFrom(e.f)
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
		func(tz *row) (err error) {
			_, err = tz.WriteTo(ft)
			return
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = os.Rename(
		ft.Name(),
		e.standort.FileVerzeichnisseEnnui(),
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
