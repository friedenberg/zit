package query

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/kilo/env_workspace"
)

type QueryBuilderModifier interface {
	ModifyBuilder(*Builder)
}

type BuilderOptionGetter interface {
	GetQueryBuilderOptions() builderOptionsOld
}

type BuilderOption interface {
	Apply(*Builder) *Builder
}

type (
	BuilderOptionsMulti []BuilderOption
	builderOptions      []BuilderOption
)

func BuilderOptions(options ...BuilderOption) builderOptions {
	return builderOptions(options)
}

func (options builderOptions) Apply(builder *Builder) *Builder {
	for _, option := range options {
		builder = option.Apply(builder)
	}

	return builder
}

type BuilderOptionWorkspace struct {
	env_workspace.Env
}

func (options BuilderOptionWorkspace) Apply(builder *Builder) *Builder {
	if options.Env == nil {
		return builder
	}

	if options.Env == nil {
		return builder
	}

	workspaceConfig := options.GetWorkspaceConfig()

	if workspaceConfig == nil {
		return builder
	}

	builder.workspaceEnabled = true
	defaultQueryGroup := workspaceConfig.GetDefaultQueryGroup()

	if defaultQueryGroup == "" {
		return builder
	}

	// TODO add after parsing as an independent query group, rather than as a
	// literal
	builder.defaultQuery = defaultQueryGroup

	return builder
}

type builderOptionsOld struct {
	QueryBuilderModifier
}

func BuilderOptionsOld(o any, remainder ...BuilderOption) builderOptions {
	var options builderOptionsOld

	if qbm, ok := o.(QueryBuilderModifier); ok {
		options.QueryBuilderModifier = qbm
	}

	return BuilderOptions(append([]BuilderOption{options}, remainder...)...)
}

func (options builderOptionsOld) Apply(b *Builder) *Builder {
	if options.QueryBuilderModifier != nil {
		options.ModifyBuilder(b)
	}

	return b
}

type options struct {
	defaultGenres  ids.Genre
	defaultSigil   ids.Sigil
	permittedSigil ids.Sigil
}

func BuilderOptionDefaultSigil(sigils ...ids.Sigil) builderOptionDefaultSigil {
	return builderOptionDefaultSigil(ids.MakeSigil(sigils...))
}

type builderOptionDefaultSigil ids.Sigil

func (option builderOptionDefaultSigil) Apply(builder *Builder) *Builder {
	builder.options.defaultSigil = ids.Sigil(option)
	return builder
}

type builderOptionDefaultGenre ids.Genre

func BuilderOptionDefaultGenres(
	genres ...genres.Genre,
) builderOptionDefaultGenre {
	return builderOptionDefaultGenre(ids.MakeGenre(genres...))
}

func (options builderOptionDefaultGenre) Apply(builder *Builder) *Builder {
	builder.options.defaultGenres = ids.Genre(options)
	return builder
}
