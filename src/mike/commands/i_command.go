package commands

import "github.com/friedenberg/zit/src/kilo/umwelt"

type Command interface {
	Run(*umwelt.Umwelt, ...string) error
}

type CommandSupportingErrors interface {
	HandleError(*umwelt.Umwelt, error)
}

type CommandWithArgPreprocessor interface {
	PreprocessArgs(*umwelt.Umwelt, []string) ([]string, error)
}

type CommandWithDescription interface {
	Description() string
}

type CommandV2 struct {
  Command
  WithCompletion
}
