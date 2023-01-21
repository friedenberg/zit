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

cmd_zit_def=(
	-abbreviate-hinweisen=true
	-predictable-hinweisen
	-print-typen=false
)

function can_update_akte { # @test
	# setup
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run zit init "${cmd_zit_def[@]}" -disable-age -yin <(cat_yin) -yang <(cat_yang)
	assert_success

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
	} >"$expected"

	run zit new "${cmd_zit_def[@]}" -edit=false -predictable-hinweisen "$expected"
	assert_output '          (new) [o/u@18 !md "bez"]'

	run zit show "${cmd_zit_def[@]}" one/uno
	assert_output "$(cat "$expected")"

	# when
	new_akte="$(mktemp)"
	{
		echo the body but new
	} >"$new_akte"

	run zit checkin-akte "${cmd_zit_def[@]}" -new-etiketten et3 one/uno "$new_akte"
	assert_output '      (updated) [o/u@6 !md "bez"]'

	# then
	{
		echo ---
		echo "# bez"
		echo - et3
		echo ! md
		echo ---
		echo
		echo the body but new
	} >"$expected"

	run zit show "${cmd_zit_def[@]}" one/uno
	assert_output "$(cat "$expected")"
}
