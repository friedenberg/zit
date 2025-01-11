package repo_blobs

type TomlRelayLocalV0 struct {
	Path string `toml:"path"`
}

func (b TomlRelayLocalV0) GetRepoBlob() Blob {
	return b
}

func (a *TomlRelayLocalV0) Reset() {
	a.Path = ""
}

func (a *TomlRelayLocalV0) ResetWith(b TomlRelayLocalV0) {
	a.Path = b.Path
}

func (a TomlRelayLocalV0) Equals(b TomlRelayLocalV0) bool {
	if a.Path != b.Path {
		return false
	}

	return true
}
