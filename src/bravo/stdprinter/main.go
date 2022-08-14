package stdprinter

import (
	"fmt"
	"os"
	"sync"

	"github.com/friedenberg/zit/src/alfa/logz"
)

type printerLine struct {
	file *os.File
	line string
}

var (
	printerChannel chan printerLine
	waitGroup      *sync.WaitGroup
)

func init() {
	printerChannel = make(chan printerLine)
	waitGroup = &sync.WaitGroup{}
	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()

		for printerLine := range printerChannel {
			fmt.Fprint(printerLine.file, printerLine.line)
		}
	}()
}

func WaitForPrinter() {
	logz.Print("waiting")
	defer logz.Print("done")

	close(printerChannel)
	waitGroup.Wait()
}
