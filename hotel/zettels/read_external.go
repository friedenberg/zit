package zettels

import (
	"os"
)

func (zs *zettels) ReadCheckedOut(options CheckinOptions, paths ...string) (out map[_Hinweis]_ZettelCheckedOut, err error) {
	out = make(map[_Hinweis]_ZettelCheckedOut)
	var external map[_Hinweis]_ZettelExternal

	if external, err = zs.ReadExternal(options, paths...); err != nil {
		err = _Error(err)
		return
	}

	for h, ez := range external {
		var named _NamedZettel

		if named, err = zs.Read(h); err != nil {
			err = _Error(err)
			return
		}

		out[h] = _ZettelCheckedOut{
			External: ez,
			Internal: named,
		}
	}

	return
}

func (zs *zettels) ReadExternal(options CheckinOptions, paths ...string) (out map[_Hinweis]_ZettelExternal, err error) {
	out = make(map[_Hinweis]_ZettelExternal)

	for _, p := range paths {
		if options.AddMdExtension {
			p = p + ".md"
		}

		var ez _ZettelExternal

		ez, err = zs.readExternalOne(options, p)

		if options.IgnoreMissingHinweis && _ErrorsIs(os.ErrNotExist, err) {
			err = nil
			out[ez.Hinweis] = _ZettelExternal{}
			continue
		} else if err != nil {
			err = _Error(err)
			return
		}

		// _Errf("[%s] (read external)\n", ez.Hinweis)
		out[ez.Hinweis] = ez
	}

	return
}

func (zs zettels) readExternalOne(options CheckinOptions, p string) (ez _ZettelExternal, err error) {
	ez.Path = p

	head, tail := _IdHeadTailFromFileName(p)

	if ez.Hinweis, err = _MakeBlindHinweis(head + "/" + tail); err != nil {
		err = _Error(err)
		return
	}

	c := _ZettelFormatContextRead{
		AkteWriterFactory: zs,
	}

	var f *os.File

	if !_FilesExist(p) {
		err = os.ErrNotExist
		return
	}

	if f, err = os.Open(p); err != nil {
		err = _Error(err)
		return
	}

	defer _Close(f)

	c.In = f

	if _, err = options.Format.ReadFrom(&c); err != nil {
		err = _Errorf("%s: %w", f.Name(), err)
		return
	}

	ez.Zettel = c.Zettel
	ez.AktePath = c.AktePath

	return
}
