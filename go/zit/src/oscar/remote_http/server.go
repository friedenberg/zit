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
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"github.com/gorilla/mux"
)

type Server struct {
	EnvLocal env_local.Env
	Repo     repo.Repo
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
		dir := server.EnvLocal.GetXDG().State

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

func (server Server) makeRouter(
	makeHandler func(handler funcHandler) http.HandlerFunc,
) http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/blobs/{sha}", makeHandler(server.handleBlobsHeadOrGet)).
		Methods("HEAD", "GET")

	router.HandleFunc("/blobs", makeHandler(server.handleBlobsPost)).
		Methods("POST")

	router.HandleFunc("/inventory_lists", makeHandler(server.handleGetInventoryList)).
		Methods("GET")

	router.HandleFunc("/inventory_lists", makeHandler(server.handlePostInventoryList)).
		Methods("POST")

	return router
}

func (server Server) Serve(listener net.Listener) (err error) {
	httpServer := http.Server{
		Handler: server.makeRouter(server.makeHandler),
	}

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

func (server Server) ServeStdio() {
	// shuts down the server when the main context is complete (on SIGHUP / SIGINT).
	server.Repo.GetEnv().After(server.Repo.GetEnv().GetIn().GetFile().Close)
	server.Repo.GetEnv().After(server.Repo.GetEnv().GetOut().GetFile().Close)

	br := bufio.NewReader(server.Repo.GetEnv().GetIn().GetFile())
	bw := bufio.NewWriter(server.Repo.GetEnv().GetOut().GetFile())

	handler := server.makeRouter(
		func(handler funcHandler) http.HandlerFunc {
			return server.makeHandlerWithRedirect(handler, bw)
		},
	)

	for {
		server.Repo.GetEnv().ContinueOrPanicOnDone()

		var request *http.Request

		{
			var err error

			if request, err = http.ReadRequest(br); err != nil {
				if errors.IsEOF(err) {
					err = nil
				} else {
					server.EnvLocal.CancelWithError(err)
				}

				return
			}
		}

		handler.ServeHTTP(nil, request)

		if err := bw.Flush(); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				server.EnvLocal.CancelWithError(err)
			}

			return
		}
	}
}

type funcHandler func(Request) Response

type handlerWrapper funcHandler

func (server *Server) makeHandlerWithRedirect(
	handler funcHandler,
	out *bufio.Writer,
) http.HandlerFunc {
	return func(_ http.ResponseWriter, req *http.Request) {
		request := Request{
			request:    req,
			MethodPath: MethodPath{Method: req.Method, Path: req.URL.Path},
			Headers:    req.Header,
			Body:       req.Body,
		}

		response := handler(request)

		if response.StatusCode == 0 {
			response.StatusCode = http.StatusOK
		}

		if response.StatusCode == 0 {
			response.StatusCode = http.StatusOK
		}

		responseModified := &http.Response{
			// ContentLength:    -1,
			TransferEncoding: []string{"chunked"},
			ProtoMajor:       req.ProtoMajor,
			ProtoMinor:       req.ProtoMinor,
			Request:          req,
			StatusCode:       response.StatusCode,
			Body:             response.Body,
		}

		if err := responseModified.Write(out); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				server.EnvLocal.CancelWithError(err)
			}
		}
	}
}

func (server *Server) makeHandler(
	handler funcHandler,
) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, req *http.Request) {
		request := Request{
			request:    req,
			MethodPath: MethodPath{Method: req.Method, Path: req.URL.Path},
			Headers:    req.Header,
			Body:       req.Body,
		}

		response := handler(request)

		if response.StatusCode == 0 {
			response.StatusCode = http.StatusOK
		}

		if response.StatusCode == 0 {
			response.StatusCode = http.StatusOK
		}

		responseWriter.WriteHeader(response.StatusCode)

		if _, err := io.Copy(responseWriter, response.Body); err != nil {
			if errors.IsEOF(err) {
				err = nil
			} else {
				server.EnvLocal.CancelWithError(err)
			}
		}
	}
}

func (server *Server) handleBlobsHeadOrGet(request Request) (response Response) {
	shString := request.Vars()["sha"]

	if shString == "" {
		response.ErrorWithStatus(http.StatusBadRequest, errors.Errorf("empty sha"))
		return
	}

	var sh *sha.Sha

	{
		var err error

		if sh, err = sha.MakeSha(shString); err != nil {
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

	return
}

func (server *Server) handleBlobsPost(request Request) (response Response) {
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

	return
}

func (server *Server) handleGetInventoryList(request Request) (response Response) {
	if repo, ok := server.Repo.(*local_working_copy.Repo); ok {
		var qgString strings.Builder

		if _, err := io.Copy(&qgString, request.Body); err != nil {
			response.Error(err)
			return
		}

		var qg *query.Group

		{
			var err error

			if qg, err = repo.MakeExternalQueryGroup(
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

			if list, err = repo.MakeInventoryList(qg); err != nil {
				response.Error(err)
				return
			}
		}

		// TODO make this more performant by returning a proper reader
		b := bytes.NewBuffer(nil)

		// TODO
		printer := repo.MakePrinterBoxArchive(b, true)

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
	} else {
	}

	return
}

func (server *Server) handlePostInventoryList(request Request) (response Response) {
	if repo, ok := server.Repo.(*local_working_copy.Repo); ok {
		response = server.writeInventoryListLocalWorkingCopy(repo, request)
	} else {
		response = server.writeInventoryList(request)
	}

	return
}
