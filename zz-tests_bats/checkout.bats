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

function checkout_simple_all { # @test
	run_zit checkout :z,t,e
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [md.typ@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
		      checked out [one/dos.zettel@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		      checked out [one/uno.zettel@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
		      checked out [tag-1.etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		      checked out [tag-2.etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		      checked out [tag-3.etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		      checked out [tag-4.etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		      checked out [tag.etikett@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
	EOM
}

function checkout_simple_zettel { # @test
	run_zit checkout :
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [one/dos.zettel@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		      checked out [one/uno.zettel@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM
}

function checkout_binary_simple_zettel { # @test
	echo "binary file" >file.bin
	run_zit add -delete file.bin
	assert_success
	assert_output_unsorted - <<-EOM
		          deleted [file.bin]
		[!bin@e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855]
		[two/uno@b20c8fea8cb3e467783c5cdadf0707124cac5db72f9a6c3abba79fa0a42df627 !bin "file"]
	EOM

	run_zit checkout !bin:z
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [two/uno.zettel@b20c8fea8cb3e467783c5cdadf0707124cac5db72f9a6c3abba79fa0a42df627 !bin "file"]
	EOM

	run cat two/uno.zettel
	assert_success
	assert_output - <<-EOM
		---
		# file
		! b20c8fea8cb3e467783c5cdadf0707124cac5db72f9a6c3abba79fa0a42df627.bin
		---
	EOM
}

function checkout_simple_zettel_akte_only { # @test
	run_zit clean .
	assert_success
	# TODO fail checkouts if working directly has incompatible checkout
	run_zit checkout -mode blob :z
	assert_success
	assert_output_unsorted - <<-EOM
		                   one/dos.md]
		                   one/uno.md]
		      checked out [one/dos@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4
		      checked out [one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4
	EOM
}

function checkout_zettel_several { # @test
	run_zit checkout one/uno one/dos
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [one/dos.zettel@2d36c504bb5f4c6cc804c63c983174a36303e1e15a3a2120481545eec6cc5f24 !md "wow ok again" tag-3 tag-4]
		      checked out [one/uno.zettel@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM
}

function checkout_simple_typ { # @test
	run_zit checkout :t
	assert_success
	assert_output_unsorted - <<-EOM
		      checked out [md.typ@102bc5f72997424cf55c6afc1c634f04d636c9aa094426c95b00073c04697384]
	EOM
}

function checkout_zettel_akte_then_objekte { # @test
	run_zit checkout -mode blob one/uno
	assert_success
	assert_output - <<-EOM
		      checked out [one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4
		                   one/uno.md]
	EOM

	run_zit checkout one/uno
	assert_success
	assert_output - <<-EOM
		      checked out [one/uno.zettel@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4]
	EOM

	run ls one/
	assert_output - <<-EOM
		uno.zettel
	EOM

	run_zit checkout -mode blob one/uno
	assert_success
	assert_output - <<-EOM
		      checked out [one/uno@11e1c0499579c9a892263b5678e1dfc985c8643b2d7a0ebddcf4bd0e0288bc11 !md "wow the first" tag-3 tag-4
		                   one/uno.md]
	EOM

	run ls one/
	assert_output - <<-EOM
		uno.md
	EOM
}
