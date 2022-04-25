#! /usr/bin/env bats

setup() {
	load 'test_helper/bats-support/load'
	load 'test_helper/bats-assert/load'
	# ... the remaining setup is unchanged

	# get the containing directory of this file
	# use $BATS_TEST_FILENAME instead of ${BASH_SOURCE[0]} or $0,
	# as those will point to the bats executable's location or the preprocessed file respectively
	DIR="$(cd "$(dirname "$BATS_TEST_FILENAME")" >/dev/null 2>&1 && pwd)"
	# make executables in src/ visible to PATH
	PATH="$DIR/../:$PATH"
	PATH="$DIR/../build/:$PATH"

	# for shellcheck SC2154
	export output
}

function can_run_zit { # @test
	run zit
}

function provides_help_with_no_params { # @test
	run zit
	assert_output --partial 'No subcommand provided.'
}

function can_initialize_without_age { # @test
	yin="$(mktemp)"
	{
		echo "one"
		echo "two"
		echo "three"
	} >>"$yin"

	yang="$(mktemp)"
	{
		echo "uno"
		echo "dos"
		echo "tres"
	} >>"$yang"

	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin "$yin" -yang "$yang"
	[ -d .zit/ ]
	[ ! -f .zit/AgeIdentity ]
}

function can_new_zettel { # @test
	yin="$(mktemp)"
	{
		echo "one"
		echo "two"
		echo "three"
	} >>"$yin"

	yang="$(mktemp)"
	{
		echo "uno"
		echo "dos"
		echo "tres"
	} >>"$yang"

	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin "$yin" -yang "$yang"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "---"
	} >>"$to_add"

	run zit new "$to_add"
	assert_output --partial '[one/uno '

	run zit show one/uno
	[ "$(cat "$to_add")" = "$output" ]
}

function can_checkout_and_checkin { # @test
	yin="$(mktemp)"
	{
		echo "one"
		echo "two"
		echo "three"
	} >>"$yin"

	yang="$(mktemp)"
	{
		echo "uno"
		echo "dos"
		echo "tres"
	} >>"$yang"

	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin "$yin" -yang "$yang"

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "---"
	} >>"$to_add"

	run zit new "$to_add"
	assert_output --partial '[one/uno '

	run zit checkout one/uno
	assert_output --partial '[one/uno '
	assert_output --partial '(checked out)'

	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "---"
		echo ""
		echo "content"
	} >"one/uno.md"

	run zit checkin one/uno
	assert_output --partial '[one/uno '
	assert_output --partial '(updated)'
}
