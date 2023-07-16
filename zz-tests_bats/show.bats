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
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first"]
		[one/uno@3aa85276929951b03184a038ca0ad67cba78ae626f2e3510426b5a17a56df955 !md "wow ok"]
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
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again"]
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first"]
	EOM

	run_zit show -format akte -- -tag-3.z
	assert_success
	assert_output_unsorted - <<-EOM
		last time
		not another one
	EOM

	run_zit show -format sku -- -tag-3.z
	assert_success
	assert_output_cut -d' ' -f2- -- --sort - <<-EOM
		2059300268.968327 Zettel one/dos c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24
		2059300269.066914 Zettel one/uno d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11
	EOM
}

function show_zettel_etikett_complex { # @test
	run_zit checkout o/u
	assert_success
	assert_output - <<-EOM
		      checked out [one/uno.zettel@d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 !md "wow the first"]
	EOM

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
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first"]
	EOM

	run_zit show -format akte tag-3.z tag-5.z
	assert_success
	assert_output_unsorted - <<-EOM
		last time
	EOM

	run_zit show -format sku tag-3.z tag-5.z
	assert_success
	assert_output_unsorted --partial - <<-EOM
		Zettel one/uno c3232cc6a4122368757d0af489e471e138eab3133ff9107372f33eaf0e284190 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11
	EOM
}

function show_complex_zettel_etikett_negation { # @test
	run_zit show -format log ^-etikett-two.z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again"]
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first"]
	EOM
}

function show_simple_all { # @test
	run_zit show -format log :z,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again"]
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first"]
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

	run_zit show -format sku :z,t
	assert_success
	assert_output_cut -d' ' -f2- -- --sort - <<-EOM
		 Typ md b986c1d21fcfb7f0fe11ae960236e3471b4001029a9e631d16899643922b2d15 102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384
		 Zettel one/dos c6b9d095358b8b26a99e90496d916ba92a99e9b75c705165df5f6d353a949ea9 2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24
		 Zettel one/uno d47c552a5299f392948258d7959fc7cf94843316a21c8ea12854ed84a8c06367 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11
	EOM
}

function show_simple_typ_schwanzen { # @test
	run_zit show -format log .t
	assert_output_unsorted - <<-EOM
		[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
	EOM
}

function show_simple_typ_history { # @test
	run_zit show -format log +t
	assert_output_unsorted - <<-EOM
		[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
	EOM
}

function show_simple_etikett_schwanzen { # @test
	run_zit show -format log .e
	assert_output_unsorted - <<-EOM
		[-tag-2@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-3@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-4@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-1@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
	EOM
}

function show_simple_etikett_history { # @test
	run_zit show -format log +e
	assert_output_unsorted - <<-EOM
		[-tag-2@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-3@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-4@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
		[-tag-1@48cae50776cad1ddf3e711579e64a1226ae188ddaa195f4eb8cf6d8f32774249]
	EOM
}

function show_konfig { # @test
	run_zit show -format log +konfig
	assert_output_unsorted - <<-EOM
		[konfig@c1a8ed3cf288dd5d7ccdfd6b9c8052a925bc56be2ec97ed0bb345ab1d961c685]
		[konfig@c1a8ed3cf288dd5d7ccdfd6b9c8052a925bc56be2ec97ed0bb345ab1d961c685]
	EOM
}

function show_history_all { # @test
	run_zit show -format bestandsaufnahme-sans-tai +konfig,kasten,typ,etikett,zettel
	assert_output_unsorted - <<-EOM
		---
		---
		---
		---
		---
		---
		---
		---
		---
		---
		---
		---
		---
		Akte 102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384
		Akte 102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384
		Akte 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11
		Akte 2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24
		Akte 3aa85276929951b03184a038ca0ad67cba78ae626f2e3510426b5a17a56df955
		Akte c1a8ed3cf288dd5d7ccdfd6b9c8052a925bc56be2ec97ed0bb345ab1d961c685
		Akte c1a8ed3cf288dd5d7ccdfd6b9c8052a925bc56be2ec97ed0bb345ab1d961c685
		Bezeichnung wow ok
		Bezeichnung wow ok again
		Bezeichnung wow the first
		Etikett tag-1
		Etikett tag-2
		Etikett tag-3
		Etikett tag-3
		Etikett tag-4
		Etikett tag-4
		Gattung Etikett
		Gattung Etikett
		Gattung Etikett
		Gattung Etikett
		Gattung Etikett
		Gattung Konfig
		Gattung Konfig
		Gattung Typ
		Gattung Typ
		Gattung Zettel
		Gattung Zettel
		Gattung Zettel
		Kennung konfig
		Kennung konfig
		Kennung md
		Kennung md
		Kennung one/dos
		Kennung one/uno
		Kennung one/uno
		Kennung tag
		Kennung tag-1
		Kennung tag-2
		Kennung tag-3
		Kennung tag-4
		Typ md
		Typ md
		Typ md
	EOM
}
