package format

import (
	"bufio"
	"io"
	"sort"
	"strings"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type KeyValuePair struct {
	Key   string
	Value string
}

type KeyValueGetter interface {
	KeyValuePairs() []KeyValuePair
}

type KeyValueGetterHeader interface {
	Header() string
}

type KeyValueSetter interface {
	SetKeyValue(k, v string) (err error)
}

type KeyValuer interface {
	KeyValueGetter
	KeyValueSetter
}

type KeyValues struct {
}

func MakeKeyValues() *KeyValues {
	return &KeyValues{}
}

func (f KeyValues) ReadFormat(
	r1 io.Reader,
	o KeyValueSetter,
) (n int64, err error) {
	r := bufio.NewReader(r1)
	header := ""
	tryHeader := false

	if h, ok := o.(KeyValueGetterHeader); ok {
		header = h.Header()
		tryHeader = header != ""
	}

	for {
		var lineOriginal string
		lineOriginal, err = r.ReadString('\n')

		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return
		}

		// line := strings.TrimSpace(lineOriginal)
		line := lineOriginal

		n += int64(len(lineOriginal))

		loc := strings.Index(line, " ")

		if line == "" {
			errors.TodoP4("handle empty lines in a more structured way")
		}

		if tryHeader {
			if line != header {
				err = errors.Errorf(
					"expected header %q but got %q, line number: %d",
					header,
					line,
					loc,
				)

				return
			}

			tryHeader = false
		} else {
			v := line[loc+1:]

			if err = o.SetKeyValue(line[:loc], v); err != nil {
				err = errors.Wrapf(err, "Line Number: %d", loc)
				return
			}
		}
	}

	return
}

func (f KeyValues) WriteFormat(
	w1 io.Writer,
	o KeyValueGetter,
) (n int64, err error) {
	w := NewWriter()

	if h, ok := o.(KeyValueGetterHeader); ok {
		header := h.Header()

		if header != "" {
			w.WriteLines(header)
		}
	}

	keyValuePairs := o.KeyValuePairs()

	sort.Slice(
		keyValuePairs,
		func(i, j int) bool {
			keysSame := keyValuePairs[i].Key == keyValuePairs[j].Key

			if keysSame {
				return keyValuePairs[i].Value < keyValuePairs[j].Value
			} else {
				return keyValuePairs[i].Key < keyValuePairs[j].Key
			}
		},
	)

	for _, kvp := range keyValuePairs {
		w.WriteFormat("%s %s", kvp.Key, kvp.Value)
	}

	if n, err = w.WriteTo(w1); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
