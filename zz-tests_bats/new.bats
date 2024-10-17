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
	assert_success
	assert_output - <<-EOM
		[two/uno @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !md]
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
		[et1 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[et2 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[two/uno @036a8e44e472523c0306946f2712f372c234f8a24532e933f1509ae4db0da064 !md "bez" et1 et2]
	EOM

	run_zit new -edit=false "$expected"
	assert_success
	assert_output - <<-EOM
		[one/tres @036a8e44e472523c0306946f2712f372c234f8a24532e933f1509ae4db0da064 !md "bez" et1 et2]
	EOM

	# when
	run_zit show -format text two/uno
	assert_success
	assert_output "$(cat "$expected")"

	run_zit show -format text one/tres
	assert_success
	assert_output "$(cat "$expected")"
}

function use_blob_shas { # @test
	run_zit write-blob - <<-EOM
		  the blob
	EOM
	assert_success
	assert_output '6a405a5e357550175234d9f5b177014984f99fe34d35fe931cf8d2e96b8b0cb0 - (checked in)'

	run_zit new -edit=false -shas 6a405a5e357550175234d9f5b177014984f99fe34d35fe931cf8d2e96b8b0cb0
	assert_success
	assert_output - <<-EOM
		[two/uno @6a405a5e357550175234d9f5b177014984f99fe34d35fe931cf8d2e96b8b0cb0 !md]
	EOM

	the_blob2_sha="ad100d00763b333c0c8af3c89dbcc1f52f9c4a8a208476c35eb9d364121301b6"
	run_zit write-blob - <<-EOM
		  the blob2
	EOM
	assert_success
	assert_output "$the_blob2_sha - (checked in)"

	run_zit new -edit=false -shas -type txt "$the_blob2_sha"
	assert_success
	assert_output - <<-EOM
		[!txt @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/tres @$the_blob2_sha !txt]
	EOM

	run_zit_stderr_unified new -edit=false -shas "$the_blob2_sha"
	assert_success
	assert_output --partial - <<-EOM
		ad100d00763b333c0c8af3c89dbcc1f52f9c4a8a208476c35eb9d364121301b6 appears in object already checked in (["one/tres"]). Ignoring
	EOM
}
