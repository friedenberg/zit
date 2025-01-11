package repo_layout

import "code.linenisgreat.com/zit/go/zit/src/delta/immutable_config"

type config struct {
	immutable_config.Config
	storeVersion      immutable_config.StoreVersion
	compressionType   immutable_config.CompressionType
	lockInternalFiles bool
}
