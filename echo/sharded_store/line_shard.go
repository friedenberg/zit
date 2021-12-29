package sharded_store

type ShardLine struct{}

func (s ShardLine) LineToEntry(line string) (entry Entry, err error) {
	entry.Key = line

	return
}

func (s ShardLine) EntryToLine(entry Entry) (line string, err error) {
	line = entry.Key
	return
}
