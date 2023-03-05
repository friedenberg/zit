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
	assert_output --partial '[!md@0966bffa92f9391ec0874fe0bd5ed77b9ceddc45e36a866c71a3ccbb31711a71]'
	assert_output --partial '[one/uno@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]'
	assert_output --partial '[one/dos@30edfed4c016580f5b69a2709b8e5ae01c2b504b8826bf2d04e6c1ecd6bb3268 !md "dos wildly different"]'
}

function checkin_simple_typ { # @test
	run_zit checkin .t
	assert_output '[!md@0966bffa92f9391ec0874fe0bd5ed77b9ceddc45e36a866c71a3ccbb31711a71]'
}
