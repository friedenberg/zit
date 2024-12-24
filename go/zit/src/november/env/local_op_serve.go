package env

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/bravo/ui"
	"code.linenisgreat.com/zit/go/zit/src/delta/sha"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/hotel/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
)

func (env *Local) InitializeListener(
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
		if _, err = strconv.Atoi(address); err != nil {
			err = errors.Wrap(err)
			return
		}

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

func (env *Local) InitializeUnixSocket(
	config net.ListenConfig,
	path string,
) (sock UnixSocket, err error) {
	sock.Path = path

	if sock.Path == "" {
		dir := env.GetDirectoryLayout().GetXDG().State

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

func (env *Local) InitializeHTTP(
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

func (env *Local) Serve(listener net.Listener) (err error) {
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

func (local *Local) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	type MethodPath struct {
		Method string
		Path   string
	}

	mp := MethodPath{Method: req.Method, Path: req.URL.Path}
	ui.Log().Printf("serving: %s %s", mp.Method, mp.Path)

	switch mp {
	case MethodPath{"HEAD", "/blobs"}, MethodPath{"GET", "/blobs"}:
		var shString strings.Builder

		if _, err := io.Copy(&shString, req.Body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		}

		var sh *sha.Sha

		{

			var err error

			if sh, err = sha.MakeSha(shString.String()); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				io.WriteString(w, err.Error())
				return
			}
		}

		ui.Log().Printf("blob requested: %q", sh)

		if mp.Method == "HEAD" {
			if local.GetDirectoryLayout().HasBlob(sh) {
				w.WriteHeader(http.StatusNoContent)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		} else {
			var rc sha.ReadCloser

			{
				var err error

				if rc, err = local.GetDirectoryLayout().BlobReader(sh); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					io.WriteString(w, err.Error())
					return
				}
			}

			if _, err := io.Copy(w, rc); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, err.Error())
				return
			}

			if err := rc.Close(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, err.Error())
				return
			}
		}

	case MethodPath{"GET", "/inventory_list"}:
		var qgString strings.Builder

		if _, err := io.Copy(&qgString, req.Body); err != nil {
			panic(err)
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
				panic(err)
			}
		}

		var list *sku.List

		{
			var err error

			if list, err = local.MakeInventoryList(qg); err != nil {
				panic(err)
			}
		}

		bw := bufio.NewWriter(w)

		printer := local.MakePrinterBoxArchive(bw, true)

		var sk *sku.Transacted
		var hasMore bool

		for {
			sk, hasMore = list.Pop()

			if !hasMore {
				break
			}

			if err := printer(sk); err != nil {
				panic(err)
			}
		}

		if err := bw.Flush(); err != nil {
			panic(err)
		}

	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
