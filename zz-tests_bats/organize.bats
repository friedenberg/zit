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

function organize_simple { # @test
	actual="$(mktemp)"
	run_zit organize -mode output-only :z,e,t >"$actual"
	assert_success
	assert_output_unsorted - <<-EOM

		- [   !md   ]
		- [   -tag  ]
		- [   -tag-1]
		- [   -tag-2]
		- [   -tag-3]
		- [   -tag-4]
		- [one/dos  ] wow ok again
		- [one/uno  ] wow the first
	EOM
}

function organize_simple_commit { # @test
	run_zit organize -mode commit-directly :z,e,t <<-EOM
		# new-etikett-for-all
		- [   !md   ]
		- [   -tag  ]
		- [   -tag-1]
		- [   -tag-2]
		- [   -tag-3]
		- [   -tag-4]
		- [one/dos  ] wow ok again
		- [one/uno  ] wow the first
	EOM
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[-new-etikett-for-all@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-new-etikett-for@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-new-etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-new@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-1@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-2@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-3@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-4@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again"]
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first"]
	EOM

	run_zit show -format log new-etikett-for-all:z,e,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384 new-etikett-for-all]
		[-new-etikett-for-all@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-1@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-2@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-3@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-4@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again"]
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first"]
	EOM
}

function organize_dry_run { # @test
	expected_show="$(mktemp)"
	# shellcheck disable=SC2154
	zit show "${cmd_zit_def[@]}" -format log :z,e,t >"$expected_show"

	run_zit organize -dry-run -mode commit-directly :z,e,t <<-EOM
		# new-etikett-for-all
		- [   !md   ]
		- [   -tag  ]
		- [   -tag-1]
		- [   -tag-2]
		- [   -tag-3]
		- [   -tag-4]
		- [one/dos  ] wow ok again
		- [one/uno  ] wow the first
	EOM
	assert_success

	run_zit show -format log :z,e,t
	assert_success
	assert_output_unsorted "$(cat "$expected_show")"
}
