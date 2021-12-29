package sharded_store

type EntryMap map[string]string

type Entry struct {
	Key, Value string
}
