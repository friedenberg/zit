package objekte_format

func MakeFormatterContextIncludeTai(
	c FormatterContext,
	include bool,
) FormatterContextIncludeTai {
	return formatterContextIncludeTai{
		FormatterContext: c,
		include:          include,
	}
}

type formatterContextIncludeTai struct {
	FormatterContext
	include bool
}

func (c formatterContextIncludeTai) IncludeTai() bool {
	return c.include
}
