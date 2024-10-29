package test_object_metadata_io

import (
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/bravo/test_logz"
	"code.linenisgreat.com/zit/go/zit/src/delta/debug"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/fs_home"
)

func Make(
	t *test_logz.T,
	contents map[string]string,
) (f fs_home.Home) {
	t = t.Skip(1)

	p := t.TempDir()

	var primitive fs_home.Primitive

	var err error

	if primitive, err = fs_home.MakePrimitiveWithHome(p, debug.Options{}); err != nil {
		t.Fatalf("failed to make fs_home.Primitive: %s", err)
	}

	if f, err = fs_home.Make(
		fs_home.Options{
			BasePath:             p,
			PermitNoZitDirectory: true,
		},
		primitive,
	); err != nil {
		t.Fatalf("failed to make fs_home: %s", err)
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
