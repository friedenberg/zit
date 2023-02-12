#! /usr/bin/env bats

setup() {
	load 'test_helper/bats-support/load'
	load 'test_helper/bats-assert/load'
	load 'common.bash'
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

function format_organize_right_align { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	to_add="$(mktemp)"
	{
		echo "# task"
		echo "## urgency"
		echo "### urgency-1"
		echo "### -2"
	} >"$to_add"

	expected="$(mktemp)"
	{
		echo
		echo "    # task"
		echo
		echo "   ## urgency"
		echo
		echo "  ###        -1"
		echo
		echo "  ###        -2"
		echo
	} >"$expected"

	run_zit format-organize -prefix-joints=true -refine "$to_add"
	assert_output "$(cat "$expected")"
}

function format_organize_left_align { # @test
	skip
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

	to_add="$(mktemp)"
	{
		echo "# task"
		echo "## urgency"
		echo "### urgency-1"
		echo "### -2"
	} >"$to_add"

	expected="$(mktemp)"
	{
		echo
		echo "# task"
		echo
		echo " ## urgency"
		echo
		echo "  ### -1"
		echo
		echo "  ### -2"
		echo
	} >"$expected"

	run_zit format-organize -prefix-joints=true -refine -right-align=false "$to_add"
	assert_output "$(cat "$expected")"
}
