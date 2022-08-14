package stored_zettel_formats

import (
	"testing"

	"github.com/friedenberg/zit/bravo/akte_ext"
	"github.com/friedenberg/zit/delta/etikett"
)

func makeEtiketten(t *testing.T, vs ...string) (es etikett.Set) {
	es = etikett.MakeSet()

	for _, v := range vs {
		if err := es.AddString(v); err != nil {
			t.Fatalf("%s", err)
		}
	}

	return
}

func makeExt(t *testing.T, v string) (es akte_ext.AkteExt) {
	if err := es.Set(v); err != nil {
		t.Fatalf("%s", err)
	}

	return
}
