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

function add_one { # @test
	run_zit checkout tag-3:e
	assert_success
	# assert_output - <<-EOM
	# 	      checked out [tag-3.etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
	# EOM

	cat - >tag-3.etikett <<-EOM
		---
		- recurse
		---

	EOM

	run_zit checkin -delete .e
	assert_success
	assert_output - <<-EOM
		[recurse@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-3@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 recurse]
		          deleted [tag-3.etikett]
	EOM

	run_zit show recurse:e,z
	assert_success
	assert_output_unsorted - <<-EOM
		[recurse@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-3@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 recurse]
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM
}

function add_one_super_etiketten { # @test
	run_zit checkout tag-3:e
	assert_success
	assert_output - <<-EOM
		      checked out [tag-3.etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
	EOM

	cat - >tag-3.etikett <<-EOM
		---
		- recurse
		---

	EOM

	run_zit checkin -delete .e
	assert_success
	assert_output - <<-EOM
		[recurse@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-3@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 recurse]
		          deleted [tag-3.etikett]
	EOM

	run_zit organize -mode commit-directly <<-EOM
		- [tag-3-sub]
	EOM

	assert_success
	assert_output - <<-EOM
		[tag-3-sub@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
	EOM

	run_zit show recurse:e,z
	assert_success
	assert_output_unsorted - <<-EOM
		[recurse@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-3@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 recurse]
		[tag-3-sub@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM
}
