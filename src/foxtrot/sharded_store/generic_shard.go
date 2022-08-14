package sharded_store

import (
	"fmt"
	"strings"

	"github.com/friedenberg/zit/src/bravo/errors"
)

type ShardGeneric struct{}

func (s ShardGeneric) LineToEntry(line string) (entry Entry, err error) {
	parts := strings.Split(strings.TrimSpace(line), " ")
	partCount := len(parts)

	if partCount != 2 {
		err = errors.Errorf("expected 2 parts, but got %d.", partCount)
		return
	}

	entry.Key = parts[0]
	entry.Value = parts[1]

	return
}

func (s ShardGeneric) EntryToLine(entry Entry) (line string, err error) {
	line = fmt.Sprintf("%s %s", entry.Key, entry.Value)
	return
}
