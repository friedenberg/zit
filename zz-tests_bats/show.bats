#! /usr/bin/env bats

setup() {
	load "$(dirname "$BATS_TEST_FILE")/common.bash"

	# for shellcheck SC2154
	export output

	version="v$(zit store-version)"
	copy_from_version "$DIR" "$version"
}

teardown() {
	rm_from_version "$version"
}

function show_simple_one_zettel { # @test
	run_zit show -format text one/uno.zettel
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

function show_simple_one_zettel_binary { # @test
	echo "binary file" >file.bin
	run_zit add -delete file.bin
	assert_success
	assert_output_unsorted - <<-EOM
		          deleted [file.bin]
		[!bin@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[two/uno@b20c8fea8cb3e467783c5cdadf0707124cac5db72f9a6c3abba79fa0a42df627 !bin "file"]
	EOM

	run_zit show -format text two/uno
	assert_success
	assert_output - <<-EOM
		---
		# file
		! b20c8fea8cb3e467783c5cdadf0707124cac5db72f9a6c3abba79fa0a42df627.bin
		---
	EOM
}

function show_history_one_zettel { # @test
	run_zit show -format log one/uno+z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		[one/uno@3aa85276929951b03184a038ca0ad67cba78ae626f2e3510426b5a17a56df955 !md "wow ok" tag-1 tag-2]
	EOM

	run_zit show -format text one/uno+z
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
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM

	run_zit show -format akte -- -tag-3.z
	assert_success
	assert_output_unsorted - <<-EOM
		last time
		not another one
	EOM

	run_zit show -format sku-metadatei-sans-tai -- -tag-3.z
	assert_success
	assert_output_unsorted - <<-EOM
		Zettel one/uno 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md tag-3 tag-4 "wow the first"
		Zettel one/dos 2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md tag-3 tag-4 "wow ok again"
	EOM
}

function show_zettel_etikett_complex { # @test
	run_zit checkout o/u
	assert_success
	assert_output - <<-EOM
		      checked out [one/uno.zettel@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first"]
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
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-5]
	EOM

	run_zit show -format akte tag-3.z tag-5.z
	assert_success
	assert_output_unsorted - <<-EOM
		last time
	EOM

	run_zit show -format sku-metadatei-sans-tai tag-3.z tag-5.z
	assert_success
	assert_output_unsorted --partial - <<-EOM
		Zettel one/uno 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md tag-3 tag-5 "wow the first"
	EOM
}

function show_complex_zettel_etikett_negation { # @test
	run_zit show -format log ^-etikett-two.z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM
}

function show_simple_all { # @test
	run_zit show -format log :z,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		[one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
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

	run_zit show -format sku-metadatei-sans-tai :z,t
	assert_success
	assert_output_unsorted - <<-EOM
		Typ md 102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384
		Zettel one/dos 2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md tag-3 tag-4 "wow ok again"
		Zettel one/uno 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md tag-3 tag-4 "wow the first"
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
	EOM
}

function show_simple_etikett_schwanzen { # @test
	run_zit show -format log .e
	assert_output_unsorted - <<-EOM
		[-tag-1@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-2@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-3@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-4@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
	EOM
}

function show_simple_etikett_history { # @test
	run_zit show -format log +e
	assert_output_unsorted - <<-EOM
		[-tag-1@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-2@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-3@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag-4@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[-tag@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
	EOM
}

function show_konfig { # @test
	run_zit show -format log +konfig
	assert_output_unsorted - <<-EOM
		[konfig@$(get_konfig_sha)]
	EOM

	run_zit show -format text :konfig
	assert_output - <<-EOM
		hidden-etiketten = []

		[defaults]
		typ = 'md'
		etiketten = []

		[file-extensions]
		zettel = 'zettel'
		organize = 'md'
		typ = 'typ'
		etikett = 'etikett'
		kasten = 'kasten'

		[cli-output]
		print-include-typen = false
		print-include-bezeichnungen = false
		print-time = false
		print-etiketten-always = false
		print-empty-shas = false
		print-matched-archiviert = false
		print-shas = false

		[cli-output.abbreviations]
		hinweisen = false
		shas = false

		[tools]
		merge = ['vimdiff']
	EOM
}

function show_history_all { # @test
	run_zit show -format sku-metadatei-sans-tai +konfig,kasten,typ,etikett,zettel
	assert_output_unsorted - <<-EOM
		Etikett tag e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
		Etikett tag-1 e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
		Etikett tag-2 e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
		Etikett tag-3 e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
		Etikett tag-4 e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
		Konfig konfig $(get_konfig_sha)
		Typ md 102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384
		Zettel one/dos 2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md tag-3 tag-4 "wow ok again"
		Zettel one/uno 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md tag-3 tag-4 "wow the first"
		Zettel one/uno 3aa85276929951b03184a038ca0ad67cba78ae626f2e3510426b5a17a56df955 !md tag-1 tag-2 "wow ok"
	EOM
}
