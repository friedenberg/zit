package remote_http

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"syscall"
	"time"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/charlie/repo_signing"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/delta/string_format_writer"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/golf/config_immutable_io"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/box_format"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"github.com/gorilla/mux"
)

type Server struct {
	EnvLocal  env_local.Env
	Repo      repo.LocalRepo
	blobCache serverBlobCache

	GetCertificate func(*tls.ClientHelloInfo) (*tls.Certificate, error)
}

func (server *Server) init() (err error) {
	server.blobCache.localBlobStore = server.Repo.GetEnvRepo().GetLocalBlobStore()
	server.blobCache.ui = server.Repo.GetEnv().GetUI()
	return
}

// TODO switch to not return error
func (server *Server) InitializeListener(
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

		server.EnvLocal.GetOut().Printf(
			"starting HTTP server on port: %q",
			strconv.Itoa(addr.Port),
		)

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

func (server *Server) InitializeUnixSocket(
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

func (server *Server) InitializeHTTP(
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

func (server *Server) makeRouter(
	makeHandler func(handler funcHandler) http.HandlerFunc,
) http.Handler {
	// TODO add errors/context middlerware for capturing errors and panics
	router := mux.NewRouter().UseEncodedPath()

	router.HandleFunc("/config-immutable", makeHandler(server.handleGetConfigImmutable)).
		Methods("GET")

	{
		router.HandleFunc("/blobs/{sha}", makeHandler(server.handleBlobsHeadOrGet)).
			Methods("HEAD", "GET")

		router.HandleFunc("/blobs/{sha}", makeHandler(server.handleBlobsPost)).
			Methods("POST")

		router.HandleFunc("/blobs", makeHandler(server.handleBlobsPost)).
			Methods("POST")
	}

	router.HandleFunc("/query/{query}", makeHandler(server.handleGetQuery)).
		Methods("GET")

	{
		router.HandleFunc("/inventory_lists", makeHandler(server.handleGetInventoryList)).
			Methods("GET")

		router.HandleFunc("/inventory_lists", makeHandler(server.handlePostInventoryList)).
			Methods("POST")

		router.HandleFunc("/inventory_lists/{box}", makeHandler(server.handlePostInventoryList)).
			Methods("POST")
	}

	if server.Repo.GetEnv().GetCLIConfig().Verbose {
		router.Use(server.loggerMiddleware)
	}

	router.Use(server.panicHandlingMiddleware)

	if len(server.Repo.GetImmutableConfigPrivate().ImmutableConfig.GetPrivateKey()) > 0 {
		router.Use(server.sigMiddleware)
	}

	return router
}

func (server *Server) sigMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) {
			nonceStringBase64 := request.Header.Get(headerChallengeNonce)

			var nonce []byte

			{
				var err error

				if nonce, err = base64.URLEncoding.DecodeString(
					nonceStringBase64,
				); err != nil {
					http.Error(responseWriter, err.Error(), http.StatusBadRequest)
					return
				}
			}

			if len(nonce) > 0 {
				key := server.Repo.GetImmutableConfigPrivate().ImmutableConfig.GetPrivateKey()

				var sig string

				{
					var err error

					if sig, err = repo_signing.SignBase64(key, nonce); err != nil {
						server.EnvLocal.CancelWithError(err)
					}
				}

				responseWriter.Header().Set(headerChallengeResponse, sig)
			}

			next.ServeHTTP(responseWriter, request)
		},
	)
}

func (server *Server) loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) {
			ui.Log().Printf("serving request: %s %s", request.Method, request.URL.Path)
			next.ServeHTTP(responseWriter, request)
			ui.Log().Printf("done serving request: %s %s", request.Method, request.URL.Path)
		},
	)
}

func (server *Server) panicHandlingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(responseWriter http.ResponseWriter, request *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					ui.Log().Print("request handler panicked", request.URL)

					switch err := r.(type) {
					default:
						panic(err)

					case error:
						http.Error(
							responseWriter,
							fmt.Sprintf("%s: %s", err, debug.Stack()),
							http.StatusInternalServerError,
						)
					}
				}
			}()

			next.ServeHTTP(responseWriter, request)
		},
	)
}

