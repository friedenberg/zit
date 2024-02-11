package kennung

import (
	"path"
	"path/filepath"
	"strings"

	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/echo/fd"
)

func GetIdLike(f *fd.FD) (il Kennung, err error) {
	var h Hinweis

	if h, err = GetHinweis(f); err == nil {
		il = h
		return
	}

	errors.TodoP1("implement Typ and Etikett")

	err = errors.Errorf("not an id")

	return
}

func AsHinweis(f *fd.FD) (h Hinweis, ok bool) {
	var err error
	h, err = GetHinweis(f)
	ok = err == nil
	return
}

func GetHinweis(f *fd.FD) (h Hinweis, err error) {
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
		err = errors.Wrap(err)
		return
	}

	return
}
