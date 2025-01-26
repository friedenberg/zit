package files

import (
	"iter"
	"os"
	"path"
	"path/filepath"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
)

func ReadDir(ps ...string) (dirEntries []os.DirEntry, err error) {
	if dirEntries, err = os.ReadDir(path.Join(ps...)); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func DirNames(p string) (slice quiter.Slice[string], err error) {
	var names []os.DirEntry

	if names, err = ReadDir(p); err != nil {
		err = errors.Wrap(err)
		return
	}

	for _, n := range names {
		slice.Append(path.Join(p, n.Name()))
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

func DirNameWriterIgnoringHidden(
	seq iter.Seq[quiter.ElementOrError[string]],
) iter.Seq[quiter.ElementOrError[string]] {
	return func(yield func(quiter.ElementOrError[string]) bool) {
		for pathOrError := range seq {
			if pathOrError.Error != nil {
				yield(pathOrError)
				return
			}

			b := filepath.Base(pathOrError.Element)

			if strings.HasPrefix(b, ".") {
				return
			}

			if !yield(pathOrError) {
				return
			}
		}
	}
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

func DirNamesLevel2(
	p string,
) iter.Seq[quiter.ElementOrError[string]] {
	return func(yield func(quiter.ElementOrError[string]) bool) {
		var topLevel quiter.Slice[string]

		{
			var err error

			if topLevel, err = DirNames(p); err != nil {
				yield(quiter.ElementOrError[string]{Error: errors.Wrap(err)})
				return
			}
		}

		for topLevelDir := range topLevel.All() {
			var secondLevel quiter.Slice[string]
			{
				var err error

				if secondLevel, err = DirNames(topLevelDir); err != nil {
					yield(quiter.ElementOrError[string]{Error: errors.Wrap(err)})
					return
				}
			}

			for secondLevelDir := range secondLevel.All() {
				if !yield(quiter.ElementOrError[string]{Element: secondLevelDir}) {
					return
				}
			}
		}
	}
}
