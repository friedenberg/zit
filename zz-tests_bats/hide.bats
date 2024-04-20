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
		[konfig@9faf629a95cd95b7c3cd3b8d69a11c0d592b68e7929fe0b5fcf8912db0aaee8a]
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
		[konfig@2add70c77fac6a831bb9339e8a43446c06d4e92506e5d01d3ee9247b16ebbdbe]
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