// TODO remove error return and use context
func (server *Server) Serve(listener net.Listener) (err error) {
	if err = server.init(); err != nil {
		err = errors.Wrap(err)
		return
	}

	httpServer := http.Server{
		Handler: server.makeRouter(server.makeHandler),
	}

	if server.GetCertificate != nil {
		httpServer.TLSConfig = &tls.Config{
			GetCertificate: server.GetCertificate,
		}
	}

	go func() {
		<-server.Repo.GetEnv().Done()
		ui.Log().Print("shutting down")

		ctx, cancel := context.WithTimeoutCause(
			context.Background(),
			1e9, // 1 second
			errors.ErrorWithStackf("shut down timeout"),
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

func (server *Server) ServeStdio() {
	if err := server.init(); err != nil {
		server.EnvLocal.CancelWithError(err)
		return
	}

	// shuts down the server when the main context is complete (on SIGHUP / SIGINT).
	server.Repo.GetEnv().After(server.Repo.GetEnv().GetIn().GetFile().Close)
	server.Repo.GetEnv().After(server.Repo.GetEnv().GetOut().GetFile().Close)

	bufferedReader := bufio.NewReader(server.Repo.GetEnv().GetIn().GetFile())
	bufferedWriter := bufio.NewWriter(server.Repo.GetEnv().GetOut().GetFile())

	var responseWriter BufferedResponseWriter

	handler := server.makeRouter(
		func(handler funcHandler) http.HandlerFunc {
			return server.makeHandlerUsingBufferedWriter(handler, bufferedWriter)
		},
	)

	for {
		server.Repo.GetEnv().ContinueOrPanicOnDone()
		responseWriter.Reset()

		var request *http.Request

		{
			var err error

			if request, err = http.ReadRequest(bufferedReader); err != nil {
				if errors.IsEOF(err) {
					err = nil
				} else {
					server.EnvLocal.CancelWithError(err)
				}

				return
			}
		}

		handler.ServeHTTP(&responseWriter, request)

		if err := request.Body.Close(); err != nil {
			server.EnvLocal.CancelWithError(err)
			return
		}

		if responseWriter.Dirty {
			if err := responseWriter.WriteResponse(bufferedWriter); err != nil {
				server.EnvLocal.CancelWithError(err)
				return
			}
		}

		if err := bufferedWriter.Flush(); err != nil {
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

// TODO switch to using responseWriter
func (server *Server) makeHandlerUsingBufferedWriter(
	handler funcHandler,
	out *bufio.Writer,
) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, req *http.Request) {
		request := Request{
			context:    errors.MakeContext(server.EnvLocal),
			request:    req,
			MethodPath: MethodPath{Method: req.Method, Path: req.URL.Path},
			Headers:    req.Header,
			Body:       req.Body,
		}

		var response Response

		if err := errors.RunContextWithPrintTicker(
			request.context,
			func(ctx errors.Context) {
				response = handler(request)

				responseModified := http.Response{
					TransferEncoding: []string{"chunked"},
					ProtoMajor:       req.ProtoMajor,
					ProtoMinor:       req.ProtoMinor,
					Request:          req,
					StatusCode:       response.StatusCode,
					Body:             response.Body,
				}

				// TODO determine why iterating thru the headers and setting them manually
				// doesn't work
				responseModified.Header = responseWriter.Header().Clone()

				if responseModified.StatusCode == 0 {
					responseModified.StatusCode = http.StatusOK
				}

				if err := responseModified.Write(out); err != nil {
					if errors.IsEOF(err) {
						err = nil
					} else {
						ctx.CancelWithError(err)
					}
				}
			},
			func(time time.Time) {
				ui.Log().Printf("Still serving request (%s): %s", time, req.URL)
			},
			3*time.Second,
		); err != nil {
			server.EnvLocal.CancelWithError(err)
		}
	}
}

func (server *Server) makeHandler(
	handler funcHandler,
) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, req *http.Request) {
		request := Request{
			context:    errors.MakeContext(server.EnvLocal),
			request:    req,
			MethodPath: MethodPath{Method: req.Method, Path: req.URL.Path},
			Headers:    req.Header,
			Body:       req.Body,
		}

		var progressWriter env_ui.ProgressWriter

		if err := errors.RunContextWithPrintTicker(
			request.context,
			func(ctx errors.Context) {
				response := handler(request)

				// header := responseWriter.Header()

				// for key, values := range response.Headers {
				// 	for _, value := range values {
				// 		header.Add(key, value)
				// 	}
				// }

				if response.StatusCode == 0 {
					response.StatusCode = http.StatusOK
				}

				responseWriter.WriteHeader(response.StatusCode)

				if response.Body == nil {
					return
				}

				if _, err := io.Copy(
					io.MultiWriter(responseWriter, &progressWriter),
					response.Body,
				); err != nil {
					if errors.IsEOF(err) {
						err = nil
					} else if errors.IsAny(
						err,
						errors.MakeIsErrno(
							syscall.ECONNRESET,
							syscall.EPIPE,
						),
						errors.IsNetTimeout,
					) {
						ui.Err().Print(errors.Unwrap(err).Error(), req.URL)
						err = nil
					} else {
						ctx.CancelWithError(err)
					}
				}
			},
			func(time time.Time) {
				ui.Log().Printf(
					"Still serving request (%s): %q (%s bytes written)",
					time,
					req.URL,
					progressWriter.GetWrittenHumanString(),
				)
			},
			3*time.Second,
		); err != nil {
			server.EnvLocal.CancelWithError(err)
		}
	}
}

func (server *Server) handleBlobsHeadOrGet(request Request) (response Response) {
	shString := request.Vars()["sha"]

	if shString == "" {
		response.ErrorWithStatus(http.StatusBadRequest, errors.ErrorWithStackf("empty sha"))
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
				if env_dir.IsErrBlobMissing(err) {
					response.StatusCode = http.StatusNotFound
				} else {
					response.Error(err)
				}

				return
			}
		}

		response.Body = rc
	}

	return
}

