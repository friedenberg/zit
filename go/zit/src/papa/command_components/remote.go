package command_components

import (
	"crypto/ed25519"
	"flag"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/xdg"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_blobs"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/lima/typed_blob_store"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/oscar/remote_http"
)

type Remote struct {
	Env
	EnvRepo
	LocalWorkingCopy
	LocalArchive

	// TODO rename to ConnectionType
	RemoteType repo.RemoteType
}

func (cmd *Remote) SetFlagSet(f *flag.FlagSet) {
	// TODO remove and replace with repo builtin type options
	f.Var(&cmd.RemoteType, "remote-type", fmt.Sprintf("%q", repo.GetAllRemoteTypes()))
}

func (cmd Remote) CreateRemoteObject(
	req command.Request,
	local repo.LocalRepo,
) (sk *sku.Transacted) {
	envRepo := cmd.MakeEnvRepo(req, false)
	typedRepoBlobStore := typed_blob_store.MakeRepoStore(envRepo)

	sk = sku.GetTransactedPool().Get()

	var blob repo_blobs.BlobMutable

	switch cmd.RemoteType {
	default:
		req.CancelWithBadRequestf("unsupported remote type: %q", cmd.RemoteType)

	case repo.RemoteTypeNativeDotenvXDG:
		xdgDotenvPath := req.PopArg("xdg-dotenv-path")

		envLocal := cmd.MakeEnvWithXDGLayoutAndOptions(
			req,
			xdgDotenvPath,
			envRepo.GetOptions(),
		)

		sk.Metadata.Type = builtin_types.GetOrPanic(builtin_types.RepoTypeXDGDotenvV0).Type
		blob = repo_blobs.TomlXDGV0FromXDG(envLocal.GetXDG())

	case repo.RemoteTypeStdioLocal:
		path := req.PopArg("path")

		sk.Metadata.Type = builtin_types.GetOrPanic(builtin_types.RepoTypeLocalPath).Type
		blob = &repo_blobs.TomlLocalPathV0{Path: envRepo.AbsFromCwdOrSame(path)}
	}

	remote := cmd.MakeRemoteFromBlob(req, local, blob.GetRepoBlob())
	remoteConfig := remote.GetImmutableConfigPublic().ImmutableConfig
	blob.SetPublicKey(remoteConfig.GetPublicKey())

	var blobSha interfaces.Sha

	{
		var err error

		if blobSha, _, err = typedRepoBlobStore.WriteTypedBlob(
			sk.Metadata.Type,
			blob,
		); err != nil {
			req.CancelWithError(err)
		}
	}

	sk.Metadata.Blob.ResetWithShaLike(blobSha)

	return
}

func (cmd Remote) MakeRemote(
	req command.Request,
	local repo.LocalRepo,
	sk *sku.Transacted,
) (remote repo.Repo) {
	envRepo := cmd.MakeEnvRepo(req, false)
	typedRepoBlobStore := typed_blob_store.MakeRepoStore(envRepo)

	var blob repo_blobs.Blob

	{
		var err error

		if blob, _, err = typedRepoBlobStore.ReadTypedBlob(
			sk.Metadata.Type,
			sk.GetBlobSha(),
		); err != nil {
			req.CancelWithError(err)
		}
	}

	return cmd.MakeRemoteFromBlob(req, local, blob)
}

