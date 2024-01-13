package commands

// TODO-P4 bring back exec command
// type Exec struct{}

// func init() {
// 	registerCommand(
// 		"exec",
// 		func(f *flag.FlagSet) Command {
// 			c := &Exec{}

// 			return c
// 		},
// 	)
// }

// func (c Exec) Run(u *umwelt.Umwelt, args ...string) (err error) {
// 	var tz zettel.Transacted
// 	var executor script_config.RemoteScript
// 	var ar io.ReadCloser

// 	if tz, ar, executor, err = c.getZettel(u, args[0]); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	var pipePath string

// 	if pipePath, err = c.makeFifoPipe(tz); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	var cmd *exec.Cmd

// 	if cmd, err = c.makeCmd(executor, pipePath, args...); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	wg := &sync.WaitGroup{}
// 	wg.Add(2)

// 	go c.feedPipe(ar, wg, pipePath, tz)
// 	go c.exec(wg, pipePath, cmd)

// 	wg.Wait()
// 	errors.Log().Print("done waiting")

// 	return
// }

// func (c Exec) getZettel(
// 	u *umwelt.Umwelt,
// 	hString string,
// ) (
// 	tz zettel.Transacted,
// 	ar io.ReadCloser,
// 	executor script_config.RemoteScript,
// 	err error,
// ) {
// 	var h kennung.Hinweis

// 	if err = h.Set(hString); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	var zt *zettel.Transacted

// 	if zt, err = u.StoreObjekten().ReadOne(&h); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	tz = *zt

// 	typ := tz.GetTyp()

// 	typKonfig := u.Konfig().GetApproximatedTyp(typ).ApproximatedOrActual()

// 	if typKonfig == nil {
// 		err = errors.Normal(
// 			errors.Errorf("Typ does not have an exec-command set: %s", typ),
// 		)
// 		return
// 	}

// 	executor = typKonfig.Akte.ExecCommand

// 	if ar, err = u.StoreObjekten().AkteReader(
// 		tz.GetAkteSha(),
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

// func (c Exec) makeFifoPipe(tz zettel.Transacted) (p string, err error) {
// 	h := tz.Sku.GetKennung()
// 	var d string

// 	if d, err = os.MkdirTemp("", h.Kopf()); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	p = path.Join(d, h.Schwanz()+"."+tz.GetTyp().String())

// 	if err = syscall.Mknod(p, syscall.S_IFIFO|0o666, 0); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if err = os.Remove(p); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

// func (c Exec) makeCmd(
// 	executor script_config.RemoteScript,
// 	p string,
// 	args ...string,
// ) (cmd *exec.Cmd, err error) {
// 	cmdArgs := append([]string{p})

// 	if len(args) > 1 {
// 		cmdArgs = append(cmdArgs, args[1:]...)
// 	}

// 	if cmd, err = executor.Cmd(cmdArgs...); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	cmd.Stdin = os.Stdin
// 	cmd.Stderr = os.Stderr
// 	cmd.Stdout = os.Stdout

// 	return
// }

// func (c Exec) feedPipe(
// 	ar io.ReadCloser,
// 	wg *sync.WaitGroup,
// 	p string,
// 	tz zettel.Transacted,
// ) (err error) {
// 	defer wg.Done()
// 	var pipeFileWriter *os.File

// 	if pipeFileWriter, err = os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_APPEND,
// 0o777); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	defer errors.Deferred(&err, pipeFileWriter.Close)

// 	defer errors.Deferred(&err, ar.Close)

// 	if _, err = io.Copy(pipeFileWriter, ar); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	errors.Log().Print("done copying")

// 	return
// }

// func (c Exec) exec(
// 	wg *sync.WaitGroup,
// 	p string,
// 	cmd *exec.Cmd,
// ) (err error) {
// 	defer wg.Done()
// 	// var pipeFileReader *os.File

// 	// if pipeFileReader, err = os.OpenFile(pipePath, os.O_CREATE,
// 	// os.ModeNamedPipe); err != nil {
// 	// 	err = errors.Error(err)
// 	// 	return
// 	// }

// 	// defer pipeFileReader.Close()

// 	// cmd.ExtraFiles = append(cmd.ExtraFiles, pipeFileReader)

// 	errors.Log().Print("start running")

// 	if err = cmd.Run(); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	errors.Log().Print("done running")

// 	return
// }
