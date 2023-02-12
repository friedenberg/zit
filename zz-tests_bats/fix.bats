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

function can_update_akte { # @test
	# setup
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

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

	run_zit new -edit=false "$expected"
	assert_output '[one/uno@18df16846a2f8bbce5f03e1041baff978a049aabd169ab9adac387867fe1706c !md "bez"]'

	run_zit show one/uno
	assert_output "$(cat "$expected")"

	# when
	new_akte="$(mktemp)"
	{
		echo the body but new
	} >"$new_akte"

	run_zit checkin-akte -new-etiketten et3 one/uno "$new_akte"
	assert_output '[one/uno@6b4905e7d7a5185f73db1e27448663fa38b3aca11d62e1dc33ecb066653791b7 !md "bez"]'

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

	run_zit show one/uno
	assert_output "$(cat "$expected")"
}