func (server *Server) handleBlobsPost(request Request) (response Response) {
	shString := request.Vars()["sha"]
	var result interfaces.Sha

	if shString == "" {
		var err error

		if result, err = server.copyBlob(request.Body); err != nil {
			response.Error(err)
			return
		}

		response.StatusCode = http.StatusCreated
		response.Body = io.NopCloser(strings.NewReader(result.GetShaString()))

		return
	}

	var sh sha.Sha

	if err := sh.Set(shString); err != nil {
		response.Error(err)
		return
	}

	if server.Repo.GetBlobStore().HasBlob(&sh) {
		response.StatusCode = http.StatusFound
		return
	}

	{
		var err error

		if result, err = server.copyBlob(request.Body); err != nil {
			response.Error(err)
			return
		}
	}

	response.StatusCode = http.StatusCreated

	if err := sh.AssertEqualsShaLike(result); err != nil {
		response.Error(err)
		return
	}

	response.StatusCode = http.StatusCreated
	response.Body = io.NopCloser(strings.NewReader(result.GetShaString()))

	return
}

func (server *Server) copyBlob(
	reader io.ReadCloser,
) (result interfaces.Sha, err error) {
	var writeCloser interfaces.ShaWriteCloser

	if writeCloser, err = server.Repo.GetBlobStore().BlobWriter(); err != nil {
		err = errors.Wrap(err)
		return
	}

	var n int64

	if n, err = io.Copy(writeCloser, reader); err != nil {
		err = errors.Wrap(err)
		return
	}

	if err = writeCloser.Close(); err != nil {
		err = errors.Wrap(err)
		return
	}

	result = writeCloser.GetShaLike()

	blobCopierDelegate := sku.MakeBlobCopierDelegate(
		server.Repo.GetEnv().GetUI(),
	)

	if err = blobCopierDelegate(
		sku.BlobCopyResult{
			Sha: result,
			N:   n,
		},
	); err != nil {
		err = errors.Wrap(err)
		return
	}

	return
}

