package organize_text

import (
	"bufio"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/alfa/schnittstellen"
	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/delta/ohio"
	"github.com/friedenberg/zit/src/echo/format"
	"github.com/friedenberg/zit/src/echo/kennung"
	"github.com/friedenberg/zit/src/foxtrot/metadatei"
	"github.com/friedenberg/zit/src/india/matcher"
)

type Metadatei struct {
	kennung.EtikettSet
	Matchers schnittstellen.SetLike[matcher.Matcher]
	Typ      kennung.Typ
}

func (m Metadatei) AsMetadatei() (m1 metadatei.Metadatei) {
	m1.Typ = m.Typ
	m1.SetEtiketten(m.EtikettSet)
	return
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
		ohio.MakeLineReaderRepeat(
			ohio.MakeLineReaderKeyValues(
				map[string]schnittstellen.FuncSetString{
					"%": ohio.MakeLineReaderNop(),
					"-": iter.MakeFuncSetString[
						kennung.Etikett,
						*kennung.Etikett,
					](mes),
					"!": m.Typ.Set,
				},
			),
		),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	m.EtikettSet = mes.CloneSetPtrLike()

	return
}

func (m Metadatei) WriteTo(w1 io.Writer) (n int64, err error) {
	w := format.NewLineWriter()

	for _, e := range iter.SortedStrings[kennung.Etikett](m.EtikettSet) {
		w.WriteFormat("- %s", e)
	}

	tString := m.Typ.String()

	if tString != "" {
		w.WriteFormat("! %s", tString)
	}

	if m.Matchers != nil {
		for _, c := range iter.SortedStrings[matcher.Matcher](m.Matchers) {
			w.WriteFormat("%% %s", c)
		}
	}

	return w.WriteTo(w1)
}
