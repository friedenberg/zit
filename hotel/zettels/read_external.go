package zettels

import "os"

type ExternalZettel struct {
	Path     string
	AktePath string
	Hinweis  _Hinweis
	Zettel   _Zettel
}

func (zs *zettels) ReadExternal(options CheckinOptions, paths ...string) (out map[_Hinweis]ExternalZettel, err error) {
	out = make(map[_Hinweis]ExternalZettel)

	for _, p := range paths {
		var ez ExternalZettel

		if ez, err = zs.readExternalOne(options, p); err != nil {
			err = _Error(err)
			return
		}

		out[ez.Hinweis] = ez
	}

	return
}

func (zs zettels) readExternalOne(options CheckinOptions, p string) (ez ExternalZettel, err error) {
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

	if f, err = _Open(p); err != nil {
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
