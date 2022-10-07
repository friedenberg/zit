package alfred

import (
	"bufio"
	"encoding/json"
	"io"

	"github.com/friedenberg/zit/src/alfa/errors"
)

type Writer struct {
	wBuf            *bufio.Writer
	jsonEncoder     *json.Encoder

	chItem          chan *Item
	chDone          chan struct{}

	afterFirstWrite bool
	ItemPool
}

func NewWriter(out io.Writer) (w *Writer, err error) {
	w = &Writer{
		wBuf:     bufio.NewWriter(out),
		chItem:   make(chan *Item),
		chDone:   make(chan struct{}),
		ItemPool: MakeItemPool(),
	}

	w.jsonEncoder = json.NewEncoder(w.wBuf)

	if err = w.open(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (w *Writer) open() (err error) {
	_, err = w.wBuf.WriteString("{\"items\":[\n")

	go func() {
		errors.Print("write item loop")
		defer errors.Print("done with write item loop")

		for i := range w.chItem {
			// var err error

			errors.Print("running")
			if _ = w.writeItem(i); err != nil {
				// err = errors.Wrap(err)
				//TODO
			}
		}

		w.chDone <- struct{}{}
	}()

	return
}

func (w *Writer) WriteItem(i *Item) {
	errors.Print("writing")
	defer errors.Print("done writing")

	w.chItem <- i
}

func (w *Writer) setAfterFirstWrite() {
	w.afterFirstWrite = true
}

func (w *Writer) writeItem(i *Item) (err error) {
	defer w.setAfterFirstWrite()
	defer w.Put(i)

	if w.afterFirstWrite {
		if _, err = io.WriteString(w.wBuf, ","); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	if err = w.jsonEncoder.Encode(i); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (w *Writer) Close() (err error) {
	errors.Print("waiting to close")
	close(w.chItem)
	<-w.chDone
	errors.Print("closing")

	if _, err = w.wBuf.WriteString("]}\n"); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = w.wBuf.Flush(); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}
