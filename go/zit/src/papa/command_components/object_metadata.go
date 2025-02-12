package command_components

import (
	"strings"

	"code.linenisgreat.com/zit/go/zit/src/alfa/errors"
	"code.linenisgreat.com/zit/go/zit/src/alfa/flag_policy"
	"code.linenisgreat.com/zit/go/zit/src/bravo/flag"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata"
)

type ObjectMetadata struct{}

func (cmd ObjectMetadata) GetFlagValueMetadataTags(
	metadata *object_metadata.Metadata,
) flag.Value {
	// TODO add support for tag_paths
	fes := flag.Make(
		flag_policy.FlagPolicyAppend,
		func() string {
			return metadata.Cache.TagPaths.String()
		},
		func(v string) (err error) {
			vs := strings.Split(v, ",")

			for _, v := range vs {
				if err = metadata.AddTagString(v); err != nil {
					err = errors.Wrap(err)
					return
				}
			}

			return
		},
		func() {
			metadata.ResetTags()
		},
	)

	return fes
}

func (cmd ObjectMetadata) GetFlagValueMetadataDescription(
	metadata *object_metadata.Metadata,
) flag.Value {
	return &metadata.Description
}

func (cmd ObjectMetadata) GetFlagValueMetadataType(
	metadata *object_metadata.Metadata,
) flag.Value {
	return &metadata.Type
}
