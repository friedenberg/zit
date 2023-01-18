package organize_text

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/echo/format"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type Metadatei struct {
	kennung.EtikettSet
	Typ kennung.Typ
}

func (m Metadatei) HasMetadateiContent() bool {
	if m.EtikettSet.Len() > 0 {
		return true
	}

	tString := m.Typ.String()

	if tString != "" {
		return true
	}

	return false
}

func (m *Metadatei) ReadFrom(r1 io.Reader) (n int64, err error) {
	r := bufio.NewReader(r1)

	mes := kennung.MakeEtikettMutableSet()

	if n, err = format.ReadLines(
		r,
		format.MakeLineReaderRepeat(
			format.MakeLineReaderKeyValues(
				map[string]format.FuncReadLine{
					"%": format.MakeLineReaderNop(),
					"-": mes.StringAdder(),
					"!": m.Typ.Set,
				},
			),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	m.EtikettSet = mes.Copy()

	return
}

func (m Metadatei) WriteTo(w1 io.Writer) (n int64, err error) {
	w := format.NewLineWriter()

	for _, e := range m.EtikettSet.SortedString() {
		w.WriteFormat("- %s", e)
	}

	tString := m.Typ.String()

	if tString != "" {
		w.WriteFormat("! %s", tString)
	}

	return w.WriteTo(w1)
}
