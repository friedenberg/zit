package errors

import (
	log_package "log"
)

var (
	//TODO add native methods
	Panic  = log_package.Panic
	Output = log_package.Output
	Fatal  = log_package.Fatal
)
