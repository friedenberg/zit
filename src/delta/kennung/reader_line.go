package kennung

// type ReaderLine struct {
// 	Set
// 	Expanders               Expanders
// 	ImplicitEtikettenGetter ImplicitEtikettenGetter
// }

// func (rl *ReaderLine) ReadFrom(r1 io.Reader) (n int64, err error) {
// 	errors.TodoP4("add expanders")
// 	rl.Set = MakeSet(nil, Expanders{}, nil, nil)
// 	r := bufio.NewReader(r1)

// 	for {
// 		var line string

// 		line, err = r.ReadString('\n')
// 		n += int64(len(line))

// 		switch {
// 		case err == nil:
// 			break

// 		case errors.IsEOF(err):
// 			err = nil
// 			return

// 		default:
// 			err = errors.Wrap(err)
// 			return
// 		}

// 		if line == "" {
// 			continue
// 		}

// 		if err = tryAddMatcher(
// 			&rl.Set,
// 			rl.Expanders,
// 			rl.ImplicitEtikettenGetter,
// 			line,
// 		); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	}
// }
