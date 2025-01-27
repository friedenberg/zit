package files

import (
	"iter"
	"os"
	"path"
	"path/filepath"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
)

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

func DirNameWriterIgnoringHidden(
	seq iter.Seq2[string, error],
) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		for path, err := range seq {
			if err != nil {
				yield(path, err)
				return
			}

			b := filepath.Base(path)

			if strings.HasPrefix(b, ".") {
				return
			}

			if !yield(path, err) {
				return
			}
		}
	}
}

func DirNamesLevel2(
	p string,
) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		var topLevel quiter.Slice[string]

		{
			var err error

			if topLevel, err = DirNames(p); err != nil {
				yield("", err)
				return
			}
		}

		for topLevelDir := range topLevel.All() {
			var secondLevel quiter.Slice[string]
			{
				var err error

				if secondLevel, err = DirNames(topLevelDir); err != nil {
					yield("", err)
					return
				}
			}

			for secondLevelDir := range secondLevel.All() {
				if !yield(secondLevelDir, nil) {
					return
				}
			}
		}
	}
}
