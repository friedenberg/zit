#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	copy_from_version "$DIR" "$version"

	run_zit checkout @z,t,e

	cat >one/uno.zettel <<-EOM
		---
		# wildly different
		- etikett-one
		! md
		---

		newest body
	EOM

	cat >one/dos.zettel <<-EOM
		---
		# dos wildly different
		- etikett-two
		! md
		---

		dos newest body
	EOM

	cat >md.typ <<-EOM
		inline-akte = false
		vim-syntax-type = "test"
	EOM

	cat >zz-archive.etikett <<-EOM
		hide = true
	EOM
}

teardown() {
	rm_from_version "$version"
}

function checkin_simple_one_zettel { # @test
	run_zit checkin one/uno.zettel
	assert_output '[one/uno@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]'
}

function checkin_complex_zettel_etikett_negation { # @test
	run_zit checkin ^-etikett-two.z
	assert_output '[one/uno@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]'
}

function checkin_simple_all { # @test
	run_zit checkin .
	assert_output_unsorted - <<-EOM
		[!md@72d654e3c7f4e820df18c721177dfad38fe831d10bca6dcb33b7cad5dc335357]
		[one/uno@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]
		[one/dos@30edfed4c016580f5b69a2709b8e5ae01c2b504b8826bf2d04e6c1ecd6bb3268 !md "dos wildly different"]
		[-zz-archive@cba019d4f889027a3485e56dd2080c7ba0fa1e27499c24b7ec08ad80ef55da9d]
	EOM
}

function checkin_simple_typ { # @test
	run_zit checkin .t
	assert_output '[!md@72d654e3c7f4e820df18c721177dfad38fe831d10bca6dcb33b7cad5dc335357]'

	run_zit show -format vim-syntax-type !md.typ
	assert_output 'test'
}

function checkin_simple_etikett { # @test
	run_zit checkin .e
	assert_output '[-zz-archive@cba019d4f889027a3485e56dd2080c7ba0fa1e27499c24b7ec08ad80ef55da9d]'

	run_zit show -format text -- -zz-archive.e
	assert_output 'hide = true'
}
