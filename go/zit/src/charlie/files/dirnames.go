package files

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

func ReadDir(ps ...string) (dirEntries []os.DirEntry, err error) {
	if dirEntries, err = os.ReadDir(path.Join(ps...)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func ReadDirNamesTo(
	wf func(string) error,
	p string,
) (err error) {
	var names []os.DirEntry

	if names, err = ReadDir(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, n := range names {
		if err = wf(path.Join(p, n.Name())); err != nil {
			err = errors.Wrapf(err, "Path: %q, Name: %q", p, n.Name())
			return
		}
	}

	return
}

func MakeDirNameWriterIgnoringHidden(
	wf func(string) error,
) func(string) error {
	return func(p string) (err error) {
		b := filepath.Base(p)
		if strings.HasPrefix(b, ".") {
			return
		}

		return wf(p)
	}
}

func ReadDirNamesLevel2(
	wf func(string) error,
	p string,
) (err error) {
	ui.TodoP3("support ErrStopIteration")
	ui.TodoP2("support concurrency")

	wfLevel2 := func(p2 string) (err error) {
		return wf(p2)
	}

	wfLevel1 := func(p1 string) (err error) {
		if err = ReadDirNamesTo(
			MakeDirNameWriterIgnoringHidden(wfLevel2),
			p1,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		return
	}

	if err = ReadDirNamesTo(
		MakeDirNameWriterIgnoringHidden(wfLevel1),
		p,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
