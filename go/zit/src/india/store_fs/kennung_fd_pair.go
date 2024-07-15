package store_fs

import (
	"path"
	"path/filepath"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type KennungFDPair struct {
	Kennung ids.ObjectId
	FDs     FDPair
}

func (a *KennungFDPair) String() string {
	return a.Kennung.String()
}

func (a *KennungFDPair) Equals(b KennungFDPair) bool {
	if a.Kennung.String() != b.Kennung.String() {
		return false
	}

	if !a.FDs.Equals(&b.FDs) {
		return false
	}

	return true
}

func (e *KennungFDPair) GetKennungLike() ids.IdLike {
	return &e.Kennung
}

func (e *KennungFDPair) GetKennungLikePtr() ids.IdLikePtr {
	return &e.Kennung
}

func (e *KennungFDPair) GetFDs() *FDPair {
	return &e.FDs
}

func (e *KennungFDPair) GetObjekteFD() *fd.FD {
	return &e.FDs.Object
}

func (e *KennungFDPair) GetAkteFD() *fd.FD {
	return &e.FDs.Blob
}

func (e *KennungFDPair) SetKennungFromFullPath(
	fullPath string,
	allowErrors bool,
) (err error) {
	var f *fd.FD

	if f, err = fd.FDFromPath(fullPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	parts := strings.Split(f.GetPath(), string(filepath.Separator))

	switch len(parts) {
	case 0, 1:
		err = errors.Errorf("not enough parts: %q", parts)
		return

	default:
		parts = parts[len(parts)-2:]
	}

	p := strings.Join(parts, string(filepath.Separator))

	p1 := p
	ext := path.Ext(p)

	if len(ext) != 0 {
		p1 = p[:len(p)-len(ext)]
	}

	e.Kennung.SetGenre(genres.Zettel)

	if err = e.Kennung.Set(p1); err != nil {
		if allowErrors {
			if err = e.Kennung.SetRaw(p1); err != nil {
				err = errors.Wrap(err)
				return
			}
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
