package type_blob

type LifecycleHooksV0 struct {
	New    interface{} `toml:"new"`
	Commit interface{} `toml:"commit"`
}

func (o *LifecycleHooksV0) Reset() {
	o.New = ""
	o.Commit = ""
}
