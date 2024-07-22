package store_fs

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/echo/fd"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

// TODO combine everyting into one function

func (c *Store) tryTag(fi os.FileInfo, dir string) (err error) {
	var h ids.Tag
	var f *fd.FD

	if f, err = fd.FileInfo(fi, dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	pathMinusExt := strings.ToLower(f.FileNameSansExt())

	if err = h.Set(pathMinusExt); err != nil {
		err = errors.Wrap(err)
		return
	}

	t, ok := c.tags.Get(h.String())

	if !ok {
		t = &ObjectIdFDPair{}
	}

	if err = t.ObjectId.SetWithIdLike(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.FDs.Object.ResetWith(f)

	return c.tags.Add(t)
}

func (c *Store) tryRepo(fi os.FileInfo, dir string) (err error) {
	var h ids.RepoId
	var f *fd.FD

	if f, err = fd.FileInfo(fi, dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	pathMinusExt := strings.ToLower(f.FileNameSansExt())

	if err = h.Set(pathMinusExt); err != nil {
		err = errors.Wrap(err)
		return
	}

	t, ok := c.repos.Get(h.String())

	if !ok {
		t = &ObjectIdFDPair{}
	}

	if err = t.ObjectId.SetWithIdLike(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.FDs.Object.ResetWith(f)

	return c.repos.Add(t)
}

func (c *Store) tryType(fi os.FileInfo, dir string) (err error) {
	var h ids.Type
	var f *fd.FD

	if f, err = fd.FileInfo(fi, dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	pathMinusExt := strings.ToLower(f.FileNameSansExt())

	if err = h.Set(pathMinusExt); err != nil {
		err = errors.Wrap(err)
		return
	}

	t, ok := c.types.Get(h.String())

	if !ok {
		t = &ObjectIdFDPair{}
	}

	if err = t.ObjectId.SetWithIdLike(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.FDs.Object.ResetWith(f)

	return c.types.Add(t)
}

func getZettelId(f *fd.FD, allowErrors bool) (h ids.ZettelId, err error) {
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

	if err = h.Set(p1); err != nil {
		if allowErrors {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (c *Store) tryZettel(
	dir string,
	name string,
	fullPath string,
) (err error) {
	var f *fd.FD

	if f, err = fd.FDFromPath(fullPath); err != nil {
		err = errors.Wrap(err)
		return
	}

	var h ids.ZettelId

	if h, err = getZettelId(f, false); err != nil {
		err = errors.Wrap(err)
		return
	}

	t, ok := c.zettels.Get(h.String())

	if !ok {
		t = &ObjectIdFDPair{}
	}

	if err = t.ObjectId.SetWithIdLike(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	ext := strings.TrimPrefix(path.Ext(name), ".")

	if ext == c.fileExtensions.Zettel {
		if err = t.FDs.Object.SetPath(fullPath); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = t.FDs.Blob.SetPath(fullPath); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = c.zettels.Add(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Store) tryZettelUnsure(
	name string,
	fullPath string,
) (err error) {
	t, ok := c.unsureZettels.Get(fullPath)

	if !ok {
		t = &ObjectIdFDPair{}
	}

	if err = t.ObjectId.SetRaw(name); err != nil {
		err = errors.Wrap(err)
		return
	}

	ext := strings.TrimPrefix(path.Ext(name), ".")

	if ext == c.fileExtensions.Zettel {
		if err = t.FDs.Object.SetPath(fullPath); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = t.FDs.Blob.SetPath(fullPath); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = c.unsureZettels.Add(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
