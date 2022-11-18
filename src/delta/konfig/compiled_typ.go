package konfig

type compiledTyp struct {
	Name           string
	FormatScript   ScriptConfig
	InlineAkte     bool
	ExecCommand    ScriptConfig
	EtikettenRules map[string]EtikettRule
	FileExtension  string
}

// TODO merge
func (ct *compiledTyp) Apply(kt KonfigTyp) {
	ct.FormatScript = kt.FormatScript
	ct.InlineAkte = kt.InlineAkte
	ct.ExecCommand = kt.ExecCommand
	ct.EtikettenRules = kt.EtikettenRules
	ct.FileExtension = kt.FileExtension
}
