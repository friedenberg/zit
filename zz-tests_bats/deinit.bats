#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit info store-version)"
	copy_from_version "$DIR" "$version"
	run_zit_init_workspace
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

  # TODO determine why this zit init is not emitting the created objects
	run_zit_init
	assert_success
	# assert_output - <<-EOM
	# 	[!md @$(get_type_blob_sha) !toml-type-v1]
	# 	[konfig @$(get_konfig_sha) !toml-config-v1]
	# EOM

	run_zit last
	assert_success
	assert_output - <<-EOM
		[!md @$(get_type_blob_sha) !toml-type-v1]
		[konfig @$(get_konfig_sha) !toml-config-v1]
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
