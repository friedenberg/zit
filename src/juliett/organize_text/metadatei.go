package organize_text

import (
	"bufio"
	"io"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
	"github.com/friedenberg/zit/src/bravo/line_format"
	"github.com/friedenberg/zit/src/charlie/kennung"
	"github.com/friedenberg/zit/src/echo/typ"
)

type Metadatei struct {
	kennung.Set
	Typ typ.Kennung
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

	mes := kennung.MakeMutableSet()

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
			if err = mes.AddString(tail); err != nil {
				err = errors.Wrap(err)
				return
			}

		case "!":
			if err = m.Typ.Set(tail); err != nil {
				err = errors.Wrap(err)
				return
			}
		}
	}

	m.Set = mes.Copy()

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
