package script_config

type ScriptConfigWithUTI struct {
	ScriptConfig
	UTI string `toml:"uti"`
}

func (a *ScriptConfigWithUTI) Equals(b ScriptConfigWithUTI) bool {
	if a.UTI != b.UTI {
		return false
	}

	return true
}

func (s *ScriptConfigWithUTI) Merge(s2 ScriptConfigWithUTI) {
	if s2.UTI != "" {
		s.UTI = s2.UTI
	}

	s.ScriptConfig.Merge(s2.ScriptConfig)
}
