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
	run_zit show -format text one/uno
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

function show_simple_one_zettel_with_sigil { # @test
	run_zit show -format text one/uno:
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

function show_simple_one_zettel_with_sigil_and_genre { # @test
	run_zit show -format text one/uno:zettel
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

function show_simple_one_zettel_checked_out { # @test
	run_zit checkout one/uno
	assert_success
	assert_output - <<-EOM
		      checked out [one/uno.zettel @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM

	run_zit show one/uno.zettel
	assert_success
	assert_output - <<-EOM
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM
}

function show_simple_one_zettel_hidden { # @test
	run_zit schlummernd-add tag-3
	assert_success
	assert_output ''

	run_zit show :z
	assert_success
	assert_output ''

	run_zit show :?z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM

	run_zit show one/uno
	assert_success
	assert_output - <<-EOM
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM
}

function show_simple_one_zettel_hidden_past { # @test
	run_zit schlummernd-add tag-1
	assert_success
	assert_output ''

	run_zit show :?z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
	EOM

	run_zit show one/uno
	assert_success
	assert_output - <<-EOM
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM
}

function show_all_mutter { # @test
	skip
	run_zit show -format mutter-sha :
	assert_success
	assert_output - <<-EOM
		5b059e2dd36c89f2c7f75b2b6f39573af94e4109ceebabe2814515c9ea30eb98
	EOM
}

function show_simple_one_zettel_binary { # @test
	echo "binary file" >file.bin
	run_zit add -delete file.bin
	assert_success
	assert_output_unsorted - <<-EOM
		          deleted [file.bin]
		[!bin @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 !toml-type-v1]
		[two/uno @b20c8fea8cb3e467783c5cdadf0707124cac5db72f9a6c3abba79fa0a42df627 !bin "file"]
	EOM

	cat >bin.type <<-EOM
		---
		! toml-type-v1
		---

		binary = true
	EOM

	run_zit checkin -delete bin.type
	assert_success
	assert_output_unsorted - <<-EOM
		          deleted [bin.type]
		[!bin @e07d72a74e0a01c23ddeb871751f6fcb43afec5fb81108c157537db96c6c1da0 !toml-type-v1]
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
	run_zit show one/uno+z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		[one/uno @3aa85276929951b03184a038ca0ad67cba78ae626f2e3510426b5a17a56df955 !md "wow ok" tag-1 tag-2]
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
	run_zit show tag-3:z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM

	run_zit show -format blob tag-3:z
	assert_success
	assert_output_unsorted - <<-EOM
		last time
		not another one
	EOM

	run_zit show -format sku-metadata-sans-tai tag-3:z
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
		      checked out [one/uno.zettel @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
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

	# TODO support . operator for checked out
	# run_zit show -verbose tag-3.z tag-5.z
	# assert_success
	# assert_output_unsorted - <<-EOM
	# 	[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-5]
	# EOM

	run_zit checkin -delete one/uno.zettel

	run_zit show [tag-3 tag-5]:z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-5]
	EOM

	run_zit show -format blob [tag-3 tag-5]:z
	assert_success
	assert_output_unsorted - <<-EOM
		last time
	EOM

	run_zit show -format sku-metadata-sans-tai [tag-3 tag-5]:z
	assert_success
	assert_output_unsorted --partial - <<-EOM
		Zettel one/uno 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md tag-3 tag-5 "wow the first"
	EOM
}

function show_complex_zettel_etikett_negation { # @test
	run_zit show ^-etikett-two:z
	assert_success
	assert_output_unsorted - <<-EOM
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM
}

function show_simple_all { # @test
	run_zit show :z,t
	assert_success
	assert_output_unsorted - <<-EOM
		[!md @$(get_type_blob_sha) !toml-type-v1]
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM

	run_zit show -format blob :z,t
	assert_success
	assert_output_unsorted - <<-EOM
		file-extension = 'md'
		last time
		not another one
		vim-syntax-type = 'markdown'
	EOM

	run_zit show -format sku-metadata-sans-tai :z,t
	assert_success
	assert_output_unsorted - <<-EOM
		Typ md $(get_type_blob_sha) !toml-type-v1
		Zettel one/dos 2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md tag-3 tag-4 "wow ok again"
		Zettel one/uno 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md tag-3 tag-4 "wow the first"
	EOM
}

function show_simple_typ_schwanzen { # @test
	run_zit show :t
	assert_output_unsorted - <<-EOM
		[!md @$(get_type_blob_sha) !toml-type-v1]
	EOM
}

function show_simple_typ_history { # @test
	run_zit show +t
	assert_output_unsorted - <<-EOM
		[!md @$(get_type_blob_sha) !toml-type-v1]
	EOM
}

function show_simple_etikett_schwanzen { # @test
	run_zit show :e
	assert_output_unsorted - <<-EOM
		[tag-1 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-2 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-3 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-4 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
	EOM
}

function show_simple_etikett_history { # @test
	run_zit show +e
	assert_output_unsorted - <<-EOM
		[tag-1 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-2 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-3 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag-4 @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
	EOM
}

function show_konfig { # @test
	run_zit show +konfig
	assert_output_unsorted - <<-EOM
		[konfig @$(get_konfig_sha) !toml-config-v1]
	EOM

	run_zit show -format text :konfig
	assert_output - <<-EOM
		[defaults]
		type = 'md'
		tags = []

		[file-extensions]
		zettel = 'zettel'
		organize = 'md'
		type = 'type'
		tag = 'tag'
		repo = 'repo'

		[cli-output]
		print-include-description = false
		print-time = false
		print-etiketten-always = false
		print-empty-shas = false
		print-include-typen = false
		print-matched-archiviert = false
		print-shas = false
		print-flush = false
		print-unchanged = false
		print-colors = false
		print-bestandsaufnahme = false

		[cli-output.abbreviations]
		hinweisen = false
		shas = false

		[tools]
		merge = ['vimdiff']
	EOM
}

function show_history_all { # @test
	run_zit show -format sku-metadata-sans-tai +konfig,kasten,typ,etikett,zettel
	assert_output_unsorted - <<-EOM
		Etikett tag e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
		Etikett tag-1 e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
		Etikett tag-2 e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
		Etikett tag-3 e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
		Etikett tag-4 e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855
		Konfig konfig $(get_konfig_sha) !toml-config-v1
		Typ md $(get_type_blob_sha) !toml-type-v1
		Zettel one/dos 2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md tag-3 tag-4 "wow ok again"
		Zettel one/uno 11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md tag-3 tag-4 "wow the first"
		Zettel one/uno 3aa85276929951b03184a038ca0ad67cba78ae626f2e3510426b5a17a56df955 !md tag-1 tag-2 "wow ok"
	EOM
}

function show_etikett_toml { # @test
	cat >true.tag <<-EOM
		---
		! toml-tag-v1
		---

		filter = """
		return {
		  contains_sku = function (sk)
		    return true
		  end
		}
		"""
	EOM

	run_zit checkin -delete true.tag
	assert_success
	assert_output_unsorted - <<-EOM
		          deleted [true.tag]
		[true @1379cb8d553a340a4d262b3be216659d8d8835ad0b4cc48005db8db264a395ed !toml-tag-v1]
	EOM

	run_zit show true
	assert_success
	assert_output_unsorted - <<-EOM
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM
}

# TODO fix race condition between stderr and stdout
function show_etikett_lua_v1 { # @test
	cat >true.tag <<-EOM
		---
		! lua-tag-v1
		---

		return {
		  contains_sku = function (sk)
		    print(Selbst.Kennung)
		    return true
		  end
		}
	EOM

	run_zit checkin -delete true.tag
	assert_success
	assert_output_unsorted - <<-EOM
		          deleted [true.tag]
		[true @67b7eb3e9ea1c4b3404b34a0b2abcc09f450797c8cc801671463a79429aead37 !lua-tag-v1]
	EOM

	run_zit show true
	assert_success
	assert_output_unsorted - <<-EOM
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		true
		true
	EOM
}

# TODO fix race condition between stderr and stdout
function show_etikett_lua_v2 { # @test
	cat >true.tag <<-EOM
		---
		! lua-tag-v2
		---

		return {
		  contains_sku = function (sk)
		    print(Self.ObjectId)
		    return true
		  end
		}
	EOM

	run_zit checkin -delete true.tag
	assert_success
	assert_output_unsorted - <<-EOM
		          deleted [true.tag]
		[true @ed8e3cf53e044fcc1ae040ed5203515d1c6d205decc745f0caafd5dee67efbab !lua-tag-v2]
	EOM

	run_zit show true
	assert_success
	assert_output_unsorted - <<-EOM
		[one/dos @2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		[one/uno @11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		true
		true
	EOM
}

function show_etiketten_paths { # @test
	run_zit show -format etiketten-path :e
	assert_success
	assert_output_unsorted - <<-EOM
		tag [Paths: [TypeSelf:[tag]], All: [tag:[TypeSelf:[tag]]]]
		tag-1 [Paths: [TypeSelf:[tag-1]], All: [tag-1:[TypeSelf:[tag-1]]]]
		tag-2 [Paths: [TypeSelf:[tag-2]], All: [tag-2:[TypeSelf:[tag-2]]]]
		tag-3 [Paths: [TypeSelf:[tag-3]], All: [tag-3:[TypeSelf:[tag-3]]]]
		tag-4 [Paths: [TypeSelf:[tag-4]], All: [tag-4:[TypeSelf:[tag-4]]]]
	EOM
}

function show_etiketten_exact { # @test
	run_zit show =tag:e
	assert_success
	assert_output_unsorted - <<-EOM
		[tag @e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
	EOM

	run_zit show =tag
	assert_success
	assert_output_unsorted ''
}

function show_inventory_lists { # @test
	run_zit show :b
	assert_success
	assert_output
}

# bats test_tags=user_story:builtin_types
function show_builtin_type_md { # @test
	run_zit show -format text !toml-type-v1:t
	assert_success
	assert_output - <<-EOM
		---
		! toml-type-v1
		---

		file-extension = 'md'
		vim-syntax-type = 'markdown'
	EOM
}
