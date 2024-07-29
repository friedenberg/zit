package store_fs

import (
	"path"
	"path/filepath"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type ObjectIdFDPair struct {
	ObjectId ids.ObjectId
	FDs      FDPair
}

func (a *ObjectIdFDPair) String() string {
	return a.ObjectId.String()
}

func (a *ObjectIdFDPair) Equals(b ObjectIdFDPair) bool {
	if a.ObjectId.String() != b.ObjectId.String() {
		return false
	}

	if !a.FDs.Equals(&b.FDs) {
		return false
	}

	return true
}

func (e *ObjectIdFDPair) GetObjectId() *ids.ObjectId {
	return &e.ObjectId
}

func (e *ObjectIdFDPair) GetFDs() *FDPair {
	return &e.FDs
}

func (e *ObjectIdFDPair) GetObjectFD() *fd.FD {
	return &e.FDs.Object
}

func (e *ObjectIdFDPair) GetBlobFD() *fd.FD {
	return &e.FDs.Blob
}

func (e *ObjectIdFDPair) SetObjectIdFromFullPath(
	fullPath string,
	allowErrors bool,
) (err error) {
	var f *fd.FD

	if f, err = fd.MakeFromPath(fullPath); err != nil {
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

	e.ObjectId.SetGenre(genres.Zettel)

	if err = e.ObjectId.Set(p1); err != nil {
		if allowErrors {
			if err = e.ObjectId.SetRaw(p1); err != nil {
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

func SetAddPairs(
	in interfaces.SetLike[*ObjectIdFDPair],
	out fd.MutableSet,
) (err error) {
	return in.Each(
		func(e *ObjectIdFDPair) (err error) {
			ofd := e.GetObjectFD()

			if !ofd.IsEmpty() {
				if err = out.Add(ofd); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			ofd = e.GetBlobFD()

			if !ofd.IsEmpty() {
				if err = out.Add(ofd); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
	)
}
