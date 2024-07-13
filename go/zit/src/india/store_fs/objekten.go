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

func (c *Store) tryEtikett(fi os.FileInfo, dir string) (err error) {
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

	t, ok := c.etiketten.Get(h.String())

	if !ok {
		t = &KennungFDPair{}
	}

	if err = t.Kennung.SetWithIdLike(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.FDs.Objekte.ResetWith(f)

	return c.etiketten.Add(t)
}

func (c *Store) tryKasten(fi os.FileInfo, dir string) (err error) {
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

	t, ok := c.kisten.Get(h.String())

	if !ok {
		t = &KennungFDPair{}
	}

	if err = t.Kennung.SetWithIdLike(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.FDs.Objekte.ResetWith(f)

	return c.kisten.Add(t)
}

func (c *Store) tryTyp(fi os.FileInfo, dir string) (err error) {
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

	t, ok := c.typen.Get(h.String())

	if !ok {
		t = &KennungFDPair{}
	}

	if err = t.Kennung.SetWithIdLike(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	t.FDs.Objekte.ResetWith(f)

	return c.typen.Add(t)
}

func getHinweis(f *fd.FD, allowErrors bool) (h ids.ZettelId, err error) {
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

	if h, err = getHinweis(f, false); err != nil {
		err = errors.Wrap(err)
		return
	}

	t, ok := c.zettelen.Get(h.String())

	if !ok {
		t = &KennungFDPair{}
	}

	if err = t.Kennung.SetWithIdLike(h); err != nil {
		err = errors.Wrap(err)
		return
	}

	ext := strings.TrimPrefix(path.Ext(name), ".")

	if ext == c.fileExtensions.Zettel {
		if err = t.FDs.Objekte.SetPath(fullPath); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = t.FDs.Akte.SetPath(fullPath); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = c.zettelen.Add(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (c *Store) tryZettelUnsure(
	name string,
	fullPath string,
) (err error) {
	t, ok := c.unsureZettelen.Get(fullPath)

	if !ok {
		t = &KennungFDPair{}
	}

	if err = t.Kennung.SetRaw(name); err != nil {
		err = errors.Wrap(err)
		return
	}

	ext := strings.TrimPrefix(path.Ext(name), ".")

	if ext == c.fileExtensions.Zettel {
		if err = t.FDs.Objekte.SetPath(fullPath); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if err = t.FDs.Akte.SetPath(fullPath); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = c.unsureZettelen.Add(t); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
