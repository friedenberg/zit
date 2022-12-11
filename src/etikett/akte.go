package etikett

type Akte struct {
	AddToNewZettels bool `toml:"add-to-new-zettels"`
	Hide            bool `toml:"hide"`
}

func (ct *Akte) Merge(ct2 *Akte) {
}
