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
	skip
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
	skip
	run_zit init-workspace -query tag-3
	assert_success

	export EDITOR="true"
	run_zit edit
	assert_success
	assert_output - <<-EOM
		      checked out [one/uno.zettel @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		[one/uno @85eb98a5c8f7ccc354f35b846bb24adc1764e88cb907f63293f6902aa105af58]
	EOM

	run_zit show -format blob one/uno
	assert_success
	assert_output - <<-EOM
		this is the body 2
	EOM
}

function workspace_checkout { # @test
	skip
	run_zit init-workspace -tags tag-3
	assert_success

	run_zit checkout
	assert_success
	assert_output - <<-EOM
		      checked out [one/uno.zettel @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		[one/uno @85eb98a5c8f7ccc354f35b846bb24adc1764e88cb907f63293f6902aa105af58]
	EOM

	run_zit show -format blob one/uno
	assert_success
	assert_output - <<-EOM
		this is the body 2
	EOM
}
