package query

import (
	"code.linenisgreat.com/zit/go/zit/src/delta/genres"
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
	"code.linenisgreat.com/zit/go/zit/src/india/env_workspace"
)

type QueryBuilderModifier interface {
	ModifyBuilder(*Builder)
}

type DefaultGenresGetter interface {
	DefaultGenres() ids.Genre
}

type BuilderOptionGetter interface {
	GetQueryBuilderOptions() builderOptionsInterfaces
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
	if options.Env != nil {
		workspaceConfig := options.GetWorkspaceConfig()

		if workspaceConfig != nil {
			defaultQueryGroup := workspaceConfig.GetDefaultQueryGroup()

			// TODO add after parsing as an independent query group, rather than as a
			// literal
			if defaultQueryGroup != "" {
				builder.defaultQuery = defaultQueryGroup
			}
		}
	}

	return builder
}

type builderOptionsInterfaces struct {
	QueryBuilderModifier
	DefaultGenresGetter
}

func MakeBuilderOptions(o any) builderOptionsInterfaces {
	var options builderOptionsInterfaces

	if dgg, ok := o.(DefaultGenresGetter); ok {
		options.DefaultGenresGetter = dgg
	}

	if qbm, ok := o.(QueryBuilderModifier); ok {
		options.QueryBuilderModifier = qbm
	}

	return options
}

func (options builderOptionsInterfaces) Apply(b *Builder) *Builder {
	if options.DefaultGenresGetter != nil {
		b = b.WithDefaultGenres(options.DefaultGenres())
	}

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
