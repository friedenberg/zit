package user_ops

//import (
//	"io"
//	"os"

//	"github.com/friedenberg/zit/alfa/errors"
//	"github.com/friedenberg/zit/foxtrot/stored_zettel"
//)

//type ReadFromPaths struct {
//	Umwelt  _Umwelt
//	Store   _Store
//	Format  _ZettelFormatsText
//	Filter  _ScriptValue
//	Options _ZettelsCheckinOptions
//}

//type ReadFromPathsResults struct {
//	Zettelen []stored_zettel.External
//}

////func (c ReadFromPaths) Run(args ...string) (results ReadFromPathsResults, err error) {
////	toCreate := make([]_Zettel, 0, len(args))

////	for _, arg := range args {
////		var toAdd []_Zettel

////		if toAdd, err = c.zettelsFromPath(arg); err != nil {
////			err = _Errorf("zettel text format error for path: %s: %w", arg, err)
////			return
////		}

////		toCreate = append(toCreate, toAdd...)
////	}

////	for _, z := range toCreate {
////		var named _NamedZettel

////		if named, err = c.Store.Create(z); err != nil {
////			//TODO add file for error handling
////			c.handleStoreError(named, "", err)
////			err = nil
////			return
////		}

////		_Outf("[%s %s]\n", named.Hinweis, named.Sha)
////	}

////	return
////}

////func (c ReadFromPaths) zettelsFromPath(p string) (out []_Zettel, err error) {
////	var r io.Reader

////	log.Print("running")

////	if r, err = c.Filter.Run(p); err != nil {
////		err = _Error(err)
////		return
////	}

////	defer c.Filter.Close()

////	ctx := zettel.FormatContextRead{
////		In:                r,
////		AkteWriterFactory: c.Store,
////	}

////	if _, err = c.Format.ReadFrom(&ctx); err != nil {
////		err = _Error(err)
////		return
////	}

////	if ctx.RecoverableError != nil {
////		var errAkteInlineAndFilePath zettel_formats.ErrHasInlineAkteAndFilePath

////		if errors.As(ctx.RecoverableError, &errAkteInlineAndFilePath) {
////			var z1 _Zettel

////			if z1, err = errAkteInlineAndFilePath.Recover(); err != nil {
////				err = _Error(err)
////				return
////			}

////			out = append(out, z1)
////		} else {
////			err = _Errorf("unsupported recoverable error: %w", ctx.RecoverableError)
////			return
////		}
////	}

////	out = append(out, ctx.Zettel)

////	return
////}

//func (c ReadFromPaths) Run(paths ...string) (results ReadFromPathsResults, err error) {
//	results.Zettelen = make([]stored_zettel.External, 0, len(paths))

//	for _, p := range paths {
//		if c.Options.AddMdExtension {
//			p = p + ".md"
//		}

//		var ez stored_zettel.External

//		ez, err = c.readExternalOne(p)

//		if c.Options.IgnoreMissingHinweis && errors.Is(os.ErrNotExist, err) {
//			err = nil
//			out[ez.Hinweis] = _ZettelExternal{}
//			continue
//		} else if err != nil {
//			err = _Error(err)
//			return
//		}

//		out[ez.Hinweis] = ez
//	}

//	return
//}

//func (c ReadFromPaths) readExternalOne(p string) (ez _ZettelExternal, err error) {
//	ez.Path = p

//	head, tail := _IdHeadTailFromFileName(p)

//	if ez.Hinweis, err = _MakeBlindHinweis(head + "/" + tail); err != nil {
//		err = _Error(err)
//		return
//	}

//	c := _ZettelFormatContextRead{
//		AkteWriterFactory: zs,
//	}

//	var f *os.File

//	if !_FilesExist(p) {
//		err = os.ErrNotExist
//		return
//	}

//	var r io.Reader

//	if r, err = c.Filter.Run(p); err != nil {
//		err = _Error(err)
//		return
//	}

//	defer c.Filter.Close()

//	if f, err = os.Open(p); err != nil {
//		err = _Error(err)
//		return
//	}

//	defer _Close(f)

//	c.In = f

//	if _, err = options.Format.ReadFrom(&c); err != nil {
//		err = _Errorf("%s: %w", f.Name(), err)
//		return
//	}

//	ez.Zettel = c.Zettel
//	ez.AktePath = c.AktePath

//	return
//}
