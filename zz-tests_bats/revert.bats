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

function revert_one { # @test
	run_zit revert one/dos.zettel
	assert_success
	assert_output - <<-EOM
		[-etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-etikett-one@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[one/uno@d2b258fadce18f2de6356bead0c773ca785237cad5009925a3cf1a77603847fc !md "wildly different" etikett-one]
	EOM
}
