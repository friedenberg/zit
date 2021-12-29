package stored_zettel_formats

import "testing"

func makeEtiketten(t *testing.T, vs ...string) (es _EtikettSet) {
	es = _EtikettNewSet()

	for _, v := range vs {
		if err := es.AddString(v); err != nil {
			t.Fatalf("%s", err)
		}
	}

	return
}

func makeExt(t *testing.T, v string) (es _AkteExt) {
	if err := es.Set(v); err != nil {
		t.Fatalf("%s", err)
	}

	return
}
