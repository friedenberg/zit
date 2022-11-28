package typ_toml

type Typ struct {
	InlineAkte     bool                   `toml:"inline-akte,omitempty"`
	FileExtension  string                 `toml:"file-extension,omitempty"`
	ExecCommand    *ScriptConfig          `toml:"exec-command,omitempty"`
	Actions        map[string]*Action     `toml:"actions,omitempty"`
	EtikettenRules map[string]EtikettRule `toml:"etiketten-rules,omitempty"`
}

func (a *Typ) Reset(b *Typ) {
	panic("not implemented")
	// if b == nil {
	// } else {
	// }
}

func (a *Typ) Equals(b *Typ) bool {
	if b == nil || a == nil {
		return false
	}

	if a.InlineAkte != b.InlineAkte {
		return false
	}

	if a.FileExtension != b.FileExtension {
		return false
	}

	if a.ExecCommand != b.ExecCommand {
		return false
	}

	if len(a.Actions) != len(b.Actions) {
		return false
	}

	for k, v := range a.Actions {
		v1, ok := b.Actions[k]

		if !ok {
			return false
		}

		if !v.Equals(v1) {
			return false
		}
	}

	if len(a.EtikettenRules) != len(b.EtikettenRules) {
		return false
	}

	for k, v := range a.EtikettenRules {
		v1, ok := b.EtikettenRules[k]

		if !ok {
			return false
		}

		if v != v1 {
			return false
		}
	}

	return true
}

func (ct *Typ) Apply(kt *Typ) {
	ct.InlineAkte = kt.InlineAkte
	ct.FileExtension = kt.FileExtension

	// if kt.Description != "" {
	// 	ct.Description = collections.MakeStringValue(kt.Description)
	// }

	if len(kt.Actions) > 0 {
		ct.Actions = kt.Actions
	}

	if kt.ExecCommand != nil {
		ct.ExecCommand = kt.ExecCommand
	}

	if len(kt.EtikettenRules) > 0 {
		ct.EtikettenRules = kt.EtikettenRules
	}
}

func (ct *Typ) Merge(ct2 *Typ) {
	if ct2.InlineAkte {
		ct.InlineAkte = true
	}

	if ct2.FileExtension != "" {
		ct.FileExtension = ct2.FileExtension
	}

	ct.ExecCommand.Merge(ct2.ExecCommand)

	for k, v := range ct2.EtikettenRules {
		ct.EtikettenRules[k] = v
	}

	for k, v := range ct2.Actions {
		sc, ok := ct.Actions[k]

		if !ok {
			sc = v
		} else {
			sc.Merge(v)
		}

		ct.Actions[k] = sc
	}
}
