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

function mergetool_none { # @test
	run_zit merge-tool .
	assert_success
	assert_output "nothing to merge"
}

function mergetool_conflict_one { # @test
	run_zit checkout one/dos
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [one/dos.zettel@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again"]
	EOM

	cat - >one/dos.zettel <<-EOM
		---
		# wow ok again
		- get_this_shit_merged
		- tag-3
		- tag-4
		! txt
		---

		not another one, conflict time
	EOM

	run_zit organize -mode commit-directly one/dos <<-EOM
		---
		! txt2
		---

		# new-etikett-for-all
		- [one/dos  ] wow ok again
	EOM
	assert_success
	assert_output_unsorted - <<-EOM
		[!txt2@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-new-etikett-for-all@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-new-etikett-for@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-new-etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-new@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !txt2 "wow ok again" new-etikett-for-all tag-3 tag-4]
	EOM

	run_zit show -format log new-etikett-for-all:z,e,t
	assert_success
	assert_output_unsorted - <<-EOM
		[-new-etikett-for-all@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !txt2 "wow ok again" new-etikett-for-all tag-3 tag-4]
	EOM

	run_zit status one/dos.zettel
	assert_success
	assert_output - <<-EOM
		       conflicted [one/dos.zettel]
	EOM

	run_zit merge-tool -merge-tool "bash -c 'echo merged >\"\$MERGED\"'" .
	assert_success
	assert_output - <<-EOM
		[one/dos@72d8264b97bee169d1844d054282694b68be7e91c8bd5d616540adf63ae7d4af]
	EOM

	run_zit last
	assert_success
	assert_output - <<-EOM
		[one/dos@72d8264b97bee169d1844d054282694b68be7e91c8bd5d616540adf63ae7d4af]
	EOM
}