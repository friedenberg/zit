#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	version="v$(zit info store-version)"
	copy_from_version "$DIR" "$version"
}

teardown() {
	rm_from_version "$version"
}

function complete_show_all { # @test
	skip
	run_zit show -complete :
	assert_success
	assert_output_unsorted - <<-EOM
		md\tTyp
		one/dos\tZettel: !md wow ok again
		one/uno\tZettel: !md wow the first
		tag\tEtikett
		tag-1\tEtikett
		tag-2\tEtikett
		tag-3\tEtikett
		tag-4\tEtikett
	EOM
}

function complete_show_zettelen { # @test
	run_zit show -verbose -complete :z
	assert_success
	assert_output_unsorted --regexp - <<-EOM
		one/dos.*Zettel: !md wow ok again
		one/uno.*Zettel: !md wow the first
	EOM
}

function complete_show_typen { # @test
	run_zit show -complete :t
	assert_success
	assert_output_unsorted --regexp - <<-EOM
		md.*Typ
	EOM
}

function complete_show_etiketten { # @test
	run_zit show -complete :e
	assert_success
	assert_output_unsorted --regexp - <<-EOM
		tag-3.*Etikett
		tag-4.*Etikett
	EOM
}
