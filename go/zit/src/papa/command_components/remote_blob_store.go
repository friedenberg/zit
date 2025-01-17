package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/golf/env"
	"code.linenisgreat.com/zit/go/zit/src/hotel/repo_layout"
)

type RemoteBlobStore struct {
	Blobs  string
	Config immutable_config.BlobStoreTomlV1
}

func (cmd *RemoteBlobStore) SetFlagSet(f *flag.FlagSet) {
	cmd.Config.CompressionType = immutable_config.CompressionTypeDefault
	cmd.Config.CompressionType.SetFlagSet(f)
	f.StringVar(&cmd.Blobs, "blobs", "", "")
}

func (cmd *RemoteBlobStore) MakeRemoteBlobStore(
	e env.LocalEnv,
) (blobStore interfaces.BlobStore, err error) {
	blobStore = repo_layout.MakeBlobStore(
		cmd.Blobs,
		env_dir.MakeConfigFromImmutableBlobConfig(&cmd.Config),
		e.GetTempLocal(),
	)

	return
}
