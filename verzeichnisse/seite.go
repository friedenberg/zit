package verzeichnisse

//func newPage(r1 io.Reader, dec gob.Decoder) (p *page, err error) {
//	p = &page{}

//	r := bufio.NewReader(r1)

//  dec.Decode(r1,

//	return
//}

//func (p *page) readAll(r *bufio.Reader) (err error) {

//}

//// func (s *shard) Remove(key string) {
//// 	s.rwLock.RLock()
//// 	defer s.rwLock.RUnlock()

//// 	s.hasChanges = true
//// 	delete(s.entries, key)
//// }

//func (s *shard) Set(key, value string) {
//	s.rwLock.RLock()
//	defer s.rwLock.RUnlock()

//	s.hasChanges = true
//	s.entries[key] = value
//}

//func (s shard) Read(key string) (value string, ok bool) {
//	s.rwLock.RLock()
//	defer s.rwLock.RUnlock()

//	if value, ok = s.entries[key]; ok {
//		return
//	}

//	return
//}

//func (s shard) flush(w io.WriteCloser) (err error) {
//	s.rwLock.RLock()
//	// maybe a mistake
//	defer s.rwLock.RUnlock()

//	if !s.hasChanges {
//		return
//	}

//	logz.Printf("flushing: %s", s.path)

//	var file *os.File

//	if file, err = open_file_guard.TempFile(); err != nil {
//		logz.Print(err)
//		err = errors.Error(err)
//		return
//	}

//	defer open_file_guard.Close(file)

//	var w io.WriteCloser

//	if w, err = s.Writer(file); err != nil {
//		logz.Print(err)
//		err = errors.Error(err)
//		return
//	}

//	defer w.Close()

//	for k, v := range s.entries {
//		var line string

//		if line, err = s.EntryToLine(Entry{k, v}); err != nil {
//			logz.Print(err)
//			err = errors.Error(err)
//			return
//		}

//		if _, err = io.WriteString(w, fmt.Sprintln(line)); err != nil {
//			logz.Print(err)
//			err = errors.Error(err)
//			return
//		}
//	}

//	//TODO-research should the file be closed before being renamed???
//	logz.Printf("renaming %s to %s", file.Name(), s.path)
//	if err = os.Rename(file.Name(), s.path); err != nil {
//		err = errors.Error(err)
//		return
//	}

//	logz.Print("done renaming")

//	return
//}

//func (s shard) All() (a []Entry) {
//	a = make([]Entry, 0, len(s.entries))

//	for k, v := range s.entries {
//		a = append(a, Entry{k, v})
//	}

//	return
//}
