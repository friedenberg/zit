package compare_map

import (
	"fmt"
	"strings"

	"code.linenisgreat.com/zit/src/echo/bezeichnung"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/foxtrot/metadatei"
)

type SetKeyToMetadatei map[string]*metadatei.Metadatei

func (m SetKeyToMetadatei) String() string {
	sb := &strings.Builder{}

	for h, es := range m {
		fmt.Fprintf(sb, "%s: %s\n", h, es)
	}

	return sb.String()
}

func (s SetKeyToMetadatei) Add(h string, b bezeichnung.Bezeichnung) {
	var m *metadatei.Metadatei
	ok := false

	if m, ok = s[h]; !ok {
		m = &metadatei.Metadatei{}
		metadatei.Resetter.Reset(m)
		m.Bezeichnung = b
	}

	s[h] = m
}

func (s SetKeyToMetadatei) AddEtikett(
	h string,
	e kennung.Etikett,
	b bezeichnung.Bezeichnung,
) {
	var m *metadatei.Metadatei
	ok := false

	if m, ok = s[h]; !ok {
		metadatei.Resetter.Reset(m)
		m.Bezeichnung = b
	}

	if !bezeichnung.Equaler.Equals(m.Bezeichnung, b) {
		panic(fmt.Sprintf("bezeichnung changes: %q != %q", m.Bezeichnung, b))
	}

	kennung.AddNormalizedEtikett(m.GetEtikettenMutable(), &e)

	s[h] = m
}

func (s SetKeyToMetadatei) ContainsEtikett(
	h string,
	e kennung.Etikett,
) (ok bool) {
	var m *metadatei.Metadatei

	if m, ok = s[h]; !ok {
		return
	}

	ok = m.GetEtiketten().Contains(e)

	return
}

type CompareMap struct {
	Named   SetKeyToMetadatei // etikett to hinweis
	Unnamed SetKeyToMetadatei // etikett to bezeichnung
}
