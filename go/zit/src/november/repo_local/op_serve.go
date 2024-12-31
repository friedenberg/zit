package repo_local

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/india/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
)

func (env *Repo) InitializeListener(
	network, address string,
) (listener net.Listener, err error) {
	var config net.ListenConfig

	switch network {
	case "unix":
		if listener, err = env.InitializeUnixSocket(config, address); err != nil {
			err = errors.Wrap(err)
			return
		}

	case "tcp":
		if listener, err = config.Listen(env.Context, network, address); err != nil {
			err = errors.Wrap(err)
			return
		}

		addr := listener.Addr().(*net.TCPAddr)

		ui.Log().Printf("starting HTTP server on port: %q", strconv.Itoa(addr.Port))

	default:
		if listener, err = config.Listen(env.Context, network, address); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

type UnixSocket struct {
	net.Listener
	Path string
}

func (env *Repo) InitializeUnixSocket(
	config net.ListenConfig,
	path string,
) (sock UnixSocket, err error) {
	sock.Path = path

	if sock.Path == "" {
		dir := env.GetRepoLayout().GetXDG().State

		if err = os.MkdirAll(dir, 0o700); err != nil {
			err = errors.Wrap(err)
			return
		}

		sock.Path = fmt.Sprintf("%s/%d.sock", dir, os.Getpid())
	}

	ui.Log().Printf("starting unix domain server on socket: %q", sock.Path)

	if sock.Listener, err = config.Listen(env.Context, "unix", sock.Path); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type HTTPPort struct {
	net.Listener
	Port int
}

func (env *Repo) InitializeHTTP(
	config net.ListenConfig,
	port int,
) (httpPort HTTPPort, err error) {
	if httpPort.Listener, err = config.Listen(
		env.Context,
		"tcp",
		fmt.Sprintf(":%d", port),
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	addr := httpPort.Addr().(*net.TCPAddr)

	ui.Log().Printf("starting HTTP server on port: %q", strconv.Itoa(addr.Port))

	return
}

func (env *Repo) Serve(listener net.Listener) (err error) {
	httpServer := http.Server{Handler: env}

	go func() {
		<-env.Done()
		ui.Log().Print("shutting down")

		ctx, cancel := context.WithTimeoutCause(
			context.Background(),
			1e9, // 1 second
			errors.Errorf("shut down timeout"),
		)

		defer cancel()

		httpServer.Shutdown(ctx)
	}()

	if err = httpServer.Serve(listener); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			err = nil
		} else {
			err = errors.Wrap(err)
			return
		}
	}

	ui.Log().Print("shutdown complete")

	return
}

func (repo *Repo) ServeStdio() (err error) {
	// shuts down the server when the main context is complete (on SIGHUP / SIGINT).
	repo.After(repo.GetIn().GetFile().Close)
	repo.After(repo.GetOut().GetFile().Close)

	br := bufio.NewReader(repo.GetIn().GetFile())
	bw := bufio.NewWriter(repo.GetOut().GetFile())

	for {
		repo.ContinueOrPanicOnDone()

		var request *http.Request

		if request, err = http.ReadRequest(br); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		response := repo.ServeRequest(
			Request{
				MethodPath: MethodPath{
					Method: request.Method,
					Path:   request.URL.Path,
				},
				Body: request.Body,
			},
		)

		if response.StatusCode == 0 {
			response.StatusCode = http.StatusOK
		}

		responseModified := &http.Response{
			// ContentLength:    -1,
			TransferEncoding: []string{"chunked"},
			ProtoMajor:       request.ProtoMajor,
			ProtoMinor:       request.ProtoMinor,
			Request:          request,
			StatusCode:       response.StatusCode,
			Body:             response.Body,
		}

		if err = responseModified.Write(bw); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		if err = bw.Flush(); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}
	}

	ui.Log().Print("shutdown complete")

	return
}

type MethodPath struct {
	Method string
	Path   string
}

type Request struct {
	MethodPath
	Body io.ReadCloser
}

type Response struct {
	StatusCode int
	Body       io.ReadCloser
}

func (r *Response) ErrorWithStatus(status int, err error) {
	r.StatusCode = status
	r.Body = io.NopCloser(strings.NewReader(err.Error()))
}

func (r *Response) Error(err error) {
	r.ErrorWithStatus(http.StatusInternalServerError, err)
}

func (local *Repo) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request := Request{
		MethodPath: MethodPath{Method: req.Method, Path: req.URL.Path},
		Body:       req.Body,
	}

	response := local.ServeRequest(request)

	if response.StatusCode == 0 {
		response.StatusCode = http.StatusOK
	}

	w.WriteHeader(response.StatusCode)

	if _, err := io.Copy(w, response.Body); err != nil {
		local.CancelWithError(err)
	}

	if err := response.Body.Close(); err != nil {
		local.CancelWithError(err)
	}
}