func (server *Server) handleGetQuery(request Request) (response Response) {
	var queryGroupString string

	{
		var err error

		if queryGroupString, err = url.QueryUnescape(
			request.Vars()["query"],
		); err != nil {
			response.Error(err)
			return
		}
	}

	if repo, ok := server.Repo.(*local_working_copy.Repo); ok {
		var queryGroup *query.Query

		{
			var err error

			if queryGroup, err = repo.MakeExternalQueryGroup(
				nil,
				sku.ExternalQueryOptions{},
				queryGroupString,
			); err != nil {
				response.Error(err)
				return
			}
		}

		var list *sku.List

		{
			var err error

			if list, err = repo.MakeInventoryList(queryGroup); err != nil {
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
		response.StatusCode = http.StatusNotImplemented
	}

	return
}

func (server *Server) handleGetInventoryList(
	request Request,
) (response Response) {
	inventoryListStore := server.Repo.GetInventoryListStore()

	// TODO make this more performant by returning a proper reader
	b := bytes.NewBuffer(nil)

	boxFormat := box_format.MakeBoxTransactedArchive(
		server.Repo.GetEnv(),
		server.Repo.GetEnv().GetCLIConfig().PrintOptions.WithPrintTai(true),
	)

	printer := string_format_writer.MakeDelim(
		"\n",
		b,
		string_format_writer.MakeFunc(
			func(
				writer interfaces.WriterAndStringWriter,
				object *sku.Transacted,
			) (n int64, err error) {
				return boxFormat.EncodeStringTo(object, writer)
			},
		),
	)

	iter := inventoryListStore.IterAllInventoryLists()

	for sk, err := range iter {
		if err != nil {
			response.Error(err)
			return
		}

		server.Repo.GetEnv().ContinueOrPanicOnDone()

		if err = printer(sk); err != nil {
			response.Error(err)
			return
		}
	}

	response.Body = io.NopCloser(b)

	return
}

func (server *Server) handlePostInventoryList(
	request Request,
) (response Response) {
	boxString := request.Vars()["box"]

	var sk *sku.Transacted

	typedInventoryListStore := server.Repo.GetTypedInventoryListBlobStore()

	if boxString != "" {

		{
			var err error

			if boxString, err = url.QueryUnescape(request.Vars()["box"]); err != nil {
				response.Error(err)
				return
			}
		}

		{
			var err error

			if sk, err = typedInventoryListStore.ReadInventoryListObject(
				ids.MustType(builtin_types.InventoryListTypeV1),
				strings.NewReader(boxString),
			); err != nil {
				response.Error(
					errors.ErrorWithStackf(
						"failed to parse inventory list sku (%q): %w",
						boxString,
						err,
					),
				)

				return
			}
		}

		defer sku.GetTransactedPool().Put(sk)
	}

	// TODO parse box into sk
	if repo, ok := server.Repo.(*local_working_copy.Repo); ok {
		response = server.writeInventoryListLocalWorkingCopy(repo, request, sk)
	} else {
		response = server.writeInventoryList(request, sk)
	}

	return
}

func (server *Server) handleGetConfigImmutable(request Request) (response Response) {
	config := server.Repo.GetImmutableConfigPublic()
	configLoaded := &config_immutable_io.ConfigLoadedPublic{
		Type:            config.Type,
		ImmutableConfig: config.ImmutableConfig,
	}

	encoder := config_immutable_io.CoderPublic{}

	var b bytes.Buffer

	// TODO modify to not have to buffer
	if _, err := encoder.EncodeTo(configLoaded, &b); err != nil {
		server.EnvLocal.CancelWithError(err)
	}

	response.Body = io.NopCloser(&b)

	return
}
