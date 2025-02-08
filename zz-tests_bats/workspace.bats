#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit info store-version)"
	copy_from_version "$DIR" "$version"
}

teardown() {
	rm_from_version "$version"
}

# bats file_tags=user_story:workspace

function workspace_show { # @test
	run_zit init-workspace -query tag-3
	assert_success

	run_zit show
	assert_success
	assert_output_unsorted - <<-eom
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
	eom

	run_zit show one/uno
	assert_success
	assert_output - <<-eom
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	eom
}

function workspace_edit { # @test
	run_zit init-workspace -query tag-3
	assert_success

	export EDITOR="true"
	run_zit edit
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [tag-3.tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		      checked out [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		      checked out [one/uno.zettel @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM

	run_zit show -format blob one/uno
	assert_success
	assert_output - <<-EOM
		last time
	EOM
}

function workspace_checkout { # @test
	run_zit init-workspace -tags tag-3
	assert_success

	run_zit checkout
	assert_success
	assert_output ''

	run_zit checkout :
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [one/dos.zettel @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		      checked out [one/uno.zettel @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM

	run_zit show -format blob one/uno.zettel
	assert_success
	assert_output - <<-EOM
		last time
	EOM
}

function workspace_organize { # @test
	run_zit init-workspace -tags tag-3
	assert_success

	run_zit organize -mode output-only
	assert_success
	assert_output - <<-EOM
		---
		- tag-3
		---
	EOM

	run_zit organize -mode output-only :
	assert_success
	assert_output - <<-EOM

		- [one/dos !md tag-3 tag-4] wow ok again
		- [one/uno !md tag-3 tag-4] wow the first
	EOM

	run_zit organize -mode output-only one/uno
	assert_success
	assert_output - <<-EOM

		- [one/uno !md tag-3 tag-4] wow the first
	EOM
}
