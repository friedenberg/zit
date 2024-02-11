package alfred

import (
	"bufio"
	"encoding/json"
	"io"

	"code.linenisgreat.com/zit-go/src/alfa/errors"
)

type Writer struct {
	wBuf        *bufio.Writer
	jsonEncoder *json.Encoder

	chItem chan *Item
	chDone chan struct{}

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
		for i := range w.chItem {
			// var err error

			if _ = w.writeItem(i); err != nil {
				// err = errors.Wrap(err)
				// TODO-P4
			}
		}

		w.chDone <- struct{}{}
	}()

	return
}

func (w *Writer) WriteItem(i *Item) {
	w.chItem <- i
}

func (w *Writer) setAfterFirstWrite() {
	w.afterFirstWrite = true
}

func (w *Writer) writeItem(i *Item) (err error) {
	if i == nil {
		panic("item was nil")
	}

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
	close(w.chItem)
	<-w.chDone

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
