#! /usr/bin/env bash -e

if [[ "$(($(ag skip zz-tests_bats/ -l | wc -l)))" -gt 0 ]]; then
	echo "The following tests have skips and so the deploy cannot go forward" >&2
	ag skip zz-tests_bats/ -l
	exit 1
fi

if [[ "$(($(ag "log.Debug" -l src/ | wc -l)))" -gt 0 ]]; then
	echo "The following files have debug logs and so the deploy cannot go forward" >&2
	ag "log.Debug" -l src/
	exit 1
fi

git pull --rebase

#TODO pause mr-build-and-watch and then resume after
cmd_make=make

if command -v gmake 2>&1 >/dev/null; then
	cmd_make=gmake
else
	make
fi

$cmd_make build/deploy
make install

go clean -cache -fuzzcache
git add .

if [[ "$(git status --porcelain=v1 2>/dev/null | wc -l)" -gt 0 ]]; then
	git commit -m update
fi

git push
