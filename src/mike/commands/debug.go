package commands

func (c command) SetDebug() (d func()) {
	//TODO-P2 this is currently broken entirely
	return nil
	// df := make([]func(), 0)
	// d = func() {
	// 	logz.Printf("running debug closers: %d", len(df))
	// 	for i := len(df) - 1; i >= 0; i-- {
	// 		logz.Printf("running debug closer: %s", df[i])
	// 		df[i]()
	// 	}
	// }

	// debug.SetGCPercent(-1)

	// f, _ := open_file_guard.Create("build/cpu1.pprof")
	// df = append(df, func() { open_file_guard.Close(f) })

	// f1, _ := open_file_guard.Create("build/trace")
	// df = append(df, func() { logz.Print(); open_file_guard.Close(f1); logz.Print() })

	// pprof.StartCPUProfile(f)
	// df = append(df, func() { pprof.StopCPUProfile() })

	// trace.Start(f1)
	// df = append(df, func() { trace.Stop() })

	// return
}
