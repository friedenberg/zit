package open_file_guard

func ReadDirNames(p string) (names []string, err error) {
	d, err := Open(p)

	if err != nil {
		err = _Error(err)
		return
	}

	defer Close(d)

	if names, err = d.Readdirnames(0); err != nil {
		err = _Error(err)
		return
	}

	return
}
