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

function bootstrap {
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
	assert_output '          (new) [o/u@3 !md "wow"]'

	run zit show one/uno
	assert_output "$(cat to_add)"
}

function clone { # @test
	wd1="$(mktemp -d)"
	cd "$wd1" || exit 1
	bootstrap "$wd1"
	assert_success

	wd="$(mktemp -d)"
	cd "$wd" || exit 1

  # TODO P0 fix issue with non-deterministic sha abbreviations
	run zit clone -all -include-history -gattung zettel,typ "$wd1"
	assert_output --partial '(updated) [!md@e1d34e9ec6d4f741d0566dbf6683d3644c3b6b3b27f718a6c09668a906c7df51]'
	assert_output --partial '(updated) [konfig@e]'
	assert_output --partial '(updated) [!md@e1]'
	assert_output --partial '(updated) [konfig@e1d6]'
	assert_output --partial '    (new) [one/uno@3 !md "wow"]'
}
