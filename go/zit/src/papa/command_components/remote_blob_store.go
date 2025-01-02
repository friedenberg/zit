package command_components

import (
	"flag"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/interfaces"
	"code.linenisgreat.com/zit/go/zit/src/delta/age"
	"code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"
	"code.linenisgreat.com/zit/go/zit/src/echo/repo_layout"
)

type RemoteBlobStore struct {
	Blobs           string
	AgeIdentity     age.Identity
	CompressionType immutable_config.CompressionType
}

func (cmd *RemoteBlobStore) SetFlagSet(f *flag.FlagSet) {
	cmd.CompressionType = immutable_config.CompressionTypeDefault
	f.StringVar(&cmd.Blobs, "blobs", "", "")
	f.Var(&cmd.AgeIdentity, "age-identity", "")
	cmd.CompressionType.SetFlagSet(f)
}

func (cmd *RemoteBlobStore) MakeRemoteBlobStore() (blobStore interfaces.BlobStore, err error) {
	var ag age.Age

	if err = ag.AddIdentity(cmd.AgeIdentity); err != nil {
		err = errors.Wrapf(err, "age-identity: %q", &cmd.AgeIdentity)
		return
	}

	blobStore = repo_layout.MakeBlobStore(
		cmd.Blobs,
		&ag,
		cmd.CompressionType,
	)

	return
}
