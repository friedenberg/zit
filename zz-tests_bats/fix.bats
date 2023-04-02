#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output
}

function can_update_akte { # @test
	# setup
	wd="$(mktemp -d)"
	cd "$wd" || exit 1

	run_zit_init_disable_age

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

	run_zit new -edit=false "$expected"
	assert_success
	assert_output_unsorted - <<-EOM
		[-et2@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-et1@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[one/uno@18df16846a2f8bbce5f03e1041baff978a049aabd169ab9adac387867fe1706c !md "bez"]
	EOM

	run_zit show one/uno
	assert_success
	assert_output "$(cat "$expected")"

	# when
	new_akte="$(mktemp)"
	{
		echo the body but new
	} >"$new_akte"

	run_zit checkin-akte -new-etiketten et3 one/uno "$new_akte"
	assert_success
	assert_output - <<-EOM
		[-et3@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[one/uno@6b4905e7d7a5185f73db1e27448663fa38b3aca11d62e1dc33ecb066653791b7 !md "bez"]
	EOM

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
	assert_success
	assert_output "$(cat "$expected")"
}
