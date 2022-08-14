package stdprinter

type Logger interface {
	// Fatal(v ...interface{})
	// Fatalf(format string, v ...interface{})

	// Panic(v ...interface{})
	// Panicf(format string, v ...interface{})
	// Panicln(v ...interface{})

	Print(v ...interface{})
	Printf(format string, v ...interface{})

	// Output(calldepth int, s string) error

	// Prefix() string
	// SetPrefix(prefix string)
}
