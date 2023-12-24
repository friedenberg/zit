package objekte_format

type Options struct {
	IncludeTai           bool
	IncludeVerzeichnisse bool
	PrintFinalSha        bool
}

func (o Options) SansVerzeichnisse() Options {
	o.IncludeVerzeichnisse = false
	return o
}
