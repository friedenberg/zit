package object_metadata

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/charlie/script_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/dir_layout"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type textFormatterCommon struct {
	dirLayout     dir_layout.DirLayout
	blobFormatter script_config.RemoteScript
}

func (f textFormatterCommon) writeComments(
	w1 io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	n1 := 0

	for _, c := range c.GetMetadata().Comments {
		n1, err = io.WriteString(w1, "% ")
		n += int64(n1)

		if err != nil {
			return
		}

		n1, err = io.WriteString(w1, c)
		n += int64(n1)

		if err != nil {
			return
		}

		n1, err = io.WriteString(w1, "\n")
		n += int64(n1)

		if err != nil {
			return
		}
	}

	return
}

func (f textFormatterCommon) writeBoundary(
	w1 io.Writer,
	_ TextFormatterContext,
) (n int64, err error) {
	return ohio.WriteLine(w1, Boundary)
}

func (f textFormatterCommon) writeNewLine(
	w1 io.Writer,
	_ TextFormatterContext,
) (n int64, err error) {
	return ohio.WriteLine(w1, "")
}

func (f textFormatterCommon) writeCommonMetadataFormat(
	w1 io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	w := format.NewLineWriter()
	m := c.GetMetadata()

	if m.Description.String() != "" || !c.DoNotWriteEmptyDescription {
		sr := bufio.NewReader(strings.NewReader(m.Description.String()))

		for {
			var line string
			line, err = sr.ReadString('\n')
			isEOF := err == io.EOF

			if err != nil && !isEOF {
				err = errors.Wrap(err)
				return
			}

			w.WriteLines(
				fmt.Sprintf("# %s", strings.TrimSpace(line)),
			)

			if isEOF {
				break
			}
		}
	}

	for _, e := range quiter.SortedValues(m.GetTags()) {
		if ids.IsEmpty(e) {
			continue
		}

		w.WriteFormat("- %s", e)
	}

	if n, err = w.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f textFormatterCommon) writeTyp(
	w1 io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	m := c.GetMetadata()

	if m.Type.IsEmpty() {
		return
	}

	return ohio.WriteLine(w1, fmt.Sprintf("! %s", m.Type))
}

func (f textFormatterCommon) writeShaTyp(
	w1 io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	m := c.GetMetadata()
	return ohio.WriteLine(w1, fmt.Sprintf("! %s.%s", &m.Blob, m.Type))
}

func (f textFormatterCommon) writePathType(
	w1 io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	var ap string

	for _, f := range c.PersistentFormatterContext.GetMetadata().Fields {
		if strings.ToLower(f.Key) == "blob" {
			ap = f.Value
			break
		}
	}

	if ap != "" {
		ap = f.dirLayout.RelToCwdOrSame(ap)
	} else {
		err = errors.Errorf("path not found in fields")
		return
	}

	return ohio.WriteLine(w1, fmt.Sprintf("! %s", ap))
}

func (f textFormatterCommon) writeBlob(
	w1 io.Writer,
	c TextFormatterContext,
) (n int64, err error) {
	var ar io.ReadCloser
	m := c.GetMetadata()

	if ar, err = f.dirLayout.BlobReader(&m.Blob); err != nil {
		err = errors.Wrap(err)
		return
	}

	if ar == nil {
		err = errors.Errorf("blob reader is nil")
		return
	}

	defer errors.DeferredCloser(&err, ar)

	if f.blobFormatter != nil {
		var wt io.WriterTo

		if wt, err = script_config.MakeWriterToWithStdin(
			f.blobFormatter,
			f.dirLayout.MakeCommonEnv(),
			ar,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		if n, err = wt.WriteTo(w1); err != nil {
			err = errors.Wrap(err)
			return
		}
	} else {
		if n, err = io.Copy(w1, ar); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}
