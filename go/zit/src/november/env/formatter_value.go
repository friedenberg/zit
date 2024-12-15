package env

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"code.linenisgreat.com/chrest/go/src/bravo/client"
	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/go/zit/src/bravo/quiter"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/checkout_options"
	"code.linenisgreat.com/zit/go/zit/src/charlie/collections"
	"code.linenisgreat.com/zit/go/zit/src/charlie/delim_io"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/format"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/hotel/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/india/sku_fmt_debug"
	"code.linenisgreat.com/zit/go/zit/src/juliett/blob_store"
)

func (u *Local) MakeFormatFunc(
	v string,
	out interfaces.WriterAndStringWriter,
) (f interfaces.FuncIter[*sku.Transacted], err error) {
	if out == nil {
		out = u.Out()
	}

	if strings.HasPrefix(v, "type.") {
		return u.makeTypFormatter(strings.TrimPrefix(v, "type."), out)
	}

	switch v {
	case "tags-path":
		f = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				out,
				tl.GetObjectId(),
				&tl.Metadata.Cache.TagPaths,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "tags-path-with-types":
		f = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				out,
				tl.GetObjectId(),
				&tl.Metadata.Cache.TagPaths,
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
				tl.Metadata.Cache.QueryPath,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "box":
		p := u.SkuFormatBoxTransactedNoColor()

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

	case "box-archive":
		p := u.MakePrinterBoxArchive(out, u.GetConfig().PrintOptions.PrintTime)

		f = func(tl *sku.Transacted) (err error) {
			if err = p(tl); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "sha":
		f = func(tl *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(out, tl.Metadata.Sha())
			return
		}

	case "sha-mutter":
		f = func(tl *sku.Transacted) (err error) {
			_, err = fmt.Fprintf(out, "%s -> %s\n", tl.Metadata.Sha(), tl.Metadata.Mutter())
			return
		}

	case "tags-all":
		f = func(tl *sku.Transacted) (err error) {
			for _, es := range tl.Metadata.Cache.TagPaths.Paths {
				if _, err = fmt.Fprintf(out, "%s: %s\n", tl.GetObjectId(), es); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			for _, es := range tl.Metadata.Cache.TagPaths.All {
				if _, err = fmt.Fprintf(out, "%s: %s -> %s\n", tl.GetObjectId(), es.Tag, es.Parents); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		}

	case "tags-expanded":
		f = func(tl *sku.Transacted) (err error) {
			esImp := tl.GetMetadata().Cache.GetExpandedTags()
			// TODO-P3 determine if empty sets should be printed or not

			if _, err = fmt.Fprintln(
				out,
				quiter.StringCommaSeparated(esImp),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "tags-implicit":
		f = func(tl *sku.Transacted) (err error) {
			esImp := tl.GetMetadata().Cache.GetImplicitTags()
			// TODO-P3 determine if empty sets should be printed or not

			if _, err = fmt.Fprintln(
				out,
				quiter.StringCommaSeparated(esImp),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "tags":
		f = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				out,
				quiter.StringCommaSeparated(
					tl.Metadata.GetTags(),
				),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "tags-newlines":
		f = func(tl *sku.Transacted) (err error) {
			if err = tl.Metadata.GetTags().EachPtr(func(e *ids.Tag) (err error) {
				_, err = fmt.Fprintln(out, e)
				return
			}); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "description":
		f = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(out, tl.GetMetadata().Description); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "text":
		fo := blob_store.MakeTextFormatter(
			u.GetStore().GetDirectoryLayout(),
			checkout_options.TextFormatterOptions{
				DoNotWriteEmptyDescription: true,
			},
			u.GetConfig(),
		)

		f = func(tl *sku.Transacted) (err error) {
			_, err = fo.WriteStringFormat(out, tl)
			return
		}

	case "object":
		fo := object_inventory_format.FormatForVersion(u.GetConfig().GetStoreVersion())
		o := object_inventory_format.Options{
			Tai: true,
		}

		f = func(tl *sku.Transacted) (err error) {
			if _, err = fo.FormatPersistentMetadata(out, tl, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "object-id-parent-tai":
		f = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintf(
				out,
				"%s^@%s\n",
				&tl.ObjectId,
				tl.Metadata.Cache.ParentTai,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "object-id-sha":
		f = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintf(
				out,
				"%s@%s\n",
				&tl.ObjectId,
				tl.GetObjectSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "object-id-blob-sha":
		f = func(tl *sku.Transacted) (err error) {
			ui.TodoP3("convert into an option")

			sh := tl.GetBlobSha()

			if sh.IsNull() {
				return
			}

			if _, err = fmt.Fprintf(
				out,
				"%s %s\n",
				&tl.ObjectId,
				sh,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "object-id":
		f = func(e *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				out,
				&e.ObjectId,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "object-id-abbreviated":
		f = func(e *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				out,
				&e.ObjectId,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "object-id-tai":
		f = func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(out, e.StringObjectIdTai())
			return
		}

	case "sku-metadata-sans-tai":
		f = func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(
				out,
				sku_fmt.StringMetadataSansTai(e),
			)
			return
		}

	case "metadata":
		fo, err := object_inventory_format.FormatForKeyError(
			object_inventory_format.KeyFormatV5Metadata,
		)

		errors.PanicIfError(err)

		f = func(e *sku.Transacted) (err error) {
			_, err = fo.WriteMetadataTo(out, e)
			return
		}

	case "metadata-plus-mutter":
		fo, err := object_inventory_format.FormatForKeyError(
			object_inventory_format.KeyFormatV5MetadataObjectIdParent,
		)

		errors.PanicIfError(err)

		f = func(e *sku.Transacted) (err error) {
			_, err = fo.WriteMetadataTo(out, e)
			return
		}

	case "debug":
		f = func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintf(out, "%#v\n", e)
			return
		}

	case "log":
		f = u.PrinterTransacted()

	case "json":
		enc := json.NewEncoder(out)

		f = func(o *sku.Transacted) (err error) {
			var j sku_fmt.Json

			if err = j.FromTransacted(o, u.GetStore().GetDirectoryLayout()); err != nil {
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
			Blob map[string]interface{} `json:"blob"`
		}

		f = func(o *sku.Transacted) (err error) {
			var j tomlJson

			if err = j.FromTransacted(o, u.GetStore().GetDirectoryLayout()); err != nil {
				err = errors.Wrap(err)
				return
			}

			if err = toml.Unmarshal([]byte(j.Json.BlobString), &j.Blob); err != nil {
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

		var resp client.ResponseWithParsedJSONBody

		req := client.BrowserRequest{
			Method: "GET",
			Path:   "/tabs",
		}

		var b client.BrowserProxy

		if err = b.Read(); err != nil {
			errors.PanicIfError(err)
		}

		if resp, err = b.Request(req); err != nil {
			errors.PanicIfError(err)
		}

		tabs := resp.ParsedJSONBody.([]interface{})

		f = func(o *sku.Transacted) (err error) {
			var j sku_fmt.JsonWithUrl

			if j, err = sku_fmt.MakeJsonTomlBookmark(
				o,
				u.GetStore().GetDirectoryLayout(),
				tabs,
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

	case "blob":
		f = func(o *sku.Transacted) (err error) {
			var r sha.ReadCloser

			if r, err = u.GetStore().GetDirectoryLayout().BlobReader(
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

			if r, err = u.GetStore().GetDirectoryLayout().BlobReader(
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

	case "blob-sku-prefix":
		cliFmt := u.SkuFormatBoxTransactedNoColor()

		f = func(o *sku.Transacted) (err error) {
			var r sha.ReadCloser

			if r, err = u.GetStore().GetDirectoryLayout().BlobReader(
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

			if _, err = delim_io.CopyWithPrefixOnDelim(
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

	case "shas":
		f = func(z *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(out, &z.Metadata.Shas)
			return
		}

	case "mutter-sha":
		f = func(z *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(out, z.Metadata.Mutter())
			return
		}

	case "probe-shas":
		f = func(z *sku.Transacted) (err error) {
			sh1 := sha.FromString(z.GetObjectId().String())
			sh2 := sha.FromString(z.GetObjectId().String() + z.GetTai().String())
			defer sha.GetPool().Put(sh1)
			defer sha.GetPool().Put(sh2)
			_, err = fmt.Fprintln(out, z.GetObjectId(), sh1, sh2)
			return
		}

	case "mutter":
		p := u.PrinterTransacted()

		f = func(z *sku.Transacted) (err error) {
			if z.Metadata.Mutter().IsNull() {
				return
			}

			if z, err = u.GetStore().GetStreamIndex().ReadOneObjectIdTai(
				z.GetObjectId(),
				z.Metadata.Cache.ParentTai,
			); err != nil {
				fmt.Fprintln(out, err)
				err = nil
				return
			}

			return p(z)
		}

	case "inventory-list":
		p := u.MakePrinterBoxArchive(u.Out(), true)

		f = func(o *sku.Transacted) (err error) {
			if err = p(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "inventory-list-sans-tai":
		p := u.MakePrinterBoxArchive(u.Out(), false)

		f = func(o *sku.Transacted) (err error) {
			if err = p(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "blob-sha":
		f = func(o *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(out, o.GetBlobSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "type":
		f = func(o *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(out, o.GetType().String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "verzeichnisse":
		p := u.PrinterTransacted()

		f = func(o *sku.Transacted) (err error) {
			sk := sku.GetTransactedPool().Get()
			defer sku.GetTransactedPool().Put(sk)

			if err = u.GetStore().GetStreamIndex().ReadOneObjectId(
				o.ObjectId.String(),
				sk,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

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

			if r, err = u.GetStore().GetDirectoryLayout().BlobReader(
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

			a["description"] = o.Metadata.Description.String()
			a["identifier"] = o.ObjectId.String()

			if err = e.Encode(&a); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "toml":
		ui.TodoP3("limit to only zettels supporting toml")
		f = func(o *sku.Transacted) (err error) {
			var a map[string]interface{}

			var r sha.ReadCloser

			if r, err = u.GetStore().GetDirectoryLayout().BlobReader(
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

			a["description"] = o.Metadata.Description.String()
			a["identifier"] = o.ObjectId.String()

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

	case "debug-sku-metadata":
		f = func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(
				out,
				sku_fmt_debug.StringMetadataTai(e),
			)
			return
		}

	case "debug-sku":
		f = func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(out, sku_fmt_debug.StringTaiGenreObjectIdShaBlob(e))
			return
		}

	default:
		err = MakeErrUnsupportedFormatterValue(v, genres.None)
	}

	return
}

func (u *Local) makeTypFormatter(
	v string,
	out io.Writer,
) (f interfaces.FuncIter[*sku.Transacted], err error) {
	typeBlobStore := u.GetStore().GetBlobStore().GetType()

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

			var ta type_blobs.Blob

			if ta, _, err = typeBlobStore.ParseTypedBlob(
				tt.GetType(),
				tt.GetBlobSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer typeBlobStore.PutTypedBlob(tt.GetType(), ta)

			lw := format.MakeLineWriter()

			for fn, f := range ta.GetFormatters() {
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
		fo := sku_fmt.MakeFormatterTypFormatterUTIGroups(u.GetStore(), typeBlobStore)

		f = func(o *sku.Transacted) (err error) {
			if _, err = fo.Format(out, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "hooks.on_pre_commit":
		f = func(o *sku.Transacted) (err error) {
			var blob type_blobs.Blob

			if blob, _, err = typeBlobStore.ParseTypedBlob(
				o.GetType(),
				o.GetBlobSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer typeBlobStore.PutTypedBlob(o.GetType(), blob)

			script := blob.GetStringLuaHooks()

			if script == "" {
				return
			}

			// TODO switch to typed variant
			var vp sku.LuaVMPoolV1

			if vp, err = u.GetStore().MakeLuaVMPoolV1(o, script); err != nil {
				err = errors.Wrap(err)
				return
			}

			var vm *sku.LuaVMV1

			if vm, err = vp.Get(); err != nil {
				err = errors.Wrap(err)
				return
			}

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

			if t == nil || t.ObjectId.IsEmpty() || t.GetBlobSha().IsNull() {
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

			var ta type_blobs.Blob

			if ta, _, err = typeBlobStore.ParseTypedBlob(
				t.GetType(),
				t.GetBlobSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer typeBlobStore.PutTypedBlob(t.GetType(), ta)

			if _, err = fmt.Fprintln(
				out,
				ta.GetVimSyntaxType(),
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
