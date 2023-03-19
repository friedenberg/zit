#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	copy_from_version "$DIR" "$version"

	# 	run_zit checkout @z,t,e

	# 	cat >one/uno.zettel <<-EOM
	# 		---
	# 		# wildly different
	# 		- etikett-one
	# 		! md
	# 		---

	# 		newest body
	# 	EOM

	# 	cat >one/dos.zettel <<-EOM
	# 		---
	# 		# dos wildly different
	# 		- etikett-two
	# 		! md
	# 		---

	# 		dos newest body
	# 	EOM

	# 	cat >md.typ <<-EOM
	# 		inline-akte = false
	# 		vim-syntax-type = "test"
	# 	EOM

	# 	cat >da-new.typ <<-EOM
	# 		inline-akte = true
	# 		vim-syntax-type = "da-new"
	# 	EOM

	# 	cat >zz-archive.etikett <<-EOM
	# 		hide = true
	# 	EOM
}

teardown() {
	rm_from_version "$version"
}

function show_simple_one_zettel { # @test
	run_zit show one/uno.zettel
	assert_output - <<-EOM
		---
		# wow the first
		- tag-3
		- tag-4
		! md
		---

		last time
	EOM
}

function show_zettel_etikett { # @test
	run_zit show -format log tag-3.z
	assert_output_unsorted - <<-EOM
		[one/uno@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
		[one/dos@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
	EOM

	run_zit show -format akte -- -tag-3.z
	assert_output_unsorted - <<-EOM
		last time
		not another one
	EOM

	run_zit show -format sku2 -- -tag-3.z
	assert_output_unsorted - <<-EOM
		2057838301.857055 Zettel one/dos c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24
		2057838301.888328 Zettel one/uno d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11
	EOM
}

function show_complex_zettel_etikett_negation { # @test
	run_zit show -format log ^-etikett-two.z
	assert_output_unsorted - <<-EOM
		[one/uno@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
		[one/dos@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
	EOM
}

function show_simple_all { # @test
	run_zit show -format log @z,t
	assert_output_unsorted - <<-EOM
		[!md@eaa85e80de6d1129a21365a8ce2a49ca752457d10932a7d73001b4ebded302c7]
		[one/uno@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
		[one/dos@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
	EOM

	run_zit show -format akte @z,t
	assert_output_unsorted - <<-EOM
		file-extension = 'md'
		inline-akte = true
		last time
		not another one
		vim-syntax-type = 'markdown'
	EOM

	run_zit show -format sku2 @z,t
	assert_output_unsorted - <<-EOM
		2057838301.803584 Typ md eaa85e80de6d1129a21365a8ce2a49ca752457d10932a7d73001b4ebded302c7 102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384
		2057838301.857055 Zettel one/dos c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24
		2057838301.888328 Zettel one/uno d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11
	EOM
}

function show_simple_typ { # @test
	skip
	run_zit show .t
	assert_output_unsorted - <<-EOM
		           (changed) [md.typ@72d654e3c7f4e820df18c721177dfad38fe831d10bca6dcb33b7cad5dc335357 !md]
		         (untracked) [da-new.typ@0ed0c5d77f38816283174202947f71460a455e81b43348bf7808e2b2d81ad120 !da-new]
	EOM
}

function show_simple_etikett { # @test
	skip
	run_zit show .e
	assert_output - <<-EOM
		         (untracked) [zz-archive.etikett@cba019d4f889027a3485e56dd2080c7ba0fa1e27499c24b7ec08ad80ef55da9d -zz-archive]
	EOM
}
