package commands

// import (
// 	"flag"

// 	"code.linenisgreat.com/zit/src/alfa/errors"
// 	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
// 	"code.linenisgreat.com/zit/src/alfa/vim_cli_options_builder"
// 	"code.linenisgreat.com/zit/src/delta/gattung"
// 	"code.linenisgreat.com/zit/src/bravo/todo"
// 	"code.linenisgreat.com/zit/src/charlie/collections"
// 	"code.linenisgreat.com/zit/src/delta/script_value"
// 	"code.linenisgreat.com/zit/src/delta/gattungen"
// 	"code.linenisgreat.com/zit/src/echo/kennung"
// 	"code.linenisgreat.com/zit/src/juliett/objekte"
// 	"code.linenisgreat.com/zit/src/kilo/zettel"
// 	"code.linenisgreat.com/zit/src/kilo/cwd"
// 	"code.linenisgreat.com/zit/src/mike/store_fs"
// 	"code.linenisgreat.com/zit/src/november/umwelt"
// 	"code.linenisgreat.com/zit/src/papa/user_ops"
// )

// type Dupe struct {
// 	Edit      bool
// 	Delete    bool
// 	Dedupe    bool
// 	Count     int
// 	PrintOnly bool
// 	Filter    script_value.ScriptValue

// 	zettel.ProtoZettel
// }

// func init() {
// 	registerCommandWithQuery(
// 		"dupe",
// 		func(f *flag.FlagSet) CommandWithQuery {
// 			c := &Dupe{
// 				ProtoZettel: zettel.MakeEmptyProtoZettel(),
// 			}

// 			f.BoolVar(&c.Edit, "edit", true, "create a new empty zettel and open EDITOR or VISUAL for editing and then commit the resulting changes")

// 			c.ProtoZettel.AddToFlagSet(f)

// 			return c
// 		},
// 	)
// }

// func (c Dupe) DefaultGattungen() gattungen.Set {
// 	return gattungen.MakeSet(
// 		gattung.Zettel,
// 	)
// }

// func (c Dupe) RunWithQuery(u *umwelt.Umwelt, q kennung.MetaSet) (err error) {
// 	f := zettel.MakeObjekteTextFormat(
// 		u.StoreObjekten(),
// 		nil,
// 	)

// 	var zsc zettel.MutableSetCheckedOut

// 	var zts schnittstellen.MutableSet[*zettel.Transacted]

// 	if zts, err = c.readExistingFilesAsZettels(u, f, args...); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if c.Edit {
// 		var cwdFiles cwd.CwdFiles

// 		if cwdFiles, err = cwd.MakeCwdFilesAll(
// 			u.Konfig(),
// 			u.Standort().Cwd(),
// 		); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}

// 		options := store_fs.CheckoutOptions{
// 			Cwd:          cwdFiles,
// 			CheckoutMode: objekte.CheckoutModeObjekteAndAkte,
// 		}

// 		if zsc, err = u.StoreWorkingDirectory().Checkout(
// 			options,
// 			collections.WriterContainer[*zettel.Transacted](zts, collections.MakeErrStopIteration()),
// 		); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	}

// 	if c.Edit {
// 		ms := u.MakeMetaIdSet(
// 			kennung.MakeMatcherAlways(),
// 			gattungen.MakeSet(gattung.Zettel),
// 		)

// 		todo.Refactor("make this more stable by not using string query")
// 		if err = ms.Set(".zettel"); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}

// 		if err = c.editZettels(
// 			u,
// 			ms,
// 			zsc,
// 		); err != nil {
// 			err = errors.Wrap(err)
// 			return
// 		}
// 	}

// 	return
// }

// func (c Dupe) readExistingFilesAsZettels(
// 	u *umwelt.Umwelt,
// 	f zettel.ObjekteParser,
// 	args ...string,
// ) (zts schnittstellen.MutableSet[*zettel.Transacted], err error) {
// 	opCreateFromPath := user_ops.CreateFromPaths{
// 		Umwelt:      u,
// 		Format:      f,
// 		Filter:      c.Filter,
// 		Delete:      c.Delete,
// 		Dedupe:      c.Dedupe,
// 		ProtoZettel: c.ProtoZettel,
// 	}

// 	if zts, err = opCreateFromPath.Run(args...); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

// func (c Dupe) writeNewZettels(
// 	u *umwelt.Umwelt,
// 	f zettel.ObjekteFormatter,
// ) (zsc zettel.MutableSetCheckedOut, err error) {
// 	var cwdFiles cwd.CwdFiles

// 	if cwdFiles, err = cwd.MakeCwdFilesAll(
// 		u.Konfig(),
// 		u.Standort().Cwd(),
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	emptyOp := user_ops.WriteNewZettels{
// 		Umwelt:   u,
// 		CheckOut: c.Edit,
// 		CheckoutOptions: store_fs.CheckoutOptions{
// 			Cwd: cwdFiles,
// 		},
// 	}

// 	var defaultEtiketten kennung.EtikettSet

// 	if defaultEtiketten, err = u.DefaultEtiketten(); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	mes := c.Etiketten.MutableClone()
// 	defaultEtiketten.Each(mes.Add)
// 	c.Etiketten = mes.ImmutableClone()

// 	if zsc, err = emptyOp.RunMany(c.ProtoZettel, c.Count); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }

// func (c Dupe) editZettels(
// 	u *umwelt.Umwelt,
// 	ms kennung.MetaSet,
// 	zsc zettel.MutableSetCheckedOut,
// ) (err error) {
// 	if !c.Edit {
// 		errors.Log().Print("edit set to false, not editing")
// 		return
// 	}

// 	var filesZettelen []string

// 	if filesZettelen, err = zettel.ToSliceFilesZettelen(zsc); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	var cwdFiles cwd.CwdFiles

// 	if cwdFiles, err = cwd.MakeCwdFilesExactly(
// 		u.Konfig(),
// 		u.Standort().Cwd(),
// 		filesZettelen...,
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	openVimOp := user_ops.OpenVim{
// 		Options: vim_cli_options_builder.New().
// 			WithCursorLocation(2, 3).
// 			WithFileType("zit-zettel").
// 			WithInsertMode().
// 			Build(),
// 	}

// 	var fs []string

// 	if fs, err = cwdFiles.ZettelFiles(); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if _, err = openVimOp.Run(u, fs...); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if err = u.Reset(); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	if cwdFiles, err = cwd.MakeCwdFilesExactly(
// 		u.Konfig(),
// 		u.Standort().Cwd(),
// 		filesZettelen...,
// 	); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	checkinOp := user_ops.Checkin{}

// 	if err = checkinOp.Run(u, ms, cwdFiles); err != nil {
// 		err = errors.Wrap(err)
// 		return
// 	}

// 	return
// }
