package files

// func MakeFifoPipe() (p string, err error) {
// 	var d string

// 	if d, err = os.MkdirTemp("", ""); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	p = path.Join(d, h.Schwanz()+"."+tz.Named.Stored.Objekte.Typ.String())

// 	if err = syscall.Mknod(p, syscall.S_IFIFO|0666, 0); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if err = os.Remove(p); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }
