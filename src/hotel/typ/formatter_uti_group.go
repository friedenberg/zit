package typ

type FormatterUTIGroup map[string]string

func (a FormatterUTIGroup) Map() map[string]string {
	return map[string]string(a)
}

func (a *FormatterUTIGroup) Equals(b *FormatterUTIGroup) bool {
	if b == nil {
		return false
	}

	if len(a.Map()) != len(b.Map()) {
		return false
	}

	for k, v := range a.Map() {
		if vb, ok := b.Map()[k]; !ok {
			return false
		} else if vb != v {
			return false
		}
	}

	return true
}

func (ct *FormatterUTIGroup) Merge(ct2 *FormatterUTIGroup) {
	for k, v := range ct2.Map() {
		ct.Map()[k] = v
	}
}
