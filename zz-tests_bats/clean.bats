#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	copy_from_version "$DIR" "$version"

	run_zit checkout @z,t,e
}

teardown() {
	rm_from_version "$version"
}

function clean_all { # @test
	run_zit clean .
	assert_output_unsorted - <<-EOM
		           (deleted) [md.typ]
		           (deleted) [one/dos.zettel]
		           (deleted) [one/uno.zettel]
		           (deleted) [one]
	EOM

	run find . -type d ! -ipath './.zit*'
	assert_output '.'
}
