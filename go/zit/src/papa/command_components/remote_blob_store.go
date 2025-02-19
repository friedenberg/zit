package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/config_immutable"
	"code.linenisgreat.com/zit/go/zit/src/echo/env_dir"
	"code.linenisgreat.com/zit/go/zit/src/hotel/blob_store"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
)

type RemoteBlobStore struct {
	Blobs  string
	Config config_immutable.BlobStoreTomlV1
}

func (cmd *RemoteBlobStore) SetFlagSet(f *flag.FlagSet) {
	cmd.Config.CompressionType = config_immutable.CompressionTypeDefault
	cmd.Config.CompressionType.SetFlagSet(f)
	f.StringVar(&cmd.Blobs, "blobs", "", "")
}

func (cmd *RemoteBlobStore) MakeRemoteBlobStore(
	e env_local.Env,
) (blobStore interfaces.BlobStore, err error) {
	blobStore = blob_store.MakeShardedFilesStore(
		cmd.Blobs,
		env_dir.MakeConfigFromImmutableBlobConfig(&cmd.Config),
		e.GetTempLocal(),
	)

	return
}
