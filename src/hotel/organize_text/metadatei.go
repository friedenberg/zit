package organize_text

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/line_format"
	"github.com/friedenberg/zit/src/charlie/etikett"
	"github.com/friedenberg/zit/src/charlie/typ"
)

type Metadatei struct {
	etikett.Set
	typ.Typ
}

func (m Metadatei) HasMetadateiContent() bool {
  if m.Set.Len() > 0 {
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

	m.Set = etikett.MakeSet()

	for {
		var rawLine, line string

		rawLine, err = r.ReadString('\n')
		n += int64(len(rawLine))

		if err != nil && err != io.EOF {
			err = errors.Wrap(err)
			return
		}

		if err == io.EOF {
			err = nil
			break
		}

		line = strings.TrimSuffix(rawLine, "\n")
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")

		if len(parts) != 2 {
			err = errors.Errorf("expected exactly two space-separated strings but got %q", parts)
			return
		}

		prefix := parts[0]
		tail := parts[1]

		switch prefix {
		case "-":
			var e etikett.Etikett

			if err = e.Set(tail); err != nil {
				err = errors.Wrap(err)
				return
			}

			m.Set.Add(e)

		case "!":
			if err = m.Typ.Set(tail); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	return
}

func (m Metadatei) WriteTo(w1 io.Writer) (n int64, err error) {
	w := line_format.NewWriter()

	for _, e := range m.Set.SortedString() {
		w.WriteFormat("- %s", e)
	}

	tString := m.Typ.String()

	if tString != "" {
		w.WriteFormat("! %s", tString)
	}

	return w.WriteTo(w1)
}
