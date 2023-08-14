#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	# for shellcheck SC2154
	export output

	# version="v$(zit store-version)"
	# copy_from_version "$DIR" "$version"
}

teardown() {
	chflags_and_rm
}

function last_after_init { # @test
	run_zit_init_disable_age

	run_zit last
	assert_success
	assert_output_cut -d' ' -f2- -- - <<-EOM
		Tai Typ md 102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384
		Tai Konfig konfig 40fcab44369d4fe18dedd39d6faf5bedf3004929e0974ee631a56895813f5f8b
	EOM
}

function last_after_typ_mutate { # @test
	run_zit_init_disable_age

	cat >md.typ <<-EOM
		inline-akte = false
		vim-syntax-type = "test"
	EOM

	run_zit checkin .t
	assert_success
	assert_output - <<-EOM
		[!md@220519ab7c918ccbd73c2d4d73502ab2ec76106662469feea2db8960b5d68217]
	EOM

	run bash -c 'find .zit/Objekten2/Bestandsaufnahmen -type f | wc -l | tr -d " "'
	assert_success
	assert_output '3'

	run_zit last
	assert_success
	assert_output_cut -d' ' -f2- -- - <<-EOM
		Tai Typ md 220519ab7c918ccbd73c2d4d73502ab2ec76106662469feea2db8960b5d68217
	EOM
}
