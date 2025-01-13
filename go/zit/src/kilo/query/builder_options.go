package query

import (
	"code.linenisgreat.com/zit/go/zit/src/echo/ids"
)

type QueryBuilderModifier interface {
	ModifyBuilder(*Builder)
}

type DefaultSigilGetter interface {
	DefaultSigil() ids.Sigil
}

type DefaultGenresGetter interface {
	DefaultGenres() ids.Genre
}

type BuilderOptionGetter interface {
	GetQueryBuilderOptions() BuilderOptions
}

type BuilderOptions struct {
	QueryBuilderModifier
	DefaultSigilGetter
	DefaultGenresGetter
}

func MakeBuilderOptions(o any) BuilderOptions {
	var options BuilderOptions

	if dgg, ok := o.(DefaultGenresGetter); ok {
		options.DefaultGenresGetter = dgg
	}

	if dsg, ok := o.(DefaultSigilGetter); ok {
		options.DefaultSigilGetter = dsg
	}

	if qbm, ok := o.(QueryBuilderModifier); ok {
		options.QueryBuilderModifier = qbm
	}

	return options
}

func (options BuilderOptions) Apply(b *Builder) {
	if options.DefaultGenresGetter != nil {
		b = b.WithDefaultGenres(options.DefaultGenres())
	}

	if options.DefaultSigilGetter != nil {
		b.WithDefaultSigil(options.DefaultSigil())
	}

	if options.QueryBuilderModifier != nil {
		options.ModifyBuilder(b)
	}
}
