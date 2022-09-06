package store_working_directory

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/files"
	"github.com/friedenberg/zit/src/bravo/sha"
)

type UntrackedFile struct {
	Path string
	sha.Sha
}

func (ut UntrackedFile) String() string {
	return fmt.Sprintf("[%s %s]", ut.Path, ut.Sha)
}

func MakeUntrackedFile(dir string, p string) (ut UntrackedFile, err error) {
	ut = UntrackedFile{
		Path: path.Join(dir, p),
	}

	hash := sha256.New()

	var f *os.File

	if f, err = files.Open(ut.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	defer files.Close(f)

	if _, err = io.Copy(hash, f); err != nil {
		err = errors.Wrap(err)
		return
	}

	ut.Sha = sha.FromHash(hash)

	if ut.Path, err = filepath.Rel(dir, ut.Path); err != nil {
		err = errors.Wrapf(err, "%s", ut.Path)
		return
	}

	return
}

type CwdFiles struct {
	dir              string
	Zettelen         []string
	Akten            []string
	UnsureAkten      []UntrackedFile
	EmptyDirectories []string
}

func MakeCwdFiles(dir string) (fs CwdFiles, err error) {
	fs = CwdFiles{
		dir:              dir,
		Zettelen:         make([]string, 0),
		Akten:            make([]string, 0),
		UnsureAkten:      make([]UntrackedFile, 0),
		EmptyDirectories: make([]string, 0),
	}

	var dirs []string

	if dirs, err = files.ReadDirNames(dir); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, d := range dirs {
		if strings.HasPrefix(d, ".") {
			continue
		}

		d2 := path.Join(dir, d)

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
	return len(c.Zettelen) + len(c.Akten)
}

func (s Store) GetPossibleZettels() (result CwdFiles, err error) {
	return MakeCwdFiles(s.path)
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

	var ut UntrackedFile

	if ut, err = MakeUntrackedFile(fs.dir, a); err != nil {
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

	//TODO-decision: should there be hinweis validation?

	//TODO-refactor: akten vs zettel file extensions
	if path.Ext(a) == ".md" {
		fs.Zettelen = append(fs.Zettelen, p)
	} else {
		fs.Akten = append(fs.Akten, p)
	}

	return
}
