package akte_ext

func (a AkteExt) Equals(b AkteExt) bool {
	if a.String() != b.String() {
		return false
	}

	return true
}
