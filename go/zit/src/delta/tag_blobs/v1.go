package tag_blobs

type V1 struct {
	Filter string `toml:"filter"`
}

func (a *V1) Reset() {
	a.Filter = ""
}

func (a *V1) ResetWith(b V1) {
	a.Filter = b.Filter
}
