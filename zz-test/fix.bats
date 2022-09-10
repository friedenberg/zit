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

function can_peek_hinweisen_after_init_fix { # @test
	# setup
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -verbose -disable-age -yin <(cat_yin) -yang <(cat_yang)

	expected="$(mktemp)"
	{
		echo ---
		echo "# bez"
		echo - et1
		echo - et2
		echo ! md
		echo ---
		echo
		echo the body
	} >>"$expected"

	run zit new -verbose -edit=false -predictable-hinweisen "$expected"
	assert_output --partial '[one/uno '

	run zit checkout one/uno
	assert_output --partial '[one/uno '
	assert_output --partial '(checked out)'

	run cat one/uno.md
	assert_output "$(cat "$expected")"

	{
		echo ---
		echo "# bez"
		echo - et1
		echo - et2
		echo ! md
		echo ---
		echo
		echo the body 2
	} >"$expected"

	cat "$expected" >"one/uno.md"

	run zit checkout -verbose one/uno
	assert_output --partial '[one/uno '
	assert_output --partial '(external has changes)'

	run cat one/uno.md
	assert_output "$(cat "$expected")"
}
