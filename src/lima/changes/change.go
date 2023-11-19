package changes

import (
	"strings"

	"github.com/friedenberg/zit/src/bravo/iter"
	"github.com/friedenberg/zit/src/bravo/values"
	"github.com/friedenberg/zit/src/echo/kennung"
)

type Change struct {
	Key     string
	added   kennung.EtikettMutableSet
	removed kennung.EtikettMutableSet
}

func (a Change) GetAdded() kennung.EtikettSet {
	return a.added
}

func (a Change) GetRemoved() kennung.EtikettSet {
	return a.removed
}

func (a Change) EqualsAny(b any) bool {
	return values.Equals(a, b)
}

func (a Change) Equals(b Change) bool {
	if a.Key != b.Key {
		return false
	}

	if !iter.SetEqualsPtr[kennung.Etikett, *kennung.Etikett](a.added, b.added) {
		return false
	}

	if !iter.SetEqualsPtr[kennung.Etikett, *kennung.Etikett](
		a.removed,
		b.removed,
	) {
		return false
	}

	return true
}

func (c Change) IsEmpty() bool {
	return c.added.Len() == 0 && c.removed.Len() == 0
}

func (c Change) String() string {
	return c.Key
}

func (c Change) Description() string {
	sb := &strings.Builder{}
	sb.WriteString(c.String())
	sb.WriteString(" add: ")
	sb.WriteString(iter.StringCommaSeparated[kennung.Etikett](c.added))
	sb.WriteString(" remove: ")
	sb.WriteString(iter.StringCommaSeparated[kennung.Etikett](c.removed))
	return sb.String()
}
