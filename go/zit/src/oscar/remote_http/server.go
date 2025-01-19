package remote_http

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
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/inventory_list_blobs"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/mike/store"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type Server struct {
	// Repo repo.LocalWorkingCopy
	Repo *local_working_copy.Repo
}

func (server Server) InitializeListener(
	network, address string,
) (listener net.Listener, err error) {
	var config net.ListenConfig

	switch network {
	case "unix":
		if listener, err = server.InitializeUnixSocket(config, address); err != nil {
			err = errors.Wrap(err)
			return
		}

	case "tcp":
		if listener, err = config.Listen(
			server.Repo.GetEnv(),
			network,
			address,
		); err != nil {
			err = errors.Wrap(err)
			return
		}

		addr := listener.Addr().(*net.TCPAddr)

		ui.Log().Printf("starting HTTP server on port: %q", strconv.Itoa(addr.Port))

	default:
		if listener, err = config.Listen(
			server.Repo.GetEnv(),
			network,
			address,
		); err != nil {
			err = errors.Wrap(err)
			return
		}
	}

	return
}

func (server Server) InitializeUnixSocket(
	config net.ListenConfig,
	path string,
) (sock repo.UnixSocket, err error) {
	sock.Path = path

	if sock.Path == "" {
		dir := server.Repo.GetRepoLayout().GetXDG().State

		if err = os.MkdirAll(dir, 0o700); err != nil {
			err = errors.Wrap(err)
			return
		}

		sock.Path = fmt.Sprintf("%s/%d.sock", dir, os.Getpid())
	}

	ui.Log().Printf("starting unix domain server on socket: %q", sock.Path)

	if sock.Listener, err = config.Listen(
		server.Repo.GetEnv(),
		"unix",
		sock.Path,
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

type HTTPPort struct {
	net.Listener
	Port int
}

func (server Server) InitializeHTTP(
	config net.ListenConfig,
	port int,
) (httpPort HTTPPort, err error) {
	if httpPort.Listener, err = config.Listen(
		server.Repo.GetEnv(),
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

func (server Server) Serve(listener net.Listener) (err error) {
	httpServer := http.Server{Handler: server}

	go func() {
		<-server.Repo.GetEnv().Done()
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

func (server Server) ServeStdio() (err error) {
	// shuts down the server when the main context is complete (on SIGHUP / SIGINT).
	server.Repo.GetEnv().After(server.Repo.GetEnv().GetIn().GetFile().Close)
	server.Repo.GetEnv().After(server.Repo.GetEnv().GetOut().GetFile().Close)

	br := bufio.NewReader(server.Repo.GetEnv().GetIn().GetFile())
	bw := bufio.NewWriter(server.Repo.GetEnv().GetOut().GetFile())

	for {
		server.Repo.GetEnv().ContinueOrPanicOnDone()

		var request *http.Request

		if request, err = http.ReadRequest(br); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				err = errors.Wrap(err)
			}

			return
		}

		response := server.ServeRequest(
			Request{
				MethodPath: MethodPath{
					Method: request.Method,
					Path:   request.URL.Path,
				},
				Headers: request.Header,
				Body:    request.Body,
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
	Headers http.Header
	Body    io.ReadCloser
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

func (server Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	request := Request{
		MethodPath: MethodPath{Method: req.Method, Path: req.URL.Path},
		Headers:    req.Header,
		Body:       req.Body,
	}

	response := server.ServeRequest(request)

	if response.StatusCode == 0 {
		response.StatusCode = http.StatusOK
	}

	w.WriteHeader(response.StatusCode)

	if _, err := io.Copy(w, response.Body); err != nil {
		server.Repo.GetEnv().CancelWithError(err)
	}

	if err := response.Body.Close(); err != nil {
		server.Repo.GetEnv().CancelWithError(err)
	}
}

// TODO add path multiplexing to handle versions
func (server Server) ServeRequest(request Request) (response Response) {
	defer server.Repo.GetEnv().ContinueOrPanicOnDone()

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
			if server.Repo.GetBlobStore().HasBlob(sh) {
				response.StatusCode = http.StatusNoContent
			} else {
				response.StatusCode = http.StatusNotFound
			}
		} else {
			var rc sha.ReadCloser

			{
				var err error

				if rc, err = server.Repo.GetBlobStore().BlobReader(sh); err != nil {
					response.Error(err)
					return
				}
			}

			response.Body = rc
		}

	case MethodPath{"POST", "/blobs"}:
		var wc interfaces.ShaWriteCloser

		{
			var err error

			if wc, err = server.Repo.GetBlobStore().BlobWriter(); err != nil {
				response.Error(err)
				return
			}
		}

		var n int64

		{
			var err error

			if n, err = io.Copy(wc, request.Body); err != nil {
				response.Error(err)
				return
			}
		}

		if err := wc.Close(); err != nil {
			response.Error(err)
			return
		}

		sh := wc.GetShaLike()

		blobCopierDelegate := sku.MakeBlobCopierDelegate(
			server.Repo.GetEnv().GetUI(),
		)

		if err := blobCopierDelegate(
			sku.BlobCopyResult{
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

			if qg, err = server.Repo.MakeExternalQueryGroup(
				query.BuilderOptions{},
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

			if list, err = server.Repo.MakeInventoryList(qg); err != nil {
				response.Error(err)
				return
			}
		}

		// TODO make this more performant by returning a proper reader
		b := bytes.NewBuffer(nil)

		// TODO
		printer := server.Repo.MakePrinterBoxArchive(b, true)

		var sk *sku.Transacted
		var hasMore bool

		for {
			server.Repo.GetEnv().ContinueOrPanicOnDone()

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
		// TODO get version from header?
		// TODO
		bf := server.Repo.GetStore().GetInventoryListStore().FormatForVersion(
			server.Repo.GetStoreVersion(),
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

		// TODO make option to read from headers
		importerOptions := store.ImporterOptions{
			// TODO
			CheckedOutPrinter: server.Repo.PrinterCheckedOutConflictsForRemoteTransfers(),
		}

		if request.Headers.Get("x-zit-remote_transfer_options-allow_merge_conflicts") == "true" {
			importerOptions.AllowMergeConflicts = true
		}

		importerOptions.BlobCopierDelegate = func(
			result sku.BlobCopyResult,
		) (err error) {
			server.Repo.GetEnv().ContinueOrPanicOnDone()

			if result.N != -1 {
				return
			}

			sh := sha.GetPool().Get()
			sha.GetPool().Put(sh)
			sh.ResetWithShaLike(result.GetBlobSha())
			fmt.Fprintf(b, "%s\n", sh)

			return
		}

		// TODO
		importer := server.Repo.MakeImporter(
			importerOptions,
			sku.GetStoreOptionsRemoteTransfer(),
		)

		// TODO
		if err := server.Repo.ImportList(
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
