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

	run_zit last -format inventory-list-without-tai
	assert_success
	assert_output_unsorted - <<-EOM
		---
		---
		---
		Akte 102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384
		Akte $(get_konfig_sha)
		Gattung Konfig
		Gattung Typ
		Kennung konfig
		Kennung md
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
		[!md @220519ab7c918ccbd73c2d4d73502ab2ec76106662469feea2db8960b5d68217]
	EOM

	run bash -c 'find .xdg/data/zit/objects/inventory_lists -type f | wc -l | tr -d " "'
	assert_success
	assert_output '2'

	run_zit last -format inventory-list-without-tai
	assert_success
	assert_output - <<-EOM
		---
		Akte 220519ab7c918ccbd73c2d4d73502ab2ec76106662469feea2db8960b5d68217
		Gattung Typ
		Kennung md
		---
	EOM
}
