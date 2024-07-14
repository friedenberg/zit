package umwelt

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"code.linenisgreat.com/chrest/go/chrest"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/delta/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/kilo/zettel"
)

func (u *Umwelt) MakeFormatFunc(
	v string,
	out interfaces.WriterAndStringWriter,
) (f interfaces.FuncIter[*sku.Transacted], err error) {
	if out == nil {
		out = u.Out()
	}

	if strings.HasPrefix(v, "typ.") {
		return u.makeTypFormatter(strings.TrimPrefix(v, "typ."), out)
	}

	switch v {
	case "etiketten-path":
		f = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				out,
				tl.GetObjectId(),
				&tl.Metadatei.Cache.TagPaths,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "etiketten-path-with-types":
		f = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				out,
				tl.GetObjectId(),
				&tl.Metadatei.Cache.TagPaths,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "query-path":
		f = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				out,
				tl.GetObjectId(),
				tl.Metadatei.Cache.QueryPath,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "organize":
		p := u.SkuFmtOrganize()

		f = func(tl *sku.Transacted) (err error) {
			if _, err = p.WriteStringFormat(out, tl); err != nil {
				err = errors.Wrap(err)
				return
			}

			if _, err = fmt.Fprintln(out); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "sha":
		f = func(tl *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(out, tl.Metadatei.Sha())
			return
		}

	case "sha-mutter":
		f = func(tl *sku.Transacted) (err error) {
			_, err = fmt.Fprintf(out, "%s -> %s\n", tl.Metadatei.Sha(), tl.Metadatei.Mutter())
			return
		}

	case "etiketten-all":
		f = func(tl *sku.Transacted) (err error) {
			for _, es := range tl.Metadatei.Cache.TagPaths.Paths {
				if _, err = fmt.Fprintf(out, "%s: %s\n", tl.GetObjectId(), es); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			for _, es := range tl.Metadatei.Cache.TagPaths.All {
				if _, err = fmt.Fprintf(out, "%s: %s -> %s\n", tl.GetObjectId(), es.Etikett, es.Parents); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		}

	case "etiketten-expanded":
		f = func(tl *sku.Transacted) (err error) {
			esImp := tl.GetMetadata().Cache.GetExpandedTags()
			// TODO-P3 determine if empty sets should be printed or not

			if _, err = fmt.Fprintln(
				out,
				iter.StringCommaSeparated(esImp),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "etiketten-implicit":
		f = func(tl *sku.Transacted) (err error) {
			esImp := tl.GetMetadata().Cache.GetImplicitTags()
			// TODO-P3 determine if empty sets should be printed or not

			if _, err = fmt.Fprintln(
				out,
				iter.StringCommaSeparated(esImp),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "etiketten":
		f = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				out,
				iter.StringCommaSeparated(
					tl.Metadatei.GetTags(),
				),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "etiketten-newlines":
		f = func(tl *sku.Transacted) (err error) {
			if err = tl.Metadatei.GetTags().EachPtr(func(e *ids.Tag) (err error) {
				_, err = fmt.Fprintln(out, e)
				return
			}); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "bezeichnung":
		f = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(out, tl.GetMetadata().Description); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "text":
		fo := blob_store.MakeTextFormatter(
			checkout_options.TextFormatterOptions{
				DoNotWriteEmptyBezeichnung: true,
			},
			u.GetStore().GetStandort(),
			u.GetKonfig(),
		)

		f = func(tl *sku.Transacted) (err error) {
			_, err = fo.WriteStringFormat(out, tl)
			return
		}

	case "objekte":
		fo := object_inventory_format.FormatForVersion(u.GetKonfig().GetStoreVersion())
		o := object_inventory_format.Options{
			Tai: true,
		}

		f = func(tl *sku.Transacted) (err error) {
			if _, err = fo.FormatPersistentMetadatei(out, tl, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "kennung-sha":
		f = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintf(
				out,
				"%s@%s\n",
				&tl.Kennung,
				tl.GetObjectSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "kennung-akte-sha":
		f = func(tl *sku.Transacted) (err error) {
			errors.TodoP3("convert into an option")

			sh := tl.GetBlobSha()

			if sh.IsNull() {
				return
			}

			if _, err = fmt.Fprintf(
				out,
				"%s %s\n",
				&tl.Kennung,
				sh,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "kennung":
		f = func(e *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				out,
				&e.Kennung,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "kennung-tai":
		f = func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(out, e.StringKennungTai())
			return
		}

	case "sku-metadatei-sans-tai":
		f = func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(
				out,
				sku_fmt.StringMetadateiSansTai(e),
			)
			return
		}

	case "sku-metadatei":
		f = func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(
				out,
				sku_fmt.StringMetadatei(e),
			)
			return
		}

	case "sku":
		f = func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(out, sku_fmt.String(e))
			return
		}

	case "metadatei":
		fo, err := object_inventory_format.FormatForKeyError("Metadatei")
		errors.PanicIfError(err)

		f = func(e *sku.Transacted) (err error) {
			_, err = fo.WriteMetadateiTo(out, e)
			return
		}

	case "metadatei-plus-mutter":
		fo, err := object_inventory_format.FormatForKeyError("MetadateiMutter")
		errors.PanicIfError(err)

		f = func(e *sku.Transacted) (err error) {
			_, err = fo.WriteMetadateiTo(out, e)
			return
		}

	case "debug":
		f = func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintf(out, "%#v\n", e)
			return
		}

	case "log":
		f = u.PrinterTransactedLike()

		// case "objekte":
		// 	f := Format{}

		// 	f = func(o TransactedLikePtr) (err error) {
		// 		if _, err = f.Format(out, &o.Objekte); err != nil {
		// 			err = errors.Wrap(err)
		// 			return
		// 		}

		// 		return
		// 	}
	case "json":
		enc := json.NewEncoder(out)

		f = func(o *sku.Transacted) (err error) {
			var j sku_fmt.Json

			if err = j.FromTransacted(o, u.GetStore().GetStandort()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = enc.Encode(j); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "toml-json":
		enc := json.NewEncoder(out)

		type tomlJson struct {
			sku_fmt.Json
			Akte map[string]interface{} `json:"akte"`
		}

		f = func(o *sku.Transacted) (err error) {
			var j tomlJson

			if err = j.FromTransacted(o, u.GetStore().GetStandort()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = toml.Unmarshal([]byte(j.Json.Akte), &j.Akte); err != nil {
				err = nil

				if err = enc.Encode(j.Json); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
				// err = errors.Wrap(err)
				// return
			}

			if err = enc.Encode(j); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "json-toml-bookmark":
		enc := json.NewEncoder(out)

		var resp chrest.ResponseWithParsedJSONBody

		req := chrest.BrowserRequest{
			Method: "GET",
			Path:   "/tabs",
		}

		var b chrest.Browser

		if err = b.Read(); err != nil {
			errors.PanicIfError(err)
		}

		if resp, err = b.Request(req); err != nil {
			errors.PanicIfError(err)
		}

		chromeTabs := resp.ParsedJSONBody.([]interface{})

		f = func(o *sku.Transacted) (err error) {
			var j sku_fmt.JsonWithUrl

			if j, err = sku_fmt.MakeJsonTomlBookmark(
				o,
				u.GetStore().GetStandort(),
				chromeTabs,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = enc.Encode(j); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "tai":
		f = func(o *sku.Transacted) (err error) {
			fmt.Fprintln(out, o.GetTai())
			return
		}

	case "akte":
		f = func(o *sku.Transacted) (err error) {
			var r sha.ReadCloser

			if r, err = u.GetStore().GetStandort().BlobReader(
				o.GetBlobSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.DeferredCloser(&err, r)

			if _, err = io.Copy(out, r); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "text-sku-prefix":
		f = func(o *sku.Transacted) (err error) {
			var r sha.ReadCloser

			if r, err = u.GetStore().GetStandort().BlobReader(
				o.GetBlobSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.DeferredCloser(&err, r)

			if _, err = io.Copy(out, r); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "akte-sku-prefix":
		cliFmt := u.StringFormatWriterSkuTransactedShort()

		f = func(o *sku.Transacted) (err error) {
			var r sha.ReadCloser

			if r, err = u.GetStore().GetStandort().BlobReader(
				o.GetBlobSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.DeferredCloser(&err, r)

			sb := &strings.Builder{}

			if _, err = cliFmt.WriteStringFormat(sb, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			if _, err = ohio.CopyWithPrefixOnDelim(
				'\n',
				sb.String(),
				out,
				r,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "bestandsaufnahme-sans-tai":
		be := sku_fmt.MakeFormatBestandsaufnahmePrinter(
			out,
			object_inventory_format.Default(),
			object_inventory_format.Options{ExcludeMutter: true},
		)

		f = func(o *sku.Transacted) (err error) {
			if _, err = be.Print(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "shas":
		f = func(z *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(out, &z.Metadatei.Shas)
			return
		}

	case "mutter-sha":
		f = func(z *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(out, z.Metadatei.Mutter())
			return
		}

	case "mutter":
		p := u.PrinterTransactedLike()

		f = func(z *sku.Transacted) (err error) {
			if z.Metadatei.Mutter().IsNull() {
				return
			}

			if z, err = u.GetStore().GetVerzeichnisse().ReadOneObjectSha(
				z.GetMetadata().Mutter(),
			); err != nil {
				fmt.Fprintln(out, err)
				err = nil
				return
			}

			return p(z)
		}

	case "bestandsaufnahme":
		fo := sku_fmt.MakeFormatBestandsaufnahmePrinter(
			out,
			object_inventory_format.Default(),
			object_inventory_format.Options{Tai: true},
		)

		f = func(o *sku.Transacted) (err error) {
			if _, err = fo.Print(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "bestandsaufnahme-verzeichnisse":
		fo := sku_fmt.MakeFormatBestandsaufnahmePrinter(
			out,
			object_inventory_format.Default(),
			object_inventory_format.Options{
				Tai:           true,
				Verzeichnisse: true,
			},
		)

		f = func(o *sku.Transacted) (err error) {
			if _, err = fo.Print(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "akte-sha":
		f = func(o *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(out, o.GetBlobSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "typ":
		f = func(o *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(out, o.GetType().String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "verzeichnisse":
		p := u.PrinterTransactedLike()

		f = func(o *sku.Transacted) (err error) {
			var sk *sku.Transacted

			if sk, err = u.GetStore().GetVerzeichnisse().ReadOneObjectId(
				&o.Kennung,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer sku.GetTransactedPool().Put(sk)

			if err = p(sk); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "json-blob":
		e := json.NewEncoder(out)

		f = func(o *sku.Transacted) (err error) {
			var a map[string]interface{}

			var r sha.ReadCloser

			if r, err = u.GetStore().GetStandort().BlobReader(
				o.GetBlobSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.Deferred(&err, r.Close)

			d := toml.NewDecoder(r)

			if err = d.Decode(&a); err != nil {
				ui.Err().Printf("%s: %s", o, err)
				err = nil
				return
			}

			a["description"] = o.Metadatei.Description.String()
			a["identifier"] = o.Kennung.String()

			if err = e.Encode(&a); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "toml":
		errors.TodoP3("limit to only zettels supporting toml")
		f = func(o *sku.Transacted) (err error) {
			var a map[string]interface{}

			var r sha.ReadCloser

			if r, err = u.GetStore().GetStandort().BlobReader(
				o.GetBlobSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.DeferredCloser(&err, r)

			d := toml.NewDecoder(r)

			if err = d.Decode(&a); err != nil {
				ui.Err().Printf("%s: %s", o, err)
				err = nil
				return
			}

			a["description"] = o.Metadatei.Description.String()
			a["identifier"] = o.Kennung.String()

			e := toml.NewEncoder(out)

			if err = e.Encode(&a); err != nil {
				err = errors.Wrap(err)
				return
			}

			if _, err = out.Write([]byte("\x00")); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	default:
		err = MakeErrUnsupportedFormatterValue(v, genres.Unknown)
	}

	return
}

func (u *Umwelt) makeTypFormatter(
	v string,
	out io.Writer,
) (f interfaces.FuncIter[*sku.Transacted], err error) {
	agp := u.GetStore().GetAkten().GetTypeV0()

	if out == nil {
		out = u.Out()
	}

	switch v {
	case "formatters":
		f = func(o *sku.Transacted) (err error) {
			var tt *sku.Transacted

			if tt, err = u.GetStore().ReadTransactedFromObjectId(o.GetType()); err != nil {
				err = errors.Wrap(err)
				return
			}

			var ta *type_blobs.V0

			if ta, err = agp.GetBlob(tt.GetBlobSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer agp.PutBlob(ta)

			lw := format.MakeLineWriter()

			for fn, f := range ta.Formatters {
				fe := f.FileExtension

				if fe == "" {
					fe = fn
				}

				lw.WriteFormat("%s\t%s", fn, fe)
			}

			if _, err = lw.WriteTo(out); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "formatter-uti-groups":
		fo := zettel.MakeFormatterTypFormatterUTIGroups(u.GetStore(), agp)

		f = func(o *sku.Transacted) (err error) {
			if _, err = fo.Format(out, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "action-names":
		fan := type_blobs.MakeFormatterActionNames()

		f = func(o *sku.Transacted) (err error) {
			var akte *type_blobs.V0

			if akte, err = agp.GetBlob(o.GetBlobSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer agp.PutBlob(akte)

			if _, err = fan.Format(out, akte); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "hooks.on_pre_commit":
		f = func(o *sku.Transacted) (err error) {
			var akte *type_blobs.V0

			if akte, err = agp.GetBlob(o.GetBlobSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer agp.PutBlob(akte)

			script, ok := akte.Hooks.(string)

			if !ok || script == "" {
				return
			}

			var vp sku_fmt.LuaVMPool

			if vp, err = u.GetStore().MakeLuaVMPool(o, script); err != nil {
				err = errors.Wrap(err)
				return
			}

			vm := vp.Get()
			defer vp.Put(vm)

			f := vm.GetField(vm.Top, "on_pre_commit")

			ui.Out().Print(f.String())

			return
		}

	case "vim-syntax-type":
		f = func(o *sku.Transacted) (err error) {
			var t *sku.Transacted

			if t, err = u.GetStore().ReadTransactedFromObjectId(o.GetType()); err != nil {
				if collections.IsErrNotFound(err) {
					err = nil
				} else {
					err = errors.Wrap(err)
					return
				}
			}

			if t == nil || t.Kennung.IsEmpty() || t.GetBlobSha().IsNull() {
				ty := ""

				switch o.GetGenre() {
				case genres.Type, genres.Tag, genres.Repo, genres.Config:
					ty = "toml"

				default:
					// TODO zettel default typ
				}

				if _, err = fmt.Fprintln(out, ty); err != nil {
					err = errors.Wrap(err)
					return
				}

				return
			}

			var ta *type_blobs.V0

			if ta, err = u.GetStore().GetAkten().GetTypeV0().GetBlob(
				t.GetBlobSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer u.GetStore().GetAkten().GetTypeV0().PutBlob(ta)

			if _, err = fmt.Fprintln(
				out,
				ta.VimSyntaxType,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	default:
		err = MakeErrUnsupportedFormatterValue(
			v,
			genres.Type,
		)

		return
	}

	return
}
