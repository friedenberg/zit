#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	copy_from_version "$DIR" "$version"
}

teardown() {
	rm_from_version "$version"
}

function new_empty_no_edit { # @test
	run_zit new -edit=false
	assert_success
	assert_output - <<-EOM
		[two/uno@e6e789716abc939fc15b8caae85ecb9c1bbe96d44d1b58d2fd42a2a8fd9d904a !md ]
	EOM
}
