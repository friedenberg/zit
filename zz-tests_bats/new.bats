#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	copy_from_version "$DIR" "$version"
}

teardown() {
	rm_from_version "$version"
}

function new_empty_no_edit { # @test
	run_zit new -edit=false
	assert_success
	assert_output - <<-EOM
		[two/uno@e6e789716abc939fc15b8caae85ecb9c1bbe96d44d1b58d2fd42a2a8fd9d904a !md ]
	EOM
}

function can_duplicate_zettel_content { # @test
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
	assert_output - <<-EOM
		[-et1@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-et2@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[two/uno@18df16846a2f8bbce5f03e1041baff978a049aabd169ab9adac387867fe1706c !md "bez"]
	EOM

	run_zit new -edit=false "$expected"
	assert_success
	assert_output - <<-EOM
		[one/tres@18df16846a2f8bbce5f03e1041baff978a049aabd169ab9adac387867fe1706c !md "bez"]
	EOM

	# when
	run_zit show two/uno
	assert_success
	assert_output "$(cat "$expected")"

	run_zit show one/tres
	assert_success
	assert_output "$(cat "$expected")"
}
