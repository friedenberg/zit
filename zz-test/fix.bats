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

function can_peek_hinweisen_after_init { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -verbose -disable-age -yin <(cat_yin) -yang <(cat_yang)

	expected="$(mktemp)"
	{
		echo 0: one/dos
		echo 1: one/tres
		echo 2: one/uno
		echo 3: three/dos
		echo 4: three/tres
		echo 5: three/uno
		echo 6: two/dos
		echo 7: two/tres
    echo 8: two/uno
	} >"$expected"

	run zit peek-hinweisen
	assert_output "$(cat "$expected")"
}
