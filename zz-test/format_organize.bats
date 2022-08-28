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

cat_yin() (
	echo "one"
	echo "two"
	echo "three"
)

cat_yang() (
	echo "uno"
	echo "dos"
	echo "tres"
)

function group_by_merges_child_into_matching_parent { # @test
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

	run zit format-organize "$to_add"
	assert_output "$(cat "$expected")"
}
