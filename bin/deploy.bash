#! /usr/bin/env nix-shell
#! nix-shell -i bash

declare t

if [[ -d build_options ]]; then
	t="$(mktemp -d)"
	mv build_options "$t"
fi

#TODO pause mr-build-and-watch and then resume after
make
go clean -cache -fuzzcache
git add .

if [[ "$(git status --porcelain=v1 2>/dev/null | wc -l)" -gt 0 ]]; then
	git commit -m update
fi

git push

if [[ -d "$t" ]]; then
	mv "$t" build_options
fi
