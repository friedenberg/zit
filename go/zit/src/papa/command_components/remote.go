package command_components

import (
	"flag"
	"fmt"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/delta/xdg"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_blobs"
	"code.linenisgreat.com/zit/go/zit/src/foxtrot/builtin_types"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/env_ui"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/lima/repo"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
	"code.linenisgreat.com/zit/go/zit/src/oscar/remote_http"
)

type Remote struct {
	Env
	EnvRepo
	LocalWorkingCopy
	LocalArchive

	RemoteType repo.RemoteType
}

func (cmd *Remote) SetFlagSet(f *flag.FlagSet) {
	f.Var(&cmd.RemoteType, "remote-type", fmt.Sprintf("%s", repo.GetAllRemoteTypes()))
}

func (cmd Remote) CreateRemote(
	req command.Request,
	local *local_working_copy.Repo,
	proto sku.Proto,
) (sk *sku.Transacted) {
	env := cmd.MakeEnv(req)

	sk = sku.GetTransactedPool().Get()
	proto.Apply(sk.GetMetadata(), genres.Repo)

	var id ids.RepoId
	var blob repo_blobs.Blob

	switch cmd.RemoteType {
	default:
		req.CancelWithBadRequestf("unsupported remote type: %q", cmd.RemoteType)

	case repo.RemoteTypeNativeDotenvXDG:
		xdgDotenvPath := req.PopArg("xdg-dotenv-path")

		if err := id.Set(req.PopArg("repo-id")); err != nil {
			req.CancelWithError(err)
		}

		envLocal := cmd.MakeEnvWithXDGLayoutAndOptions(
			req,
			xdgDotenvPath,
			env.GetOptions(),
		)

		sk.Metadata.Type = builtin_types.GetOrPanic(builtin_types.RepoTypeXDGDotenvV0).Type
		blob = repo_blobs.TomlXDGV0FromXDG(envLocal.GetXDG())

	case repo.RemoteTypeStdioLocal:
		path := req.PopArg("path")

		if err := id.Set(req.PopArg("repo-id")); err != nil {
			req.CancelWithError(err)
		}

		sk.Metadata.Type = builtin_types.GetOrPanic(builtin_types.RepoTypeLocalPath).Type
		blob = repo_blobs.TomlLocalPathV0{Path: local.AbsFromCwdOrSame(path)}
	}

	req.AssertNoMoreArgs()

	var blobSha interfaces.Sha

	{
		var err error

		if blobSha, _, err = local.GetStore().GetTypedBlobStore().Repo.WriteTypedBlob(
			sk.Metadata.Type,
			blob,
		); err != nil {
			req.CancelWithError(err)
		}
	}

	sk.Metadata.Blob.ResetWithShaLike(blobSha)

	if err := sk.ObjectId.SetWithIdLike(&id); err != nil {
		req.CancelWithError(err)
	}

	if err := local.GetStore().CreateOrUpdate(
		sk,
		sku.StoreOptions{
			ApplyProto: true,
		},
	); err != nil {
		req.CancelWithError(err)
	}

	return
}

func (cmd Remote) MakeRemote(
	req command.Request,
	local *local_working_copy.Repo,
	sk *sku.Transacted,
) (remote repo.Repo) {
	var blob repo_blobs.Blob

	{
		var err error

		if blob, _, err = local.GetStore().GetTypedBlobStore().Repo.ReadTypedBlob(
			sk.Metadata.Type,
			sk.GetBlobSha(),
		); err != nil {
			req.CancelWithError(err)
		}
	}

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
		)

	case repo_blobs.V0:
		req.CancelWithErrorf("unsupported repo blob type: %T", blob)

	default:
		req.CancelWithErrorf("unsupported repo blob type: %T", blob)
	}

	return
}

func (cmd Remote) MakeArchiveFromArg(
	req command.Request,
	remoteArg string,
	local repo.Repo,
) (remote repo.Repo) {
	env := cmd.MakeEnv(req)

	switch cmd.RemoteType {
	case repo.RemoteTypeNativeDotenvXDG:
		remote = cmd.MakeLocalWorkingCopyFromConfigAndXDGDotenvPath(
			req,
			remoteArg,
			env.GetOptions(),
		)

	case repo.RemoteTypeStdioLocal:
		remote = cmd.MakeRemoteStdioLocal(
			req,
			env,
			remoteArg,
			local,
		)

	case repo.RemoteTypeStdioSSH:
		remote = cmd.MakeRemoteStdioSSH(
			req,
			env,
			remoteArg,
			local,
		)

	case repo.RemoteTypeSocketUnix:
		remote = cmd.MakeRemoteHTTPFromXDGDotenvPath(
			req,
			remoteArg,
			env.GetOptions(),
			local,
		)

	default:
		req.CancelWithBadRequestf("unsupported remote type: %q", cmd.RemoteType)
	}

	return
}

func (cmd Remote) MakeRemoteWorkingCopyFromArg(
	req command.Request,
	remoteArg string,
	local repo.Repo,
) (remote repo.WorkingCopy) {
	switch cmd.RemoteType {
	case repo.RemoteTypeNativeDotenvXDG:
		remote = cmd.MakeLocalWorkingCopyFromConfigAndXDGDotenvPath(
			req,
			remoteArg,
			env_ui.Options{},
		)

	case repo.RemoteTypeStdioLocal:
		remote = cmd.MakeRemoteStdioLocal(
			req,
			cmd.MakeEnv(req),
			remoteArg,
			local,
		)

	case repo.RemoteTypeStdioSSH:
		remote = cmd.MakeRemoteStdioSSH(
			req,
			cmd.MakeEnv(req),
			remoteArg,
			local,
		)

	case repo.RemoteTypeSocketUnix:
		remote = cmd.MakeRemoteHTTPFromXDGDotenvPath(
			req,
			remoteArg,
			env_ui.Options{},
			local,
		)

	default:
		req.CancelWithBadRequestf("unsupported remote type: %q", cmd.RemoteType)
	}

	return
}

func (cmd *Remote) MakeRemoteHTTPFromXDGDotenvPath(
	req command.Request,
	xdgDotenvPath string,
	options env_ui.Options,
	localRepo repo.Repo,
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

	if err := httpRoundTripper.Initialize(server); err != nil {
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
		localRepo.GetInventoryListStore(),
		cmd.MakeTypedInventoryListBlobStore(envRepo),
	)

	return
}

func (cmd *Remote) MakeRemoteStdioSSH(
	req command.Request,
	env env_local.Env,
	arg string,
	local repo.Repo,
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
		local.GetInventoryListStore(),
		cmd.MakeTypedInventoryListBlobStore(envRepo),
	)

	return
}

func (cmd *Remote) MakeRemoteStdioLocal(
	req command.Request,
	env env_local.Env,
	dir string,
	localRepo repo.Repo,
) (remoteHTTP repo.WorkingCopy) {
	envRepo := cmd.MakeEnvRepo(req, false)

	var httpRoundTripper remote_http.RoundTripperStdio

	httpRoundTripper.Dir = dir

	if err := httpRoundTripper.InitializeWithLocal(envRepo); err != nil {
		env.CancelWithError(err)
	}

	remoteHTTP = remote_http.MakeClient(
		env,
		&httpRoundTripper,
		localRepo.GetInventoryListStore(),
		cmd.MakeTypedInventoryListBlobStore(envRepo),
	)

	return
}