// TODO add path multiplexing to handle versions
func (local *Repo) ServeRequest(request Request) (response Response) {
	defer local.ContinueOrPanicOnDone()

	ui.Log().Printf("serving: %s %s", request.Method, request.Path)

	switch request.MethodPath {
	case MethodPath{"HEAD", "/blobs"}, MethodPath{"GET", "/blobs"}:
		var shString strings.Builder

		if _, err := io.Copy(&shString, request.Body); err != nil {
			response.ErrorWithStatus(http.StatusBadRequest, err)
			return
		}

		var sh *sha.Sha

		{
			var err error

			if sh, err = sha.MakeSha(shString.String()); err != nil {
				response.ErrorWithStatus(http.StatusBadRequest, err)
				return
			}
		}

		ui.Log().Printf("blob requested: %q", sh)

		if request.Method == "HEAD" {
			if local.GetRepoLayout().HasBlob(sh) {
				response.StatusCode = http.StatusNoContent
			} else {
				response.StatusCode = http.StatusNotFound
			}
		} else {
			var rc sha.ReadCloser

			{
				var err error

				if rc, err = local.GetRepoLayout().BlobReader(sh); err != nil {
					response.Error(err)
					return
				}
			}

			response.Body = rc
		}

	case MethodPath{"POST", "/blobs"}:
		defer func() {
			if r := recover(); r != nil {
				local.GetUI().Printf("panicked: %s", r)
				panic(r)
			}
		}()
		var wc interfaces.ShaWriteCloser

		{
			var err error

			if wc, err = local.GetRepoLayout().BlobWriter(); err != nil {
				response.Error(err)
				return
			}
		}

		local.GetUI().Print("made blob writer")

		var n int64

		{
			var err error

			if n, err = io.Copy(wc, request.Body); err != nil {
				response.Error(err)
				return
			}
		}

		local.GetUI().Printf("copied %d bytes to blob writer", n)

		if err := wc.Close(); err != nil {
			response.Error(err)
			return
		}
		local.GetUI().Printf("closed writer")

		sh := wc.GetShaLike()
		local.GetUI().Printf("got sha: %s", sh)

		blobCopierDelegate := local.MakeBlobCopierDelegate()

		if err := blobCopierDelegate(
			store.BlobCopyResult{
				Sha: sh,
				N:   n,
			},
		); err != nil {
			response.Error(err)
			return
		}

		response.StatusCode = http.StatusCreated
		response.Body = io.NopCloser(strings.NewReader(sh.GetShaString()))

		// 	case MethodPath{"GET", "/object"}:

		// 	case MethodPath{"POST", "/object"}:

	case MethodPath{"GET", "/inventory_lists"}:
		var qgString strings.Builder

		if _, err := io.Copy(&qgString, request.Body); err != nil {
			response.Error(err)
			return
		}

		var qg *query.Group

		{
			var err error

			if qg, err = local.MakeQueryGroup(
				nil,
				ids.RepoId{},
				sku.ExternalQueryOptions{},
				qgString.String(),
			); err != nil {
				response.Error(err)
				return
			}
		}

		var list *sku.List

		{
			var err error

			if list, err = local.MakeInventoryList(qg); err != nil {
				response.Error(err)
				return
			}
		}

		// TODO make this more performant by returning a proper reader
		b := bytes.NewBuffer(nil)

		printer := local.MakePrinterBoxArchive(b, true)

		var sk *sku.Transacted
		var hasMore bool

		for {
			local.ContinueOrPanicOnDone()

			sk, hasMore = list.Pop()

			if !hasMore {
				break
			}

			if err := printer(sk); err != nil {
				response.Error(err)
				return
			}
		}

		response.Body = io.NopCloser(b)

	case MethodPath{"POST", "/inventory_lists"}:
		bf := local.GetStore().GetInventoryListStore().FormatForVersion(
			local.GetConfig().GetStoreVersion(),
		)

		list := sku.MakeList()

		if err := inventory_list_blobs.ReadInventoryListBlob(
			bf,
			bufio.NewReader(request.Body),
			list,
		); err != nil {
			response.Error(err)
			return
		}

		b := bytes.NewBuffer(nil)

		importer := local.MakeImporter(false)
		importer.BlobCopierDelegate = func(result store.BlobCopyResult) (err error) {
			local.ContinueOrPanicOnDone()

			if result.N != -1 {
				return
			}

			sh := sha.GetPool().Get()
			sha.GetPool().Put(sh)
			sh.ResetWithShaLike(result.GetBlobSha())
			fmt.Fprintf(b, "%s\n", sh)

			return
		}

		if err := local.ImportList(
			list,
			importer,
		); err != nil {
			response.Error(err)
			return
		}

		response.StatusCode = http.StatusCreated

		if b.Len() > 0 {
			response.Body = io.NopCloser(b)
		}

	default:
		response.StatusCode = http.StatusNotFound
	}

	return
}
