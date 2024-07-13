package object_inventory_format

import (
	"bufio"
	"io"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type v0 struct{}

func (f v0) FormatPersistentMetadatei(
	w1 io.Writer,
	c FormatterContext,
	o Options,
) (n int64, err error) {
	m := c.GetMetadata()
	w := format.NewLineWriter()

	if o.Tai {
		w.WriteFormat("Tai %s", m.Tai)
	}

	w.WriteFormat("%s %s", genres.Blob, &m.Blob)
	w.WriteFormat("%s %s", genres.Type, m.GetTyp())
	w.WriteFormat("Bezeichnung %s", m.Description)

	for _, e := range iter.SortedValues[ids.Tag](m.GetEtiketten()) {
		w.WriteFormat("%s %s", genres.Tag, e)
	}

	if n, err = w.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (f v0) ParsePersistentMetadatei(
	r1 io.Reader,
	c ParserContext,
	_ Options,
) (n int64, err error) {
	m := c.GetMetadata()

	etiketten := ids.MakeTagMutableSet()

	r := bufio.NewReader(r1)

	typLineReader := ohio.MakeLineReaderIgnoreErrors(m.Type.Set)

	esa := iter.MakeFuncSetString[ids.Tag, *ids.Tag](
		etiketten,
	)

	var g genres.Genre

	lr := format.MakeLineReaderConsumeEmpty(
		ohio.MakeLineReaderIterate(
			g.Set,
			ohio.MakeLineReaderKeyValues(
				map[string]interfaces.FuncSetString{
					"Tai":                m.Tai.Set,
					genres.Blob.String(): m.Blob.Set,
					genres.Type.String(): typLineReader,
					"AkteTyp":            typLineReader,
					"Bezeichnung":        m.Description.Set,
					genres.Tag.String():  esa,
				},
			),
		),
	)

	if n, err = lr.ReadFrom(r); err != nil {
		err = errors.Wrap(err)
		return
	}

	m.SetEtiketten(etiketten)

	return
}
