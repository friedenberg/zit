#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	version="v$(zit info store-version)"
	copy_from_version "$DIR" "$version"
}

teardown() {
	rm_from_version "$version"
}

function complete_show { # @test
	run_zit show -complete
	assert_success
	assert_output_unsorted --regexp - <<-'EOM'
		tag-1.*Etikett
		tag-2.*Etikett
		tag-3.*Etikett
		tag-4.*Etikett
		tag.*Etikett
	EOM
}

function complete_show_all { # @test
	run_zit show -complete :z,t,b,e
	assert_success
	assert_output_unsorted --regexp - <<-EOM
		.*Bestandsaufnahme
		.*Bestandsaufnahme
		.*Bestandsaufnahme
		.*Bestandsaufnahme
		!md.*Typ
		one/dos.*Zettel: !md wow ok again
		one/uno.*Zettel: !md wow the first
		tag.*Etikett
		tag.1.*Etikett
		tag.2.*Etikett
		tag.3.*Etikett
		tag.4.*Etikett
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
