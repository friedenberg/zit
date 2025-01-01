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

# TODO add a preview of what would be deleted
function deinit_force() { # @test
	run_zit deinit -force
	assert_success
	assert_output - <<-EOM
	EOM

	run_zit status
	assert_failure
	assert_output - <<-EOM
		not in a zit directory
	EOM

	run_zit_init
	assert_success
	assert_output - <<-EOM
		[!md @b7ad8c6ccb49430260ce8df864bbf7d6f91c6860d4d602454936348655a42a16 !toml-type-v1]
		[konfig @d904d322213ed86cdc0eabd58d44f55385f9665280f6c03a01e396f22ba2333b !toml-config-v1]
	EOM
}

function deinit() { # @test
	run_zit deinit
	assert_success
	assert_output --regexp - <<-EOM
		are you sure you want to deinit in ".*"? \(y/\*)
		failed to read answer: EOF
	EOM
}
