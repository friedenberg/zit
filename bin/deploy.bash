#! /bin/bash -e

declare t

if [[ -d build_options ]]; then
	t="$(mktemp -d)"
	mv build_options "$t"
fi

#TODO pause mr-build-and-watch and then resume after
gmake
go clean -cache -fuzzcache
git add .
git commit -m update || true
git push

if [[ -d "$t" ]]; then
	mv "$t" build_options
fi
