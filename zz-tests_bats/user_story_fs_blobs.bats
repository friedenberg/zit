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

function user_story_fs_blobs_organize() { # @test
	skip
	cat >test.md <<-EOM
		newest body
	EOM

	run_zit status .
	assert_success
	assert_output - <<-EOM
		        untracked [test.md @d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc]
	EOM

	run_zit organize -mode output-only .
	assert_success
	assert_output - <<-EOM
		        untracked [test.md @d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc
		                   test.md]
	EOM

	echo "test file" >"test file.md"
	run_zit status "test file.md"
	assert_success
	assert_output - <<-EOM
		       conflicted [one/dos.zettel]
	EOM
}

function user_story_fs_blobs_status_recognized() { # @test
	run_zit show -format blob one/uno
	echo "$output" >test.md

	run_zit status .
	assert_success
	assert_output - <<-EOM
		       recognized [one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4
		                   test.md]
	EOM
}
