package local_working_copy

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
	"code.linenisgreat.com/zit/go/zit/src/hotel/object_inventory_format"
	"code.linenisgreat.com/zit/go/zit/src/hotel/type_blobs"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
)

// TODO switch to using fd.Std
func (repo *Repo) MakeFormatFunc(
	format string,
	writer interfaces.WriterAndStringWriter,
) (output interfaces.FuncIter[*sku.Transacted], err error) {
	if writer == nil {
		writer = repo.GetUIFile()
	}

	if strings.HasPrefix(format, "type.") {
		return repo.makeTypFormatter(strings.TrimPrefix(format, "type."), writer)
	}

	switch format {
	case "tags-path":
		output = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				writer,
				tl.GetObjectId(),
				&tl.Metadata.Cache.TagPaths,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "tags-path-with-types":
		output = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				writer,
				tl.GetObjectId(),
				&tl.Metadata.Cache.TagPaths,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "query-path":
		output = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				writer,
				tl.GetObjectId(),
				tl.Metadata.Cache.QueryPath,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "box":
		p := repo.SkuFormatBoxTransactedNoColor()

		output = func(tl *sku.Transacted) (err error) {
			if _, err = p.EncodeStringTo(tl, writer); err != nil {
				err = errors.Wrap(err)
				return
			}

			if _, err = fmt.Fprintln(writer); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "box-archive":
		p := repo.MakePrinterBoxArchive(writer, repo.GetConfig().GetCLIConfig().PrintOptions.PrintTime)

		output = func(tl *sku.Transacted) (err error) {
			if err = p(tl); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "sha":
		output = func(tl *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(writer, tl.Metadata.Sha())
			return
		}

	case "sha-mutter":
		output = func(tl *sku.Transacted) (err error) {
			_, err = fmt.Fprintf(writer, "%s -> %s\n", tl.Metadata.Sha(), tl.Metadata.Mutter())
			return
		}

	case "tags-all":
		output = func(tl *sku.Transacted) (err error) {
			for _, es := range tl.Metadata.Cache.TagPaths.Paths {
				if _, err = fmt.Fprintf(writer, "%s: %s\n", tl.GetObjectId(), es); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			for _, es := range tl.Metadata.Cache.TagPaths.All {
				if _, err = fmt.Fprintf(writer, "%s: %s -> %s\n", tl.GetObjectId(), es.Tag, es.Parents); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		}

	case "tags-expanded":
		output = func(tl *sku.Transacted) (err error) {
			esImp := tl.GetMetadata().Cache.GetExpandedTags()
			// TODO-P3 determine if empty sets should be printed or not

			if _, err = fmt.Fprintln(
				writer,
				quiter.StringCommaSeparated(esImp),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "tags-implicit":
		output = func(tl *sku.Transacted) (err error) {
			esImp := tl.GetMetadata().Cache.GetImplicitTags()
			// TODO-P3 determine if empty sets should be printed or not

			if _, err = fmt.Fprintln(
				writer,
				quiter.StringCommaSeparated(esImp),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "tags":
		output = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				writer,
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
		output = func(tl *sku.Transacted) (err error) {
			if err = tl.Metadata.GetTags().EachPtr(func(e *ids.Tag) (err error) {
				_, err = fmt.Fprintln(writer, e)
				return
			}); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "description":
		output = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(writer, tl.GetMetadata().Description); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "text":
		fo := typed_blob_store.MakeTextFormatter(
			repo.GetStore().GetEnvRepo(),
			checkout_options.TextFormatterOptions{
				DoNotWriteEmptyDescription: true,
			},
			repo.GetConfig(),
		)

		output = func(tl *sku.Transacted) (err error) {
			_, err = fo.EncodeStringTo(tl, writer)
			return
		}

	case "object":
		fo := object_inventory_format.FormatForVersion(
			repo.GetConfig().GetImmutableConfig().GetStoreVersion(),
		)

		o := object_inventory_format.Options{
			Tai: true,
		}

		output = func(tl *sku.Transacted) (err error) {
			if _, err = fo.FormatPersistentMetadata(writer, tl, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "object-id-parent-tai":
		output = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintf(
				writer,
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
		output = func(tl *sku.Transacted) (err error) {
			if _, err = fmt.Fprintf(
				writer,
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
		output = func(tl *sku.Transacted) (err error) {
			ui.TodoP3("convert into an option")

			sh := tl.GetBlobSha()

			if sh.IsNull() {
				return
			}

			if _, err = fmt.Fprintf(
				writer,
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
		output = func(e *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				writer,
				&e.ObjectId,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "object-id-abbreviated":
		output = func(e *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(
				writer,
				&e.ObjectId,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "object-id-tai":
		output = func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(writer, e.StringObjectIdTai())
			return
		}

	case "sku-metadata-sans-tai":
		output = func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(
				writer,
				sku_fmt.StringMetadataSansTai(e),
			)
			return
		}

	case "metadata":
		fo, err := object_inventory_format.FormatForKeyError(
			object_inventory_format.KeyFormatV5Metadata,
		)

		errors.PanicIfError(err)

		output = func(e *sku.Transacted) (err error) {
			_, err = fo.WriteMetadataTo(writer, e)
			return
		}

	case "metadata-plus-mutter":
		fo, err := object_inventory_format.FormatForKeyError(
			object_inventory_format.KeyFormatV5MetadataObjectIdParent,
		)

		errors.PanicIfError(err)

		output = func(e *sku.Transacted) (err error) {
			_, err = fo.WriteMetadataTo(writer, e)
			return
		}

	case "debug":
		output = func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintf(writer, "%#v\n", e)
			return
		}

	case "log":
		output = repo.PrinterTransacted()

	case "json":
		enc := json.NewEncoder(writer)

		output = func(o *sku.Transacted) (err error) {
			var j sku_fmt.Json

			if err = j.FromTransacted(o, repo.GetStore().GetEnvRepo()); err != nil {
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
		enc := json.NewEncoder(writer)

		type tomlJson struct {
			sku_fmt.Json
			Blob map[string]interface{} `json:"blob"`
		}

		output = func(o *sku.Transacted) (err error) {
			var j tomlJson

			if err = j.FromTransacted(o, repo.GetStore().GetEnvRepo()); err != nil {
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
		enc := json.NewEncoder(writer)

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

		output = func(o *sku.Transacted) (err error) {
			var j sku_fmt.JsonWithUrl

			if j, err = sku_fmt.MakeJsonTomlBookmark(
				o,
				repo.GetStore().GetEnvRepo(),
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
		output = func(o *sku.Transacted) (err error) {
			fmt.Fprintln(writer, o.GetTai())
			return
		}

	case "blob":
		output = func(o *sku.Transacted) (err error) {
			var r sha.ReadCloser

			if r, err = repo.GetStore().GetEnvRepo().BlobReader(
				o.GetBlobSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.DeferredCloser(&err, r)

			if _, err = io.Copy(writer, r); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "text-sku-prefix":
		cliFmt := repo.SkuFormatBoxTransactedNoColor()

		output = func(o *sku.Transacted) (err error) {
			sb := &strings.Builder{}

			if _, err = cliFmt.EncodeStringTo(o, sb); err != nil {
				err = errors.Wrap(err)
				return
			}

			if repo.GetConfig().IsInlineType(o.GetType()) {
				var r sha.ReadCloser

				if r, err = repo.GetStore().GetEnvRepo().BlobReader(
					o.GetBlobSha(),
				); err != nil {
					err = errors.Wrap(err)
					return
				}

				defer errors.DeferredCloser(&err, r)

				if _, err = delim_io.CopyWithPrefixOnDelim(
					'\n',
					sb.String(),
					repo.GetOut(),
					r,
					true,
				); err != nil {
					err = errors.Wrap(err)
					return
				}
			} else {
				if _, err = io.WriteString(writer, sb.String()); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		}

	case "blob-sku-prefix":
		cliFmt := repo.SkuFormatBoxTransactedNoColor()

		output = func(o *sku.Transacted) (err error) {
			var r sha.ReadCloser

			if r, err = repo.GetStore().GetEnvRepo().BlobReader(
				o.GetBlobSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.DeferredCloser(&err, r)

			sb := &strings.Builder{}

			if _, err = cliFmt.EncodeStringTo(o, sb); err != nil {
				err = errors.Wrap(err)
				return
			}

			if _, err = delim_io.CopyWithPrefixOnDelim(
				'\n',
				sb.String(),
				repo.GetOut(),
				r,
				true,
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "shas":
		output = func(z *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(writer, &z.Metadata.Shas)
			return
		}

	case "mutter-sha":
		output = func(z *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(writer, z.Metadata.Mutter())
			return
		}

	case "probe-shas":
		output = func(z *sku.Transacted) (err error) {
			sh1 := sha.FromStringContent(z.GetObjectId().String())
			sh2 := sha.FromStringContent(z.GetObjectId().String() + z.GetTai().String())
			defer sha.GetPool().Put(sh1)
			defer sha.GetPool().Put(sh2)
			_, err = fmt.Fprintln(writer, z.GetObjectId(), sh1, sh2)
			return
		}

	case "mutter":
		p := repo.PrinterTransacted()

		output = func(z *sku.Transacted) (err error) {
			if z.Metadata.Mutter().IsNull() {
				return
			}

			if z, err = repo.GetStore().GetStreamIndex().ReadOneObjectIdTai(
				z.GetObjectId(),
				z.Metadata.Cache.ParentTai,
			); err != nil {
				fmt.Fprintln(writer, err)
				err = nil
				return
			}

			return p(z)
		}

	case "inventory-list":
		p := repo.MakePrinterBoxArchive(repo.GetUIFile(), true)

		output = func(o *sku.Transacted) (err error) {
			if err = p(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "inventory-list-sans-tai":
		p := repo.MakePrinterBoxArchive(repo.GetUIFile(), false)

		output = func(o *sku.Transacted) (err error) {
			if err = p(o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "blob-sha":
		output = func(o *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(writer, o.GetBlobSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "type":
		output = func(o *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(writer, o.GetType().String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "verzeichnisse":
		p := repo.PrinterTransacted()

		output = func(o *sku.Transacted) (err error) {
			sk := sku.GetTransactedPool().Get()
			defer sku.GetTransactedPool().Put(sk)

			if err = repo.GetStore().GetStreamIndex().ReadOneObjectId(
				o.GetObjectId(),
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
		e := json.NewEncoder(writer)

		output = func(o *sku.Transacted) (err error) {
			var a map[string]interface{}

			var r sha.ReadCloser

			if r, err = repo.GetStore().GetEnvRepo().BlobReader(
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
		output = func(o *sku.Transacted) (err error) {
			var a map[string]interface{}

			var r sha.ReadCloser

			if r, err = repo.GetStore().GetEnvRepo().BlobReader(
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

			e := toml.NewEncoder(writer)

			if err = e.Encode(&a); err != nil {
				err = errors.Wrap(err)
				return
			}

			if _, err = writer.Write([]byte("\x00")); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "debug-sku-metadata":
		output = func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(
				writer,
				sku.StringMetadataTai(e),
			)
			return
		}

	case "debug-sku":
		output = func(e *sku.Transacted) (err error) {
			_, err = fmt.Fprintln(writer, sku.StringTaiGenreObjectIdShaBlob(e))
			return
		}

	default:
		err = MakeErrUnsupportedFormatterValue(format, genres.None)
	}

	return
}

func (u *Repo) makeTypFormatter(
	v string,
	out io.Writer,
) (f interfaces.FuncIter[*sku.Transacted], err error) {
	typeBlobStore := u.GetStore().GetTypedBlobStore().Type

	if out == nil {
		out = u.GetUIFile()
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

			u.GetUI().Print(f.String())

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
