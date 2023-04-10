#! /usr/bin/env bats

setup() {
	load "$BATS_CWD/zz-tests_bats/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	copy_from_version "$DIR" "$version"
}

teardown() {
	rm_from_version "$version"
}

function show_simple_one_zettel { # @test
	run_zit show one/uno.zettel
	assert_success
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

function show_history_one_zettel { # @test
	run_zit show -format log one/uno+z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/uno@797cbdf8448a2ea167534e762a5025f5a3e9857e1dd06a3b746d3819d922f5ce !md "wow ok"]
		[one/uno@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
	EOM

	run_zit show one/uno+z
	assert_success
	assert_output_unsorted - <<-EOM
		---
		# wow ok
		- tag-1
		- tag-2
		! md
		---

		this is the body aiiiiight
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
	assert_success
	assert_output_unsorted - <<-EOM
		[one/uno@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
		[one/dos@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
	EOM

	run_zit show -format akte -- -tag-3.z
	assert_success
	assert_output_unsorted - <<-EOM
		last time
		not another one
	EOM

	run_zit show -format sku2 -- -tag-3.z
	assert_success
	assert_output_cut -d' ' -f2- -- --sort - <<-EOM
		2059300268.968327 Zettel one/dos c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24
		2059300269.066914 Zettel one/uno d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11
	EOM
}

function show_zettel_etikett_complex { # @test
	run_zit checkout o/u
	assert_success

	cat >one/uno.zettel <<-EOM
		---
		# wow the first
		- tag-3
		- tag-5
		! md
		---

		last time
	EOM
	run_zit checkin -delete one/uno.zettel

	run_zit show -format log tag-3.z tag-5.z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/uno@c3232cc6a4122368757d0af489e471e138eab3133ff9107372f33eaf0e284190 !md "wow the first"]
	EOM

	run_zit show -format akte tag-3.z tag-5.z
	assert_success
	assert_output_unsorted - <<-EOM
		last time
	EOM

	run_zit show -format sku2 tag-3.z tag-5.z
	assert_success
	assert_output_unsorted --partial - <<-EOM
		Zettel one/uno c3232cc6a4122368757d0af489e471e138eab3133ff9107372f33eaf0e284190 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11
	EOM
}

function show_complex_zettel_etikett_negation { # @test
	run_zit show -format log ^-etikett-two.z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/uno@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
		[one/dos@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
	EOM
}

function show_simple_all { # @test
	run_zit show -format log :z,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@eaa85e80de6d1129a21365a8ce2a49ca752457d10932a7d73001b4ebded302c7]
		[one/uno@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
		[one/dos@c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 !md "wow ok again"]
	EOM

	run_zit show -format akte :z,t
	assert_success
	assert_output_unsorted - <<-EOM
		file-extension = 'md'
		inline-akte = true
		last time
		not another one
		vim-syntax-type = 'markdown'
	EOM

	run_zit show -format sku2 :z,t
	assert_success
	assert_output_cut -d' ' -f2- -- --sort - <<-EOM
		2059300268.824939 Typ md eaa85e80de6d1129a21365a8ce2a49ca752457d10932a7d73001b4ebded302c7 102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384
		2059300268.968327 Zettel one/dos c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24
		2059300269.066914 Zettel one/uno d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11
	EOM
}

function show_simple_typ { # @test
	run_zit show -format log .t
	assert_output_unsorted - <<-EOM
		[!md@eaa85e80de6d1129a21365a8ce2a49ca752457d10932a7d73001b4ebded302c7]
	EOM
}

function show_simple_etikett { # @test
	run_zit show -format log .e
	assert_output_unsorted - <<-EOM
		[-tag-2@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-tag-3@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-tag-4@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-tag@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
		[-tag-1@5dbb297b5bde513be49fde397499eb89af8f5295f5137d75b52b015802b73ae0]
	EOM
}
