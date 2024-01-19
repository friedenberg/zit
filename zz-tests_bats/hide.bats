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

function mark_one_as_hidden { # @test
	run_zit show :z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
	EOM

	run_zit edit-konfig -hide-etikett tag-3
	assert_success
	assert_output - <<-EOM
		[konfig@6fab60685dafa09e8dc6b714cc50ea605622bff5f1754919d0487a682270c49c]
	EOM

	run_zit show :z
	assert_success
	assert_output ''

	run_zit edit-konfig -unhide-etikett tag-3
	assert_success
	assert_output - <<-EOM
		[konfig@b6f18376dd601afd3c7d8692ec5b407727f5268760b7bc43884aba804cca9f88]
	EOM

	run_zit show :z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
	EOM
}
