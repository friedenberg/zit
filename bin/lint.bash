#! /usr/bin/env bash -e

do_lint() (
	msg="$1"
	shift
	cmd=("$@")

	if [[ "$(($("${cmd[@]}" | wc -l)))" -gt 0 ]]; then
		echo "deploy aborted!" >&2
		echo "$msg" >&2
		"${cmd[@]}"
		echo
		echo "cmd: " "$(echo -n "${cmd[@]@Q}")" >&2
		exit 1
	fi
)

do_lint \
	"todos without priorities in the following files:" \
	ag \
	--go \
	'// todo(?!(\s*|-)p\d)' \
	--ignore-case \
	-l

do_lint \
	"bats test files have skips" \
	ag \
	skip \
	zz-tests_bats/ \
	-l

do_lint \
	"debug logs remaining" \
	ag \
	"log.Debug" \
	-l \
	src/

do_lint \
	"P0 comment todos" \
	ag \
	--go \
	'// todo(\s*|-)p0' \
	--ignore-case \
	-l

do_lint \
	"P0 compiled todos" \
	ag \
	--go \
	'errors\.TodoP0' \
	--ignore-case \
	-l
