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
	echo "four"
	echo "five"
	echo "six"
)

cat_yang() (
	echo "uno"
	echo "dos"
	echo "tres"
	echo "quatro"
	echo "cinco"
	echo "seis"
)

function init_and_deinit { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	[[ -f .zit/KonfigCompiled ]]

	run zit deinit
	assert_success
}

function init_and_init { # @test
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init -yin <(cat_yin) -yang <(cat_yang)
	assert_success

	{
		echo "---"
		echo "# wow"
		echo "- tag"
		echo "! md"
		echo "---"
		echo
		echo "body"
	} >to_add

	run zit new -edit=false -predictable-hinweisen to_add
	assert_output '          (new) [o/u@8 !md "wow"]'

	run zit show one/uno
	assert_output "$(cat to_add)"

	run zit init -yin <(cat_yin) -yang <(cat_yang)
	assert_failure

	run zit init
	assert_output --partial '.zit/Kennung/Counter already exists, not overwriting'
	assert_output --partial '.zit/Konfig already exists, not overwriting'
	assert_output --partial '.zit/KonfigCompiled already exists, not overwriting'
	assert_output --partial '          (new) [o/u@8 !md "wow"]'

	# run zit reindex
	# assert_output "$(cat to_add)"

	# 	run tree .zit
	# 	assert_output "$(cat to_add)"

	run zit show one/uno
	assert_output "$(cat to_add)"
}
