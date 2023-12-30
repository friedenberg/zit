#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

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
	assert_failure
	assert_output - <<-EOM
		zettel is empty
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
		[-et1@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-et2@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[two/uno@036a8e44e472523c0306946f2712f372c234f8a24532e933f1509ae4db0da064 !md "bez" et1 et2]
	EOM

	run_zit new -edit=false "$expected"
	assert_success
	assert_output - <<-EOM
		[one/tres@036a8e44e472523c0306946f2712f372c234f8a24532e933f1509ae4db0da064 !md "bez" et1 et2]
	EOM

	# when
	run_zit show -format text two/uno
	assert_success
	assert_output "$(cat "$expected")"

	run_zit show -format text one/tres
	assert_success
	assert_output "$(cat "$expected")"
}