func (cmd Remote) MakeRemoteFromBlob(
	req command.Request,
	local repo.LocalRepo,
	blob repo_blobs.Blob,
) (remote repo.Repo) {
	env := cmd.MakeEnv(req)

	switch blob := blob.(type) {
	case repo_blobs.TomlXDGV0:
		envDir := env_dir.MakeWithXDG(
			req,
			req.Config.Debug,
			xdg.XDG{
				Data:    blob.Data,
				Config:  blob.Config,
				Cache:   blob.Cache,
				Runtime: blob.Runtime,
				State:   blob.State,
			},
		)

		envUI := env_ui.Make(
			req,
			req.Config,
			env.GetOptions(),
		)

		remote = local_working_copy.Make(
			env_local.Make(envUI, envDir),
			local_working_copy.OptionsEmpty,
		)

	case repo_blobs.TomlLocalPathV0:
		remote = cmd.MakeRemoteStdioLocal(
			req,
			env,
			blob.Path,
			local,
			blob.GetPublicKey(),
		)

	// case repo.RemoteTypeStdioSSH:
	// 	remote = cmd.MakeRemoteStdioSSH(
	// 		req,
	// 		env,
	// 		remoteArg,
	// 		local,
	// 	)

	// case repo.RemoteTypeSocketUnix:
	// 	remote = cmd.MakeRemoteHTTPFromXDGDotenvPath(
	// 		req,
	// 		remoteArg,
	// 		env.GetOptions(),
	// 		local,
	// 	)

	case repo_blobs.TomlUriV0:
		req.CancelWithErrorf("unsupported repo blob type: %T", blob)

	default:
		req.CancelWithErrorf("unsupported repo blob type: %T", blob)
	}

	return
}

func (cmd *Remote) MakeRemoteHTTPFromXDGDotenvPath(
	req command.Request,
	xdgDotenvPath string,
	options env_ui.Options,
	localRepo repo.LocalRepo,
	pubkey ed25519.PublicKey,
) (remoteHTTP repo.WorkingCopy) {
	envLocal := cmd.MakeEnvWithXDGLayoutAndOptions(
		req,
		xdgDotenvPath,
		options,
	)

	envRepo := cmd.MakeEnvRepoFromEnvLocal(envLocal)

	remote := cmd.MakeLocalArchive(envRepo)

	server := &remote_http.Server{
		EnvLocal: envLocal,
		Repo:     remote,
	}

	var httpRoundTripper remote_http.RoundTripperUnixSocket

	if err := httpRoundTripper.Initialize(
		server,
		pubkey,
	); err != nil {
		req.CancelWithError(err)
	}

	go func() {
		if err := server.Serve(httpRoundTripper.UnixSocket); err != nil {
			req.CancelWithError(err)
		}
	}()

	remoteHTTP = remote_http.MakeClient(
		envLocal,
		&httpRoundTripper,
		localRepo,
		cmd.MakeTypedInventoryListBlobStore(envRepo),
	)

	return
}

func (cmd *Remote) MakeRemoteStdioSSH(
	req command.Request,
	env env_local.Env,
	arg string,
	local repo.LocalRepo,
) (remoteHTTP repo.WorkingCopy) {
	envRepo := cmd.MakeEnvRepo(req, false)

	var httpRoundTripper remote_http.RoundTripperStdio

	if err := httpRoundTripper.InitializeWithSSH(
		envRepo,
		arg,
	); err != nil {
		env.CancelWithError(err)
	}

	remoteHTTP = remote_http.MakeClient(
		envRepo,
		&httpRoundTripper,
		local,
		cmd.MakeTypedInventoryListBlobStore(envRepo),
	)

	return
}

func (cmd *Remote) MakeRemoteStdioLocal(
	req command.Request,
	env env_local.Env,
	dir string,
	localRepo repo.LocalRepo,
	pubkey ed25519.PublicKey,
) (remoteHTTP repo.WorkingCopy) {
	envRepo := cmd.MakeEnvRepo(req, false)

	var httpRoundTripper remote_http.RoundTripperStdio

	httpRoundTripper.Dir = dir

	if err := httpRoundTripper.InitializeWithLocal(
		envRepo,
		pubkey,
	); err != nil {
		env.CancelWithError(err)
	}

	remoteHTTP = remote_http.MakeClient(
		env,
		&httpRoundTripper,
		localRepo,
		cmd.MakeTypedInventoryListBlobStore(envRepo),
	)

	return
}
