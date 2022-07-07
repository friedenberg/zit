package log

import (
	"fmt"
	"log"
)

func Print(vs ...interface{}) {
	log.Print(vs...)
}

func Printf(f string, vs ...interface{}) {
	log.Printf(f, vs...)
}

func PrintDebug(vs ...interface{}) {
	for _, v := range vs {
		log.Output(2, fmt.Sprintf("%#v", v))
	}
}
