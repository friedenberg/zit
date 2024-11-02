package tag_blobs

type TomlV1 struct {
	Filter string `toml:"filter"`
}

func (a *TomlV1) Reset() {
	a.Filter = ""
}

func (a *TomlV1) ResetWith(b TomlV1) {
	a.Filter = b.Filter
}
