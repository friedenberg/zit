package test_object_metadata_io

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_layout"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
)

func Make(
	t *test_logz.T,
	contents map[string]string,
) (f repo_layout.Layout) {
	t = t.Skip(1)

	p := t.TempDir()

	var dirLayout dir_layout.Layout

	var err error

	if dirLayout, err = dir_layout.MakePrimitiveWithHome(
		p,
		debug.Options{},
	); err != nil {
		t.Fatalf("failed to make dir_layout.Primitive: %s", err)
	}

	if f, err = repo_layout.Make(
		env.MakeDefault(dirLayout),
		repo_layout.Options{
			BasePath:             p,
			PermitNoZitDirectory: true,
		},
	); err != nil {
		t.Fatalf("failed to make dir_layout: %s", err)
	}

	if contents == nil {
		return
	}

	for k, e := range contents {
		var w sha.WriteCloser

		w, err := f.BlobWriter()
		if err != nil {
			t.Fatalf("failed to make blob writer: %s", err)
		}

		_, err = io.Copy(w, strings.NewReader(e))
		if err != nil {
			t.Fatalf("failed to write string to blob writer: %s", err)
		}

		err = w.Close()
		if err != nil {
			t.Fatalf("failed to write string to blob writer: %s", err)
		}

		sh := w.GetShaLike()
		expected := sha.Must(k)

		err = expected.AssertEqualsShaLike(sh)
		if err != nil {
			t.Fatalf("sha mismatch: %s. %s, %q", err, k, e)
		}
	}

	return
}
