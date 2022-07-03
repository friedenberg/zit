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

function outputs_organize_v2 { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)

	to_add="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- ok"
		echo "---"
	} >>"$to_add"

	run zit new "$to_add"
	assert_output --partial '[one/uno '

	expected_organize="$(mktemp)"
	{
		echo "# ok"
		echo ""
		echo "- [one/uno] wow"
	} >>"$expected_organize"

	run zit organize -group-by-unique ok
	assert_output "$(cat "$expected_organize")"

	{
		echo "# wow"
		echo ""
		echo "- [one/uno] wow"
	} >"$expected_organize"

	run zit organize -group-by-unique ok <"$expected_organize"

	expected_zettel="$(mktemp)"
	{
		echo "---"
		echo "# wow"
		echo "- wow"
		echo "---"
	} >>"$expected_zettel"

	run zit show one/uno
	assert_output "$(cat "$expected_zettel")"
}
