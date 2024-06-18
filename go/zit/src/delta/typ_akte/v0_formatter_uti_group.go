package typ_akte

import "code.linenisgreat.com/zit/go/zit/src/bravo/equality"

type FormatterUTIGroup map[string]string

func (a FormatterUTIGroup) Map() map[string]string {
	return map[string]string(a)
}

func (a *FormatterUTIGroup) Equals(b FormatterUTIGroup) bool {
	if b == nil {
		return false
	}

	if len(a.Map()) != len(b.Map()) {
		return false
	}

	if !equality.MapsOrdered(a.Map(), b.Map()) {
		return false
	}

	return true
}

func (ct *FormatterUTIGroup) Merge(ct2 FormatterUTIGroup) {
	for k, v := range ct2.Map() {
		if v != "" {
			ct.Map()[k] = v
		}
	}
}
