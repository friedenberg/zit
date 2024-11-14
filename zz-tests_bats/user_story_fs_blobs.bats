#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	copy_from_version "$DIR" "$version"

	cat >test.md <<-EOM
		newest body
	EOM
}

teardown() {
	rm_from_version "$version"
}

function user_story_fs_blobs_status() { # @test
	run_zit status .
	assert_success
	assert_output - <<-EOM
		        untracked [test.md @d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc]
	EOM
}

function user_story_fs_blobs_organize_output_only() { # @test
	run_zit organize -mode output-only .
	assert_success
	assert_output - <<-EOM

		- [test.md]
	EOM
}

function user_story_fs_blobs_checkin_dot() { # @test
	run_zit checkin .
	assert_success
	assert_output - <<-EOM
		[two/uno @d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc !md "test"]
	EOM
}

function user_story_fs_blobs_checkin_explicit() { # @test
	run_zit checkin test.md
	assert_success
	assert_output - <<-EOM
		[two/uno @d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc !md "test"]
	EOM
}

function user_story_fs_blobs_checkin_organize_dot_exclude() { # @test
	export EDITOR="bash -c 'echo > \"\$0\"'"
	run_zit checkin -organize .
	assert_success
	assert_output ''
}

function user_story_fs_blobs_checkin_organize_dot_include() { # @test
	export EDITOR="bash -c 'true'"
	run_zit checkin -organize . </dev/null
	assert_success
	assert_output - <<-EOM
		[two/uno @d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc !md "test"]
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
