package umwelt

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"code.linenisgreat.com/chrest"
	"code.linenisgreat.com/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/src/alfa/schnittstellen"
	"code.linenisgreat.com/zit/src/alfa/toml"
	"code.linenisgreat.com/zit/src/bravo/iter"
	"code.linenisgreat.com/zit/src/bravo/log"
	"code.linenisgreat.com/zit/src/charlie/ohio"
	"code.linenisgreat.com/zit/src/delta/gattung"
	"code.linenisgreat.com/zit/src/delta/lua"
	"code.linenisgreat.com/zit/src/delta/sha"
	"code.linenisgreat.com/zit/src/delta/typ_akte"
	"code.linenisgreat.com/zit/src/echo/format"
	"code.linenisgreat.com/zit/src/echo/kennung"
	"code.linenisgreat.com/zit/src/golf/objekte_format"
	"code.linenisgreat.com/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/src/india/akten"
	"code.linenisgreat.com/zit/src/india/sku_fmt"
	"code.linenisgreat.com/zit/src/kilo/zettel"
)

func (u *Umwelt) MakeFormatFunc(
	v string,
	out schnittstellen.WriterAndStringWriter,
) (f schnittstellen.FuncIter[*sku.Transacted], err error) {
	if out == nil {
		out = u.Out()
	}

	if strings.HasPrefix(v, "typ.") {
		return u.makeTypFormatter(strings.TrimPrefix(v, "typ."), out)
	}

	switch v {
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
			for _, es := range tl.Metadatei.Verzeichnisse.Etiketten {
				if _, err = fmt.Fprintln(out, tl.GetKennung(), "->", es); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		}

	case "etiketten-expanded":
		f = func(tl *sku.Transacted) (err error) {
			esImp := tl.GetMetadatei().Verzeichnisse.GetExpandedEtiketten()
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
			esImp := tl.GetMetadatei().Verzeichnisse.GetImplicitEtiketten()
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
					tl.Metadatei.GetEtiketten(),
				),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "etiketten-newlines":
		f = func(tl *sku.Transacted) (err error) {
			if err = tl.Metadatei.GetEtiketten().EachPtr(func(e *kennung.Etikett) (err error) {
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
			if _, err = fmt.Fprintln(out, tl.GetMetadatei().Bezeichnung); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "text":
		fo := akten.MakeTextFormatter(u.GetStore().GetStandort(), u.Konfig())

		f = func(tl *sku.Transacted) (err error) {
			_, err = fo.WriteStringFormat(out, tl)
			return
		}

	case "objekte":
		fo := objekte_format.FormatForVersion(u.Konfig().GetStoreVersion())
		o := objekte_format.Options{
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
				tl.GetObjekteSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "kennung-akte-sha":
		f = func(tl *sku.Transacted) (err error) {
			errors.TodoP3("convert into an option")

			sh := tl.GetAkteSha()

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
		fo, err := objekte_format.FormatForKeyError("Metadatei")
		errors.PanicIfError(err)

		f = func(e *sku.Transacted) (err error) {
			_, err = fo.WriteMetadateiTo(out, e)
			return
		}

	case "metadatei-plus-mutter":
		fo, err := objekte_format.FormatForKeyError("MetadateiMutter")
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
				err = errors.Wrap(err)
				return
			}

			if err = enc.Encode(j); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "json-toml-bookmark":
		enc := json.NewEncoder(out)

		var chromeTabsRaw interface{}
		var req *http.Request

		if req, err = http.NewRequest("GET", "http://localhost/tabs", nil); err != nil {
			errors.PanicIfError(err)
		}

		var chrestConfig chrest.Config

		if err = chrestConfig.Read(); err != nil {
			errors.PanicIfError(err)
		}

		if chromeTabsRaw, err = chrest.AskChrome(chrestConfig, req); err != nil {
			errors.PanicIfError(err)
		}

		chromeTabs := chromeTabsRaw.([]interface{})

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

			if r, err = u.GetStore().GetStandort().AkteReader(
				o.GetAkteSha(),
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

			if r, err = u.GetStore().GetStandort().AkteReader(
				o.GetAkteSha(),
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

			if r, err = u.GetStore().GetStandort().AkteReader(
				o.GetAkteSha(),
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

			if _, err = ohio.CopyWithPrefixOnDelim('\n', sb.String(), out, r); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "bestandsaufnahme-sans-tai":
		be := sku_fmt.MakeFormatBestandsaufnahmePrinter(
			out,
			objekte_format.Default(),
			objekte_format.Options{ExcludeMutter: true},
		)

		f = func(o *sku.Transacted) (err error) {
			if _, err = be.Print(o); err != nil {
				err = errors.Wrap(err)
				return
			}

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

			if z, err = u.GetStore().GetVerzeichnisse().ReadOneEnnui(
				z.GetMetadatei().Mutter(),
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
			objekte_format.Default(),
			objekte_format.Options{Tai: true},
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
			objekte_format.Default(),
			objekte_format.Options{
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
			if _, err = fmt.Fprintln(out, o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "typ":
		f = func(o *sku.Transacted) (err error) {
			if _, err = fmt.Fprintln(out, o.GetTyp().String()); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "verzeichnisse":
		p := u.PrinterTransactedLike()

		f = func(o *sku.Transacted) (err error) {
			var sk *sku.Transacted

			if sk, err = u.GetStore().GetVerzeichnisse().ReadOneKennung(
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

			if r, err = u.GetStore().GetStandort().AkteReader(
				o.GetAkteSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.Deferred(&err, r.Close)

			d := toml.NewDecoder(r)

			if err = d.Decode(&a); err != nil {
				log.Err().Printf("%s: %s", o, err)
				err = nil
				return
			}

			a["description"] = o.Metadatei.Bezeichnung.String()
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

			if r, err = u.GetStore().GetStandort().AkteReader(
				o.GetAkteSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer errors.DeferredCloser(&err, r)

			d := toml.NewDecoder(r)

			if err = d.Decode(&a); err != nil {
				log.Err().Printf("%s: %s", o, err)
				err = nil
				return
			}

			a["description"] = o.Metadatei.Bezeichnung.String()
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
		err = MakeErrUnsupportedFormatterValue(v, gattung.Unknown)
	}

	return
}

func (u *Umwelt) makeTypFormatter(
	v string,
	out io.Writer,
) (f schnittstellen.FuncIter[*sku.Transacted], err error) {
	agp := u.GetStore().GetAkten().GetTypV0()

	if out == nil {
		out = u.Out()
	}

	switch v {
	case "formatters":
		f = func(o *sku.Transacted) (err error) {
			t := u.Konfig().GetApproximatedTyp(o.GetTyp())

			if !t.HasValue {
				return
			}

			tt := t.ActualOrNil()

			var ta *typ_akte.V0

			if ta, err = agp.GetAkte(tt.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer agp.PutAkte(ta)

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
		fo := zettel.MakeFormatterTypFormatterUTIGroups(u.Konfig(), agp)

		f = func(o *sku.Transacted) (err error) {
			if _, err = fo.Format(out, o); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "action-names":
		fan := typ_akte.MakeFormatterActionNames()

		f = func(o *sku.Transacted) (err error) {
			var akte *typ_akte.V0

			if akte, err = agp.GetAkte(o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer agp.PutAkte(akte)

			if _, err = fan.Format(out, akte); err != nil {
				err = errors.Wrap(err)
				return
			}

			return
		}

	case "hooks.on_pre_commit":
		f = func(o *sku.Transacted) (err error) {
			var akte *typ_akte.V0

			if akte, err = agp.GetAkte(o.GetAkteSha()); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer agp.PutAkte(akte)

			script, ok := akte.Hooks.(string)

			if !ok || script == "" {
				return
			}

			var vp *lua.VMPool

			if vp, err = u.GetStore().MakeLuaVMPool(script); err != nil {
				err = errors.Wrap(err)
				return
			}

			vm := vp.Get()
			defer vp.Put(vm)

			var tt *lua.LTable

			if tt, err = vm.GetTopTableOrError(); err != nil {
				err = errors.Wrap(err)
				return
			}

			f := vm.GetField(tt, "on_pre_commit")

			log.Out().Print(f.String())

			return
		}

	case "vim-syntax-type":
		f = func(o *sku.Transacted) (err error) {
			t := u.Konfig().GetApproximatedTyp(o.GetTyp()).ApproximatedOrActual()

			if t == nil || t.Kennung.IsEmpty() || t.GetAkteSha().IsNull() {
				ty := ""

				switch o.GetGattung() {
				case gattung.Typ, gattung.Etikett, gattung.Kasten, gattung.Konfig:
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

			var ta *typ_akte.V0

			if ta, err = u.GetStore().GetAkten().GetTypV0().GetAkte(
				t.GetAkteSha(),
			); err != nil {
				err = errors.Wrap(err)
				return
			}

			defer u.GetStore().GetAkten().GetTypV0().PutAkte(ta)

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
			gattung.Typ,
		)

		return
	}

	return
}
