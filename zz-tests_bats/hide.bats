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
		[konfig@7c7e3c0dbfef03649c5c15f52617be8f75d3836971b8720d8441207a13ab2598]
	EOM

	run_zit show :z
	assert_success
	assert_output ''

	run_zit show :?z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
	EOM

	run_zit edit-konfig -unhide-etikett tag-3
	assert_success
	assert_output - <<-EOM
		[konfig@67709f55748b0aaab9723e8ad63f6877d1a2b04d0124461a4e6e313260486dad]
	EOM

	# run_zit reindex
	# assert_success

	#TODO [act/zu "fix issue with -unhide-etikett erasing the hidden zettels from the index on flush and requiring this reindex to work"

	run_zit show :z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
	EOM
}
