package command_components

import (
	"os"

	"code.linenisgreat.com/zit/go/zit/src/bravo/flag"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type Complete struct {
	ObjectMetadata
	QueryGroup
}

func (cmd Complete) GetFlagValueMetadataTags(
	metadata *object_metadata.Metadata,
) flag.Value {
	return command.FlagValueCompleter{
		Value: cmd.ObjectMetadata.GetFlagValueMetadataTags(metadata),
		FuncCompleter: func(
			req command.Request,
			envLocal env_local.Env,
			commandLine command.CommandLine,
		) {
			local := LocalWorkingCopy{}.MakeLocalWorkingCopy(req)

			cmd.CompleteObjects(
				req,
				local,
				query.BuilderOptionDefaultGenres(genres.Tag),
			)
		},
	}
}

func (cmd Complete) GetFlagValueStringTags(
	value *values.String,
) flag.Value {
	return command.FlagValueCompleter{
		Value: value,
		FuncCompleter: func(
			req command.Request,
			envLocal env_local.Env,
			commandLine command.CommandLine,
		) {
			local := LocalWorkingCopy{}.MakeLocalWorkingCopy(req)

			cmd.CompleteObjects(
				req,
				local,
				query.BuilderOptionDefaultGenres(genres.Tag),
			)
		},
	}
}

func (cmd Complete) GetFlagValueMetadataType(
	metadata *object_metadata.Metadata,
) flag.Value {
	return command.FlagValueCompleter{
		Value: cmd.ObjectMetadata.GetFlagValueMetadataType(metadata),
		FuncCompleter: func(
			req command.Request,
			envLocal env_local.Env,
			commandLine command.CommandLine,
		) {
			local := LocalWorkingCopy{}.MakeLocalWorkingCopy(req)

			cmd.CompleteObjects(
				req,
				local,
				query.BuilderOptionDefaultGenres(genres.Type),
			)
		},
	}
}

func (cmd Complete) CompleteObjects(
	req command.Request,
	local *local_working_copy.Repo,
	queryBuilderOptions query.BuilderOption,
	args ...string,
) {
	completionWriter := sku_fmt.MakeWriterComplete(os.Stdout)
	defer local.MustClose(completionWriter)

	queryGroup := cmd.MakeQueryGroup(
		req,
		queryBuilderOptions,
		local,
		args,
	)

	if err := local.GetStore().QueryTransacted(
		queryGroup,
		completionWriter.WriteOneTransacted,
	); err != nil {
		local.CancelWithError(err)
	}
}
