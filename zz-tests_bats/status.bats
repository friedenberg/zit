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

	cat >da-new.typ <<-EOM
		inline-akte = true
		vim-syntax-type = "da-new"
	EOM

	cat >zz-archive.etikett <<-EOM
		hide = true
	EOM
}

teardown() {
	rm_from_version "$version"
}

function status_simple_one_zettel { # @test
	run_zit status one/uno.zettel
	assert_output - <<-EOM
		           (changed) [one/uno.zettel@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]
	EOM
}

function status_zettelen_typ { # @test
	run_zit status !md.z
	assert_output_unsorted - <<-EOM
		           (changed) [one/dos.zettel@30edfed4c016580f5b69a2709b8e5ae01c2b504b8826bf2d04e6c1ecd6bb3268 !md "dos wildly different"]
		           (changed) [one/uno.zettel@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]
	EOM
}

function status_complex_zettel_etikett_negation { # @test
	run_zit status ^-etikett-two.z
	assert_output_unsorted - <<-EOM
		           (changed) [one/uno.zettel@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]
	EOM
}

function status_simple_all { # @test
	run_zit status .
	assert_output_unsorted - <<-EOM
		           (changed) [md.typ@72d654e3c7f4e820df18c721177dfad38fe831d10bca6dcb33b7cad5dc335357 !md]
		           (changed) [one/dos.zettel@30edfed4c016580f5b69a2709b8e5ae01c2b504b8826bf2d04e6c1ecd6bb3268 !md "dos wildly different"]
		           (changed) [one/uno.zettel@689c6787364899defa77461ff6a3f454ca667654653f86d5d44f2826950ff4f9 !md "wildly different"]
		         (untracked) [da-new.typ@0ed0c5d77f38816283174202947f71460a455e81b43348bf7808e2b2d81ad120 !da-new]
		         (untracked) [zz-archive.etikett@cba019d4f889027a3485e56dd2080c7ba0fa1e27499c24b7ec08ad80ef55da9d -zz-archive]
	EOM
}

function status_simple_typ { # @test
	run_zit status .t
	assert_output_unsorted - <<-EOM
		           (changed) [md.typ@72d654e3c7f4e820df18c721177dfad38fe831d10bca6dcb33b7cad5dc335357 !md]
		         (untracked) [da-new.typ@0ed0c5d77f38816283174202947f71460a455e81b43348bf7808e2b2d81ad120 !da-new]
	EOM
}

function status_simple_etikett { # @test
	run_zit status .e
	assert_output - <<-EOM
		         (untracked) [zz-archive.etikett@cba019d4f889027a3485e56dd2080c7ba0fa1e27499c24b7ec08ad80ef55da9d -zz-archive]
	EOM
}
