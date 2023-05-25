#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

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
		Tai Typ md b986c1d21fcfb7f0fe11ae960236e3471b4001029a9e631d16899643922b2d15 102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384
		Tai Konfig konfig 62c02b6f59e6de576a3fcc1b89db6e85b75c2ff7820df3049a5b12f9db86d1f5 c1a8ed3cf288dd5d7ccdfd6b9c8052a925bc56be2ec97ed0bb345ab1d961c685
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
		[!md@acbfc0e07b1be4bf1b12020d8316fe9629518b015041b7120db5a9f2012c84fa]
	EOM

	run bash -c 'find .zit/Objekten/Bestandsaufnahme -type f | wc -l | tr -d " "'
	assert_success
	assert_output '3'

	run_zit last
	assert_success
	assert_output_cut -d' ' -f2- -- - <<-EOM
		Tai Typ md acbfc0e07b1be4bf1b12020d8316fe9629518b015041b7120db5a9f2012c84fa 220519ab7c918ccbd73c2d4d73502ab2ec76106662469feea2db8960b5d68217
	EOM
}
