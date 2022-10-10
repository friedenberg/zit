package store_working_directory

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/charlie/hinweis"
)

type CwdZettel struct {
	hinweis.Hinweis
	Zettel, Akte File
}

type CwdFiles struct {
	dir              string
	Zettelen         map[string]CwdZettel
	UnsureAkten      []File
	EmptyDirectories []string
}

func (fs CwdFiles) ZettelFiles() (out []string) {
	out = make([]string, 0, len(fs.Zettelen))

	for _, z := range fs.Zettelen {
		out = append(out, z.Zettel.Path)
	}

	return
}

func makeCwdFiles(dir string) (fs CwdFiles) {
	fs = CwdFiles{
		dir:              dir,
		Zettelen:         make(map[string]CwdZettel, 0),
		UnsureAkten:      make([]File, 0),
		EmptyDirectories: make([]string, 0),
	}

	return
}

func MakeCwdFilesAll(dir string) (fs CwdFiles, err error) {
	fs = makeCwdFiles(dir)
	err = fs.readAll()
	return
}

func MakeCwdFilesExactly(dir string, files ...string) (fs CwdFiles, err error) {
	fs = makeCwdFiles(dir)
	err = fs.readInputFiles(files...)
	return
}

func (fs *CwdFiles) readInputFiles(args ...string) (err error) {
	for _, f := range args {
		f = filepath.Clean(f)

		if filepath.IsAbs(f) {
			if f, err = filepath.Rel(fs.dir, f); err != nil {
				err = errors.Wrap(err)
				return
			}
		}

		parts := strings.Split(f, string(filepath.Separator))

		switch len(parts) {
		case 0:

		case 1:
			if err = fs.readFirstLevelFile(parts[0]); err != nil {
				err = errors.Wrap(err)
				return
			}

		case 2:
			p := path.Join(parts[len(parts)-2], parts[len(parts)-1])

			if err = fs.readSecondLevelFile(fs.dir, p); err != nil {
				err = errors.Wrap(err)
				return
			}

		default:
			h := path.Join(parts[:len(parts)-3]...)
			p := path.Join(parts[len(parts)-2], parts[len(parts)-1])

			if err = fs.readSecondLevelFile(h, p); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func (fs *CwdFiles) readAll() (err error) {
	var dirs []string

	if dirs, err = files.ReadDirNames(fs.dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, d := range dirs {
		if strings.HasPrefix(d, ".") {
			continue
		}

		d2 := path.Join(fs.dir, d)

		var fi os.FileInfo

		if fi, err = os.Stat(d); err != nil {
			err = errors.Wrap(err)
			return
		}

		if fi.Mode().IsDir() {
			var dirs2 []string

			if dirs2, err = files.ReadDirNames(d2); err != nil {
				err = errors.Wrap(err)
				return
			}

			if len(dirs2) == 0 {
				fs.EmptyDirectories = append(fs.EmptyDirectories, d2)
			}

			for _, a := range dirs2 {
				if err = fs.readSecondLevelFile(d2, a); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

		} else if fi.Mode().IsRegular() {
			if err = fs.readFirstLevelFile(d); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func (c CwdFiles) Len() int {
	return len(c.Zettelen)
}

func (fs *CwdFiles) readFirstLevelFile(a string) (err error) {
	if strings.HasPrefix(a, ".") {
		return
	}

	var fi os.FileInfo

	if fi, err = os.Stat(path.Join(fs.dir, a)); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !fi.Mode().IsRegular() {
		return
	}

	var ut File

	if ut, err = MakeFile(fs.dir, a); err != nil {
		err = errors.Wrap(err)
		return
	}

	fs.UnsureAkten = append(fs.UnsureAkten, ut)

	return
}

func (fs *CwdFiles) readSecondLevelFile(d string, a string) (err error) {
	if strings.HasPrefix(a, ".") {
		return
	}

	var fi os.FileInfo

	p := path.Join(d, a)

	if fi, err = os.Stat(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	if !fi.Mode().IsRegular() {
		return
	}

	var h hinweis.Hinweis

	kopf := filepath.Base(d)

	if h, err = fs.hinweisFromPath(path.Join(kopf, a)); err != nil {
		err = errors.Wrap(err)
		return
	}

	var zcw CwdZettel
	ok := false

	if zcw, ok = fs.Zettelen[h.String()]; !ok {
		zcw.Hinweis = h
	}

	//TODO-refactor: akten vs zettel file extensions
	//TODO read zettels
	if path.Ext(a) == ".md" {
		zcw.Zettel.Path = p
	} else {
		zcw.Akte.Path = p
	}

	fs.Zettelen[h.String()] = zcw

	return
}

func (c CwdFiles) hinweisFromPath(p string) (h hinweis.Hinweis, err error) {
	parts := strings.Split(p, string(filepath.Separator))

	switch len(parts) {
	case 0:
		fallthrough

	case 1:
		err = errors.Errorf("not enough parts: %q", parts)
		return

	default:
		parts = parts[len(parts)-2:]

	case 2:
		break
	}

	p = strings.Join(parts, string(filepath.Separator))

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
