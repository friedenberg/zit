package command_components

import (
	"code.linenisgreat.com/zit/go/zit/src/bravo/flag"
	"code.linenisgreat.com/zit/go/zit/src/bravo/values"
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/golf/command"
	"code.linenisgreat.com/zit/go/zit/src/golf/object_metadata"
	"code.linenisgreat.com/zit/go/zit/src/hotel/env_local"
	"code.linenisgreat.com/zit/go/zit/src/juliett/sku"
	"code.linenisgreat.com/zit/go/zit/src/kilo/query"
	"code.linenisgreat.com/zit/go/zit/src/kilo/sku_fmt"
	"code.linenisgreat.com/zit/go/zit/src/november/local_working_copy"
)

type Complete struct {
	ObjectMetadata
	Query
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

func (cmd Complete) SetFlagsProto(
	proto *sku.Proto,
	flagSet *flag.FlagSet,
	descriptionUsage string,
	tagUsage string,
	typeUsage string,
) {
	proto.SetFlagSetDescription(
		flagSet,
		descriptionUsage,
	)

	flagSet.Var(
		cmd.GetFlagValueMetadataTags(&proto.Metadata),
		"tags",
		tagUsage,
	)

	flagSet.Var(
		cmd.GetFlagValueMetadataType(&proto.Metadata),
		"type",
		typeUsage,
	)
}

func (cmd Complete) CompleteObjects(
	req command.Request,
	local *local_working_copy.Repo,
	queryBuilderOptions query.BuilderOption,
	args ...string,
) {
	printerCompletions := sku_fmt.MakePrinterComplete(local)

	queryGroup := cmd.MakeQuery(
		req,
		queryBuilderOptions,
		local,
		args,
	)

	if err := local.GetStore().QueryTransacted(
		queryGroup,
		printerCompletions.PrintOne,
	); err != nil {
		local.CancelWithError(err)
	}
}
