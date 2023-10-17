package checkout_options

type (
	Path interface {
		isPath()
	}

	path int
)

func (path) isPath() {}

const (
	PathDefault = path(iota)
	PathLeft
	PathMiddle
	PathRight
	PathTempLocal
	PathTempOS
)
