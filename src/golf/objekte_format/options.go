package objekte_format

type Options struct {
	IncludeTai           bool
	IncludeVerzeichnisse bool
}

func (o Options) SansVerzeichnisse() Options {
	o.IncludeVerzeichnisse = false
	return o
}
