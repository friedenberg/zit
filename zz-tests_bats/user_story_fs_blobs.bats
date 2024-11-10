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
